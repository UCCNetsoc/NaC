#!/bin/bash

hostname="$1"
if [[ $hostname == "" ]] ; then
	arr=( $(./proxmox_inventory.py --list | jq -r '._meta.hostvars | keys | .[]') )
	select opt in "${arr[@]}"; do
		export hostname="$opt"
		break
	done
fi
ip=`./proxmox_inventory.py --list | jq '._meta.hostvars["'$"$hostname"'"].ansible_host' | sed 's/"//g'`
key=`./proxmox_inventory.py --list | jq '._meta.hostvars["'$"$hostname"'"].ansible_ssh_private_key_file' | sed 's/"//g'`
user=`./proxmox_inventory.py --list | jq '._meta.hostvars["'$"$hostname"'"].ansible_ssh_user' | sed 's/"//g'`
echo "Connectng to $ip using netsoc with $key"
shift
ssh $user@$ip -i $key -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o LogLevel=ERROR $@
