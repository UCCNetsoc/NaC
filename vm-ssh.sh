#!/bin/bash

hostname="$1.vm.netsoc.co"
ip=`./proxmox_inventory.py --list | jq '._meta.hostvars["'$"$hostname"'"].ansible_host' | sed 's/"//g'`
key=`./proxmox_inventory.py --list | jq '._meta.hostvars["'$"$hostname"'"].ansible_ssh_private_key_file' | sed 's/"//g'`
echo "Connectng to $ip using netsoc with $key"
shift
ssh netsoc@$ip -i $key -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o LogLevel=ERROR $@
