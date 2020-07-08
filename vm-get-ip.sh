#!/bin/bash

hostname="$1.vm.netsoc.co"
echo `./proxmox_inventory.py --list | jq '._meta.hostvars["'$"$hostname"'"].ansible_host' | sed 's/"//g'`

