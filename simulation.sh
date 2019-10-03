#!/bin/bash

# This spins up a simulation of the Netsoc infrastructure (just the 4 main servers so far) using (currently) LXD

SERVERS=(leela bigbertha boole lovelace)

if [[ $1 == 'teardown' ]]; then
    for server in ${SERVERS[@]}; do
        ip=$(lxc list --format json $server | jq -r .[0].state.network.eth0.addresses[0].address)
        ssh-keygen -R $ip > /dev/null
        printf '.'
        lxc stop $server
        printf '.'
        lxc rm $server
    done
    echo 'Done!'
    exit 0
fi


if [[ $(lxc image list --format json | jq -r '.[] | select(.aliases[0].name == "netsoc_infra").aliases[0].name') != 'netsoc_infra' ]]; then
    echo 'LXD image aliased "netsoc_infra" not found. Please build the LXD image using Packer and the packer.json file in ./packer'
    exit 1
fi

for server in ${SERVERS[@]}; do
    lxc launch netsoc_infra $server -c security.nesting=true
done

echo 'Sleeping to give containers time to obtain ipv4 address...'
sleep 5

IPS=()

echo ''
echo "Change the IP in ./host_vars to the following IPs. Don't forget to not commit those changes!"
for server in ${SERVERS[@]}; do
    ip=$(lxc list --format json $server | jq -r .[0].state.network.eth0.addresses[0].address)
    echo "$server: $ip"
    ssh-keyscan -H $ip >> ~/.ssh/known_hosts
    IPS+=($ip)
done

for i in "${!SERVERS[@]}"; do 
    server=${SERVERS[$i]}
    ip=${IPS[$i]}

    mkdir -p ./keys/$server

    echo -e 'y\n' | ssh-keygen -b 2048 -t rsa -N '' -f ./keys/$server/id_rsa > /dev/null
    echo borger | sshpass ssh-copy-id -i ./keys/$server/id_rsa netsoc@$ip > /dev/null
done