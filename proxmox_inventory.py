#!/usr/bin/env python3
# This is an Ansible custom inventory script that pulls the inventory from 
# a Proxmox cluster, pulling hostvars and group from the VM description
# 
# Your VM descriptions should look something like this:
#
#   groups:
#       - user
#   host_vars:
#       description: "feynman-user"


import sys
import yaml
import os
import json
from optparse import OptionParser

from proxmoxer import ProxmoxAPI, ResourceException

if not 'PM_PASS' in os.environ:
    print("You need to run source ./proxmox_secrets.sh first")
    quit(-1)

host = "localhost"
if 'PM_HOST' in os.environ:
    host = os.environ['PM_HOST']

proxmox = ProxmoxAPI(host, user=os.environ['PM_USER'], password=os.environ['PM_PASS'], verify_ssl=False)

# If you're maintaining this script for future use
# Here is a good tutorial on Ansible custom inventories
# https://www.jeffgeerling.com/blog/creating-custom-dynamic-inventories-ansible

inventory = {
    '_meta': {
        'hostvars': {
        }
    }
}

for node in proxmox.nodes.get():
    for vm in proxmox.nodes(node['node']).qemu.get():
        # Get the config off the VM
        vm_config = proxmox.nodes(node['node']).qemu(vm['vmid']).get('config')

        # We need to discover it's local 10.x.x.x IP 
        # so that we can ssh into the machine via ansible
        ssh_ip = ''

        if vm['status'] == 'running':
            # Request the IPs that the VM is exposing via it's guest agent
            try:
                config = proxmox.nodes(node['node']).qemu(f"{vm['vmid']}/config").get()
                try:
                    macaddress = config["net0"].split("=")[1].split(',')[0]
                except Exception as e:
                    macaddress = None

                ifaces = proxmox.nodes(node['node']).qemu(vm['vmid']).agent.get('network-get-interfaces')
                for iface in ifaces['result']:
                    #print(iface)
                    if 'hardware-address' in iface and iface['hardware-address'].lower() == macaddress.lower():
                        
                        if 'ip-addresses' in iface and len(iface['ip-addresses']) > 0:
                            ssh_ip = iface['ip-addresses'][0]['ip-address']
                            break

            except ResourceException:
                # If the agent isn't running on the machine, ignore that machine
                continue
                pass
            

            inventory['_meta']['hostvars'][vm['name']] = {
                "ansible_host": ssh_ip
            }
            
            ssh_ip = ''

            # Parse yaml from vm description
            if 'description' in vm_config:
                desc_config = yaml.safe_load(vm_config['description'])

                if type(desc_config) == dict:
                    # Host-specific vars
                    if 'host_vars'in desc_config:
                        inventory['_meta']['hostvars'][vm['name']] = {
                            **inventory['_meta']['hostvars'][vm['name']],
                            **desc_config['host_vars']
                        }

                    # Groups
                    if 'groups' in desc_config:
                        for group in desc_config['groups']:
                            if group not in inventory:
                                inventory[group] = {
                                    'hosts': [ ],
                                    'host_vars': {
                                        # group specific vars would be in here if specific
                                     }
                                }

                            inventory[group]['hosts'].append(vm['name'])


if __name__ == '__main__':
    parser = OptionParser()
    parser.add_option('--list', action="store_true", help="Output JSON inventory")
    
    # We don't need to support the --host parameter as have _meta in our inventory
    # Ansible will not make calls with --host if it can get the full inventory
    # This is commented in case we want to implement this in the future
    # parser.add_option('--host', action="store")
    (options, args) = parser.parse_args()

    if options.list:
        print(json.dumps(inventory))
        quit(0)
