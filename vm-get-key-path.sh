#!/bin/bash
hostname="$1.vm.netsoc.co"
echo `./proxmox_inventory.py --list | jq '._meta.hostvars["'$"$hostname"'"].ansible_ssh_private_key_file' | sed 's/"//g'`
