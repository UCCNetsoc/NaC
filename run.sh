#!/bin/bash
ansible-playbook -i proxmox_inventory.py -i hosts --ask-vault-pass ${@:1}
