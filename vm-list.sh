#!/bin/bash

./proxmox_inventory.py --list | jq '._meta.hostvars | with_entries( select(.key|contains("vm.netsoc.co") ) )'
