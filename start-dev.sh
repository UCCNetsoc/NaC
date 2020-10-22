#!/bin/bash

if [ ! -d "./bin" ]; then
	virtualenv -p /usr/bin/python2.7 .
	source bin/activate
	pip install -r ./requirements.txt
	ansible-galaxy install -r requirements.yml
	ansible-galaxy collection install -r requirements.yml
fi

if [[ $# -eq 1 ]] ; then
    if [[ $1 == 'zsh' ]] ; then
		zsh -c 'source bin/activate; source proxmox_secrets.sh; zsh -i'
    fi
    exit 0
fi

bash --init-file <(echo "source bin/activate; source proxmox_secrets.sh")
