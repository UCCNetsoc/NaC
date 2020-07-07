#!/bin/bash

if [ ! -d "./bin" ]; then
	virtualenv -p /usr/bin/python2.7 .
	source bin/activate
	pip install requirements.txt
fi

bash --init-file <(echo "source bin/activate; source proxmox_secrets.sh")
