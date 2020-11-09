#!/bin/bash
editor=$EDITOR
export EDITOR=cat
if [[ ! -z "${VAULT_PASS}" ]]; then
    echo "$VAULT_PASS" > ./_vault_pass
    ansible-vault edit vars/secrets.yml --vault-password-file ./_vault_pass
    rm _vault_pass
else
    ansible-vault edit vars/secrets.yml
fi
export EDITOR=$editor
