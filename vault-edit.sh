#!/bin/bash
printf "Run export EDITOR=<your choice of editor> to change the editor in use\n"
printf "If you get permission denied when running this script, you don't have an editor set\n"
if [[ ! -z "${VAULT_PASS}" ]]; then
    echo "$VAULT_PASS" > ./_vault_pass
    ansible-vault edit vars/secrets.yml --vault-password-file ./_vault_pass
    rm _vault_pass
else
    ansible-vault edit vars/secrets.yml
fi
