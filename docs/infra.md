# How UCC Netsoc Does Infrastructure

This file contains documentation of the redesigned infrastructure for UCC Netsoc (2020/2021). It's ideally targeted at new SysAdmin's or anyone wanting to learn about how we operate

## Historical Infrastructure

**You can ignore this section if you wish, it's just to provide context if you come across old stuff we missed/have not yet migrated**

Prior to 2020/2021, Netsoc ran a completely different setup on bare metal which consisted of 5 bare-metal servers & 2 virtual machines provided by UCC

Bare-metal machines:

* **leela**
  * 84.39.234.51
  * Our main user server that any of our members can ssh into
  * Hosted netsocadmin (v2)
  * UCC Express & Motley
* **bigbertha**
  * 84.39.234.50
  * CI/related
* **lovelace**
  * 84.39.234.52
  * Game server hosting
* **boole** 
  * 84.39.234.54
  * Monitoring
* **feynman**
  * 84.39.234.53
  * Root server for any user

These machines are hosted in a rack in Cork International eXchange (CIX).
Many previous and current Netsoc members have worked at CIX and are familiar with the data center

Our point-of-contact with CIX is Jerry McSweeney and noc@cix.ie, lovely chap

Our UCC VMs:
* **elon / netsoc1.ucc.ie**
  * 143.239.87.40
  * LDAP
  * MySQL
* **tesla / netsoc2.ucc.ie**
  * 143.239.87.41
  * I honestly can't remember, it's semi-perma borked

We need to figure out our UCC point-of-contact

The setup consisted of a Docker swarm all of our bare metal servers.

We did not run a router, just a switch. As such each bare metal server was bound directly to their public IPs

We decided it was time to completely manage our infrastructure with Infrastructure-As-Code as it would help preserve the server setup over time. We therefore migrated a large amount of our setup to be configured using Ansible in 2019/2020.

## Core Infrastructure

We run [Proxmox VE](https://pve.proxmox.com/), a system that abstracts KVM virtual machines and LXC containers into a common set of powerful tools.

It's also got a _really_ nice web UI. Think of it like an open-source ESXi/vSphere

Our setup consists of a cluster of Proxmox machines, there currently are as follows:

* feynman
  * Designated as our Ansible control server
    * Designated as such because it's a pretty terrible server
    * All Ansible playbooks that make changes must be run from here
    * Each SysAdmin has a Linux account created on the server to allow them to do development work
    * Fair warning: it's disks are slow af
* lovelace

Servers that still need to become Proxmox cluster hosts (these still run the historical setups):

* bigbertha
* leela
* boole

We have not yet decided on what will be done to the UCC VMs.
We currently do _not_ plan on using Ceph (distributed storage) on our cluster. This is because Ceph requires a minimum quorum of 3 servers to operate redundantly. We may use it in the future

### Short Intro to Ansible (for new SysAdmins)

[Ansible](https://www.ansible.com/overview/how-ansible-works) is an incredibly powerful IT automation tool we use to configure and install software on servers

* Servers to target are defined in a file (normally called `hosts`, or by a script), and given aliases and groups.
  * Our `hosts` looks a little like this:
    ```
    control

    [proxmox_master]
    feynman

    [proxmox_hosts]
    feynman
    lovelace
    ```
  * These hosts & groups can have `host_vars` and `group_vars` which apply to any server / server in that group respectively
    * Here's an example of `host_vars` for `feynman` (`host_vars/feynman.yml`):
      ```
      ansible_host: 10.0.10.53
      ansible_ssh_user: root
      ansible_ssh_private_key_file: "/root/.ssh/id_rsa"
      ansible_python_interpreter: "/usr/bin/python3"
      ```
    * Here's an example of `group_vars` for a group called `vm` (`group_vars/vm.yml`):
      ```
      ansible_ssh_common_args: '-o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null'
      ansible_python_interpreter: "/usr/bin/python3"
      ansible_user: netsoc
      ```
    * Not only can you make use of explicit variables specify info about the host for Ansible, you can specify your own and access them later
* A _playbook_ consists of a list of _plays_ to run on a host / group of hosts
    * An _play_ could be a bunch of _tasks_ or _roles_
      * A _task_ is a call to an Ansible _module_ that does some functionality
        * Example, to install the `jq` package:
          ```
            apt:
              name: jq
              state: latest
              update_cache: yes
          ```
        * It's important to keep the Ansible module documentation open
          * [Example](https://docs.ansible.com/ansible/latest/modules/apt_module.html)
        * There are modules for *almost everything*, before implementing functionality using a script or shell commands try and see if it can be done using an existing module
        * [An example list of some of the most common modules you may use](https://opensource.com/article/19/9/must-know-ansible-modules)
      * A _role_ is a grouping of functionality 
        * It has lists of tasks, files, variables that can be specified when the role is ran
        * Think of it almost like a script that accepts arguments
        * Example of a role that installs and starts nginx:
        ```
          roles/nginx/vars/defaults.yml:
            start_on_boot: no

          roles/nginx/tasks/main.yml:
          - name: Install nginx package
            apt:
              name: nginx
              state: latest
              update_cache: yes
          
          - name: Start Nginx service (and enable service start on boot if specified)
            service:
              name: nginx
              enabled: "{{ start_on_boot }}"
              state: started
        ```
      * You could call the previous role as a _play_ in a playbook like so:
        ```
        - name: "Ensure nginx on lovelace"
          hosts: lovelace
          roles:
            - role: nginx
              vars:
                start_on_boot: yes
        ```
          * This is the explicit var syntax that scopes the vars to the role, you can also scope the vars to the play
      * Here is a list of tasks as a _play_:
        ```
        - name: Setup server MOTD
          hosts: portal
          tasks:
            - become: yes
              copy:
                content: "{{ lookup('file', './portal-banner.txt') }}"
                dest: /etc/netsoc-banner
            - become: yes
              copy:
                content: "{{ lookup('file', './portal-motd.txt') }}"
                dest: /etc/netsoc-motd
            - become: yes
              shell: /bin/rm -rf /etc/update-motd.d/*
            - become: yes
              copy:
                content: |
                  #!/bin/sh
                  cat /etc/netsoc-motd
                mode: "0755"
                dest: /etc/update-motd.d/00-netsoc-motd
            - become: yes
              copy:
                content: ""
                dest: /etc/legal
            - become: yes
              lineinfile:
                regex: "^#PrintLastLog"
                line: "PrintLastLog no"
                path: /etc/ssh/sshd_config
            - become: yes
              lineinfile:
                regex: "^#Banner none"
                line: "Banner /etc/netsoc-banner"
                path: /etc/ssh/sshd_config
            - become: yes
              service:
                name: sshd
                state: restarted
        ```
        * Don't forget that tasks run after roles if you use them together in a play
      * Ansible supports Jinja variable interpolation almost everywhere
        * Don't forget to look at [filters](https://docs.ansible.com/ansible/latest/user_guide/playbooks_filters.html#list-filters)
  * You can run a playbook by doing `ansible-playbook -i hosts playbook.yml`
  * The playbook functionality is a complete IaC addition to Ansible, but you don't actually need to use it
    * Ping every server: `ansible -i hosts all -m ping`
    * Chmod a.txt to 600:  `ansible -i hosts all -m file -a "dest=/srv/foo/a.txt mode=600"`
* Ansible also has a _vault_ for secrets
  * This is like a yaml file full of variables we don't reveal, you can edit the vault and when you finish editing it Ansible will encrypt the file
  * This lets us store our secrets in GitHub but not reveal them to anyone
  * You can edit the vault by doing `ansible-vault edit vars/secrets.yml`
  * When you run a playbook, you can enable vault by doing `ansible-playbook -i hosts playbook.yml --ask-vault-pass`
  * This will decrypt every vault variable and prefix them with `vault_`
  * We normally realias them with the files in `vars/`
    * Example:
    ```
    cloudflare_api_email: "{{ vault_cloudflare_api_email }}"
    cloudflare_api_key: "{{ vault_cloudflare_api_key }}"
    cloudflare_dns_api_token: "{{ vault_cloudflare_dns_api_token }}"
    ```

* **[Keep all operations idempotent, this is the core philosophy to Ansible](https://docs.ansible.com/ansible/latest/reference_appendices/glossary.html#term-idempotency)**

### How we use IaC (summary)

* We can define and create a "template" VM that contains pre-installed tools and a base image
  * Template VMs are created using [Hashicorp Packer](https://www.packer.io)
  * We install the distro OS and then do further provisioning via Ansible/shell scripts
  * We do not call the Packer executable directly, we instead use an Ansible role we created `packer-proxmox-template`
    * Packer doesn't support saving the VM template to a distributed storage
      * Not that we're using one, yet
    * For a given template i.e `netsoc-centos-6`, we will first build the template on a specified server under the name `netsoc-centos-6-<hostname>`
    * The template is then dumped and copied to every other Proxmox server and hostname renamed
    * When we want to create instances of the template, we clone on that server
    * i.e for `netsoc-ubuntu-server`, we have `netsoc-ubuntu-server-lovelace` and `netsoc-ubuntu-server-feynman`
  * See `template-netsoc-ubuntu-server.yml`

* We can define and create VMs at will
  * VM creation is done by our custom VM Ansible role `proxmox-infra-cloudinit-vm`
  * We create a clone of a template VM defined previously
  * It's expected that the template VM we cloned has cloud-init & qemu-guest-agent installed
    * [cloud-init](https://cloudinit.readthedocs.io/en/latest/) is the de-facto industry standard for configuring Cloud VMs
      * We use it to inject maintenance user info (ssh keys, etc) and network configuration into our VMs
      * Note, the role we use accepts only _plain-text_ yaml, this is because not every field (userdata/instance data/network) only uses yaml. [Look here in the cloud-init docs](https://cloudinit.readthedocs.io/en/latest/topics/format.html) to see what I mean
    * [qemu-guest-agent](https://pve.proxmox.com/wiki/Qemu-guest-agent) is a package that exposes the current network setup and stats to the KVM machine
      * i.e. it let's us see what VMs are _actually_ binding to what static IPs we asked them to
      * also used for some Proxmox related functionality
  * We then modify the clone we created
    * Set specifications:
      * Cores, Memory, Network stuff
    * Resize existing disks
    * Set new disks
      * **If you rerun the playbook, we will only recreate the disks if you explicitly say we should**
        * Prevents loss of user data
    * Set the description:
      * All our IaC managed VMs are expected have descriptions containing YAML.
      * We use this YAML later to "discover" info about any running VMs
      * Example from (`create-auth.yml`):
        ```
        groups:
          - vm
          - auth
        vars:
          ansible_ssh_private_key_file: ./keys/auth/id_rsa
        ```
      * We will use this to pass info about running VMs to Ansible later
    * Supply Cloud-Init metadata / userdata / network config
      * If you're reading cloud-init docs, we use the [NoCloud data source](https://cloudinit.readthedocs.io/en/latest/topics/datasources/nocloud.html)
      * We generally don't use metadata/instance metdata
        * More targeted towards actual clouds like EC2
      * We use userdata to inject the Netsoc maintenance user and set the hostname (using [`#cloud-config`](https://cloudinit.readthedocs.io/en/latest/topics/examples.html))
      * Network config [`(we use v2)`](https://cloudinit.readthedocs.io/en/latest/topics/network-config-format-v2.html#network-config-v2) is used to tell the VM what static IP it should bind to, what DNS server is should use & what router it should send packets to
      * Example from (`create-auth.yml`):
        ```
        cloudinit:
          drive_device: ide2
          metadata: ""
          userdata: |
            #cloud-config
            users:
              - name: netsoc
                gecos: Netsoc Management User
                primary_group: netsoc
                sudo: ALL=(ALL) NOPASSWD:ALL
                ssh_authorized_keys:
                  - {{ lookup('file', './keys/auth/id_rsa.pub') }}
            preserve_hostname: false
            manage_etc_hosts: true
            fqdn: auth.vm.netsoc.co
          networkconfig: |
            version: 2
            ethernets:
              30a:
                match:
                  name: ens18
                addresses:
                  - {{ ip_allocation.freeipa_base }}
                gateway4: 10.0.30.1
                nameservers:
                  search: {{ ip_allocation.search }}
                  addresses: {{ ip_allocation.nameservers }}
        ```
      * `drive_device` specifies what device the NoCloud CD-ROM drive is mounted on
  * See `create-auth.yml`

* We need to be able to discover every VM we have running in order to customize the base image we cloned further
  * Proxmox has an API
  * Description field was set with our previous Ansible-related groups and vars
  * `qemu-guest-agent` will tell us what IPs the machines are running
  * Ansible has a feature known as [custom inventory scripts](https://docs.ansible.com/ansible/latest/user_guide/intro_dynamic_inventory.html#dynamic-inventory)
    * Ansible normally searches for machines to connect to in predefined files (i.e ./hosts, ./hosts.yml, ansible.cfg, etc...)
    * A custom inventory script is a script that instead can be called to return the list of possible machines to connect to
    * It's really common to use with Cloud providers (i.e discover every AWS EC2 instance you have running and run a command on them all)
    * Our custom inventory script for Proxmox is `proxmox_inventory.py`
    * Before you can use it, you must have loaded the environmental variables from the Ansible vault (these contain the API password)
    * You can do this by running `source proxmox_secrets.sh` and entering the vault password
    * You can now see we can set `host_vars` from the VM's description `vars`
    * This is how we discover which ssh key in the NaC directory to use, and other info
    * The inventory script will now work, try `ansible -i proxmox_inventory.py all -m ping` and it should be able to ping all the VMs
    * Similarly you will see that `ansible-playbook -i proxmox_inventory.py ...` will use that inventory
    * You can target the cluster hosts at the same time too by using `ansible-playbook -i hosts -i proxmox_inventory ...`
    * If you want to get the json passed to Ansible for yourself, you can run `./proxmox_inventory.py --list`

* We can then further provision that VM with tools that only it specifically needs
  * See `provision-*.yml`
  * i.e if the vm is databases.vm.netsoc.co, install MySQL
  * They target the groups we specified in the description of the VM
    * Don't forget you can target by name too, like `databases.vm.netsoc.co`

* We can manage our router setup
  * We have more VMs than we have public IP addresses, therefore we must make use of a router and a [NAT](https://www.reddit.com/r/explainlikeimfive/comments/1wqc30/eli5_how_does_nat_network_address_translation_work/) configuration
  * We do not use a hardware router, instead we use a virtualized router called [VyOS](https://www.vyos.io/)
  * We can map any port on any of our available IP addresses to any port on a VM
  * We can use VLANs to serperate infra and user server traffic
    * We currently only make (real) use of 1-2 VLANs
    * The router has a prescence on every VLAN as 10.0.x.1
    * Send traffic to 10.0.x.1 when you're on VLAN x and it'l get routed to it's destination i.e other VMs on the same VLAN or the internet
    * Currently in use VLANs:
      * 10
      * 20 
      * 30 - Internal infrastructure (the majority of VMs will be on this)
      * 40 - User infrastructure (not currently in use-yet)
      * 100 - DHCP for Packer
  * All of this is configured via Ansible
    * TODO(ocanty)-elaborate more on this, mloc too pls
  * See `configure-router.yml`
  * Watch these (excuse the cat meme, they'll be better than any lecture you'll ever get in college):
    * [A Cat Explains the Internet](https://www.youtube.com/watch?v=jHyaFBsxy1s)
    * [A Cat Explains Subnetting](https://www.youtube.com/watch?v=QgT1s2fOfiE)
    * [A Cat Explains DNS](https://www.youtube.com/watch?v=4ZtFk2dtqv0)

* We can manage our DNS records
  * All (external) DNS records are hosted on Cloudflare and completely managed by the a playbook
  * Internal DNS records are managed by FreeIPA
  * See `configure-dns.yml`

* All of this should tie into eachother
  * We use `vars/ip_allocation.yml` to set IP allocations for VMs
  * We can feed that to `create-*.yml` to create a VM that uses that IP
  * We can feed that to `configure-router.yml` and create a port-mapping that maps an external port to that VM's IP 
  * We can then set a DNS record to point at the IP that the external port is on

* All of this is stuck in time forever because git is awesome

### Important policies

* **All volatile data should be on a 2nd disk of the VM, not the boot disk**
  * Configure your VMs tolerate the boot disk being completely destroyed
    * i.e if we want to restore from a backup, we can delete the VM, run the playbook to create a new VM, detach the e,empty 2nd disk it created and attach the disk we backed up
  * By defining a VM entirely by playbooks we can potentially save a LOT of disk space on backups
    * Needs to be investigated

* **Do NOT put a VM on a trunk port unless it's a router**
  * Always give the VMs a NIC with a tagged VLAN
  * If you want a prescence on multiple VLANs, attach multiple NICs
  * If you give the VM with a NIC with no tagging, it will recieve ALL traffic.
      * This is a huge security risk

### Finally, our actual VM setup

**IP allocations for these VMs can be found in `vars/ip_allocation.yml` and are subject to change**

* auth.vm.netsoc.co
  * Hosts a FreeIPA server on NIC 1
    * Container hostname: 
      * ipa.vm.netsoc.co
    * LDAP for user accounts
    * DNS server that serves *.vm.netsoc.co
    * Kerberos
    * Web UI
  * Hosts Keycloak on NIC 2
    * Provides OAuth/OpenID Connect for users in LDAP

**Every single other VM bar the router should be enrolled in the FreeIPA server**

* nfs.vm.netsoc.co
  * Provides an authenticated Kerberos NFS server
  * Stores home directories and docker data
  * ZFS is used on the data disk
  * A script exists that periodically should scan FreeIPA and create zfs directories for users

* databases.vm.netsoc.co
  * Hosts MySQL database for infra on NIC 1
  * Hosts MYSQL databsae for users on NIC 2

* managerN.vm.netsoc.co, workerN.vm.netsoc.co
  * Docker swarm managers and workers
  * Hosts the following containers:
    * Traefik (this reverse proxies the **majority** of our web services. It is **vital**)

* router.vm.netsoc.co
  * The VyOS router
  * This is the only VM that should bind public IPs (when migrating from our old infra is done)
  * Notable mappings:
    * TODO

* portal.vm.netsoc.co
  * The user server