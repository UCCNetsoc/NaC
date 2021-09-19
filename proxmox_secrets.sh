#!/bin/bash
# Source this script extract ENV from the Ansible Vault for Proxmox secrets
# PM_API_URL, PM_USER, PM_PASS

# Wow! https://stackoverflow.com/a/21189044
parse_yaml() {
   local prefix=$2
   local s='[[:space:]]*' w='[a-zA-Z0-9_]*' fs=$(echo @|tr @ '\034')
    sed -ne "s|^\($s\)\($w\)$s:$s\"\(.*\)\"$s\$|\1$fs\2$fs\3|p" \
        -e "s|^\($s\)\($w\)$s:$s\(.*\)$s\$|\1$fs\2$fs\3|p"  $1 |
    awk -F$fs '$3 !~ /\(/{
    indent = length($1)/2;
    vname[indent] = $2;
    for (i in vname) {if (i > indent) {delete vname[i]}}
        if (length($3) > 0) {
            vn=""; for (i=0; i<indent; i++) {vn=(vn)(vname[i])("_")}
            printf("%s%s%s=\"%s\"\n", "'$prefix'",vn, $2, $3);
        }
    }'c
}

editor=$EDITOR
export EDITOR=cat

if [[ -z "${VAULT_PASS}" ]]; then
    echo -n "Vault password: "
    read -s vault_pass
    echo "$vault_pass" > ./_vault_pass
    export VAULT_PASS=$vault_pass
else
    echo "$VAULT_PASS" > ./_vault_pass
fi

echo ""
ansible-vault edit vars/secrets.yml --vault-password-file ./_vault_pass > ./_secrets.yml 
eval $(parse_yaml _secrets.yml "yaml_")
rm _secrets.yml
rm _vault_pass
export EDITOR=$editor

export PM_USER=$yaml_vault_proxmox_username
export PM_PASS=$yaml_vault_proxmox_password

for f in *.sh; do
	alias "${f::-3}"="$(pwd)/$f"
done

echo "All .sh scripts in NaC are global. EG: 'run provision-infra-web.yml' maps to './run.sh provision-infra-web.yml'"
