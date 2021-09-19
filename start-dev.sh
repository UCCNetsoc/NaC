#!/bin/bash

if [ ! -d "./bin" ]; then
	virtualenv -p /usr/bin/python3 .
	source bin/activate
	pip install -r ./requirements.txt
	ansible-galaxy install -r requirements.yml
fi

git config --local merge.ansible-vault.driver "./vault-merge.sh %O %A %B %L %P"
git config --local merge.ansible-vault.name "Ansible Vault merge driver"

if [[ $# -eq 1 ]] ; then
    if [[ $1 == 'fish' ]] ; then
		fish -C "source bin/activate.fish; source bass.fish; bass source proxmox_secrets.sh"
    fi
    exit 0
fi

bash --init-file <(echo "source bin/activate; source proxmox_secrets.sh")