#!/bin/bash
printf "Run export EDITOR=<your choice of editor> to change the editor in use\n"
printf "If you get permission denied when running this script, you don't have an editor set\n"
ansible-vault edit vars/secrets.yml
