#!/bin/bash
ansible-playbook -i proxmox_inventory.py -i hosts ${@:1}