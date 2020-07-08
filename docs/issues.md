# Common Issues


## `Could not open lock file /var/lib/dpkg/lock-frontend`
* `{"cache_update_time": 1594165492, "cache_updated": false, "changed": false, "msg": "'/usr/bin/apt-get -y -o \"Dpkg::Options::=--force-confdef\" -o \"Dpkg::Options::=--force-confold\"      install 'autofs'' failed: E: Could not open lock file /var/lib/dpkg/lock-frontend - open (13: Permission denied)\nE: Unable to acquire the dpkg frontend lock (/var/lib/dpkg/lock-frontend), are you root?\n", "rc": 100, "stderr": "E: Could not open lock file /var/lib/dpkg/lock-frontend - open (13: Permission denied)\nE: Unable to acquire the dpkg frontend lock (/var/lib/dpkg/lock-frontend), are you root?\n", "stderr_lines": ["E: Could not open lock file /var/lib/dpkg/lock-frontend - open (13: Permission denied)", "E: Unable to acquire the dpkg frontend lock (/var/lib/dpkg/lock-frontend), are you root?"], "stdout": "", "stdout_lines": []}
`
* `E: Could not get lock /var/lib/dpkg/lock-frontend - open (11: Resource temporarily unavailable)\nE: Unable to acquire the dpkg frontend lock (/var/lib/dpkg/lock-frontend), is another process using it?\n", "rc": 100, "stderr": "E: Could not get lock /var/lib/dpkg/lock-frontend - open (11: Resource temporarily unavailable)\nE: Unable to acquire the dpkg frontend lock (/var/lib/dpkg/lock-frontend), is another process using it?\n", "stderr_lines": ["E: Could not get lock /var/lib/dpkg/lock-frontend - open (11: Resource temporarily unavailable)", "E: Unable to acquire the dpkg frontend lock (/var/lib/dpkg/lock-frontend), is another process using it?"], "stdout": "", "stdout_lines": []}`

* This usually means the VM is running automated security updates, (enabled by default in Ubuntu). i.e dpkg is in use
* Leave it for awhile and run the playbook again
* It could also mean you aren't using `become: yes` to elevate to root to run `apt`/`dpkg`

## `ERROR! Unable to execute the command "/root/.ansible/tmp/ansible-local-215645qyJDZD/tmp8c26XD.yml": [Errno 13] Permission denied` when running `ansible-vault edit`

  * Your `$EDITOR` isn't set.
    * Run `export EDITOR=nano` / `export EDITOR=vim`

## `Access denied` when reading/writing/creating files as `root` on an nfsv4

* Make sure that the name of the Kerberos ticket you are using is statically mapped to `root` in `vars/idmap.yml`
* Make sure idmap.conf has GSS-Methods set correctly
* Check `/etc/idmapd.conf` on the host
* Example of a working configuration:
  * ```
    [General]

    Verbosity = 0
    Pipefs-Directory = /run/rpc_pipefs
    # set your own domain here, if it differs from FQDN minus hostname
    # Domain = localdomain

    Domain = vm.netsoc.co

    [Mapping]
    Nobody-User = nobody
    Nobody-Group = nogroup

    [Translation]
    GSS-Methods = static,nsswitch

    [Static]
    nfs/manager0.vm.netsoc.co@VM.NETSOC.CO = root
    nfs/manager1.vm.netsoc.co@VM.NETSOC.CO = root
    nfs/manager2.vm.netsoc.co@VM.NETSOC.CO = root
    nfs/manager3.vm.netsoc.co@VM.NETSOC.CO = root
    nfs/manager4.vm.netsoc.co@VM.NETSOC.CO = root
    nfs/worker0.vm.netsoc.co@VM.NETSOC.CO  = root
    nfs/worker1.vm.netsoc.co@VM.NETSOC.CO  = root
    nfs/worker2.vm.netsoc.co@VM.NETSOC.CO  = root
    nfs/worker3.vm.netsoc.co@VM.NETSOC.CO  = root
    ```
