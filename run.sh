#!/bin/bash
if [[ ! -z "${VAULT_PASS}" ]]; then
    echo "$VAULT_PASS" > ./_vault_pass
    ansible-playbook -i proxmox_inventory.py -i hosts --vault-password-file ./_vault_pass ${@:1} || true
    rm _vault_pass
else
    ansible-playbook -i proxmox_inventory.py -i hosts --ask-vault-pass ${@:1}
fi