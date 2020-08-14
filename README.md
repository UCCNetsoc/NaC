# Netsoc Ansible

## I'm a new SysAdmin, what the heck do I do?

* `ssh` into the Ansible control server (currently `feynman.netsoc.co:2222`)
  * `ssh <username>@feynman.netsoc.co -p 2222 -i <path to ssh key>`
  * If you have not supplied an SSH key to the Head SysAdmin already:
    * Open a PR adding your username and key to `setup-control-host.yml`

* Clone this repo

* Run `./start-dev.sh` inside the cloned folder
  * You will need to run `./start-dev.sh` to setup the correct Python packages and environment variables. You must do this before beginning any development/deployment
  * You will be able to tell you have done this when your terminal prompt looks like this:
    * `(NaC) <user>@feynman:~/NaC#`

* You will need a `keys/` folder which contains SSH keys to target all physical and virtual machines. Ask the Head SysAdmin for this.
  * Do **NOT** commit them or remove the `keys/` clause from `.gitignore`.
  * Do **NOT** share them with people who are not SysAdmins
  * Do **NOT** leave them sitting on a random server somewhere, it's happened in the past

* You can peek and edit the vault using `./vault-peek.sh` and `./vault-edit.sh`

* You can list *.vm.netsoc.co by using `./vm-list.sh`
* You can ssh into *.vm.netsoc.co by running `./vm-ssh.sh <hostname>`, i.e `./vm-ssh.sh nfs` for `nfs.vm.netsoc.co`

* The Proxmox Web UI is available at [`feynman.netsoc.co:8006`](https://feynman.netsoc.co:8006). You may need to type `thisisunsafe` (if using Chrome) to get past the SSL warning 

* Once you've gotten to grips with that, have a look at `docs/` especially `infra.md`

* For your development, you can use `sshfs` / VSCode Remote / `vim` on the control server / a git branch.
  * You _will_ need to run your playbooks on the control server


## I want to contribute but I'm not a SysAdmin?

* Consider making an issue or contact us in `#servers` in our [Discord](https://discord.netsoc.co)
  * We'll welcome any help!

## I like this repo and want to learn more about UCC Netsoc

* Check out our [wiki](https://wiki.netsoc.co)

## **Important**

This repo currently contains both playbooks for managing 2019/2020 bare-metal infra and 2020/2021 Proxmox infra. Do not get them confused, have a look at `./hosts` to see what's going on.

### **Read [wiki.netsoc.co](wiki.netsoc.co)**
