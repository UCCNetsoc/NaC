#!/bin/bash
# Source this script extract ENV from the Ansible Vault for Proxmox secrets
# PM_API_URL, PM_USER, PM_PASS

# Wow! https://stackoverflow.com/a/21189044
parse_yaml() {
   local prefix=$2
   local s='[[:space:]]*' w='[a-zA-Z0-9_]*' fs=$(echo @|tr @ '\034')
    sed -ne "s|^\($s\)\($w\)$s:$s\"\(.*\)\"$s\$|\1$fs\2$fs\3|p" \
        -e "s|^\($s\)\($w\)$s:$s\(.*\)$s\$|\1$fs\2$fs\3|p"  $1 |
    awk -F$fs '{
    indent = length($1)/2;
    vname[indent] = $2;
    for (i in vname) {if (i > indent) {delete vname[i]}}
        if (length($3) > 0) {
            vn=""; for (i=0; i<indent; i++) {vn=(vn)(vname[i])("_")}
            printf("%s%s%s=\"%s\"\n", "'$prefix'",vn, $2, $3);
        }
    }'
}

export EDITOR=cat
ansible-vault edit vars/secrets.yml > ./_secrets.yml 
eval $(parse_yaml _secrets.yml "vault_")
rm _secrets.yml
export EDITOR=vi

export PM_API_URL="https://$vault_proxmox_host:8006/api2/json"
export PM_HOST=$vault_proxmox_host
export PM_USER=root@pam
export PM_PASS=$vault_proxmox_password

