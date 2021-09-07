#!/bin/bash

if [ ! -d "./bin" ]; then
	virtualenv -p /usr/bin/python3 .
	source bin/activate
	pip install -r ./requirements.txt
	ansible-galaxy install -r requirements.yml
fi

if [[ $# -eq 1 ]] ; then
    if [[ $1 == 'zsh' ]] ; then
		zsh -c 'source bin/activate; source proxmox_secrets.sh; zsh -i'
    fi
    exit 0
fi

git config --local merge.ansible-vault.driver "./vault-merge.sh %O %A %B %L %P"
git config --local merge.ansible-vault.name "Ansible Vault merge driver"

bash --init-file <(echo "source bin/activate; source proxmox_secrets.sh")
