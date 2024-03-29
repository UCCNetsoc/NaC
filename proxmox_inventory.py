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

FILTERED_GROUPS = {"netsoc_cloud_vps", "netsoc_cloud_instance"}

CONTAINER_SUFFIX = ".infra.netsoc.co"

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
    for con in proxmox.nodes(node['node']).lxc.get():
        if not con['name'].endswith(CONTAINER_SUFFIX):
            continue

        # Get the config off the container
        vm_config = proxmox.nodes(node['node']).lxc(con['vmid']).get('config')

        # We need to discover it's local 10.x.x.x IP 
        # so that we can ssh into the machine via ansible
        ssh_ip = ''

        if con['status'] == 'running':
            # Request the IPs that the VM is exposing via it's guest agent
            try:
                config = proxmox.nodes(node['node']).lxc(f"{con['vmid']}/config").get()
                try:
                    macaddress = config["net0"].split("=")[5].split(',')[0]
                except Exception as e:
                    macaddress = None

                try:
                    ssh_ip = config["net0"].split("=")[6].split(',')[0]
                    ssh_ip = ssh_ip[:-3]
                except Exception as e:
                    continue

            except ResourceException:
                # If the agent isn't running on the machine, ignore that machine
                continue
                pass 


            inventory['_meta']['hostvars'][con['name']] = {
                "ansible_host": ssh_ip,
                "ansible_ssh_user": "root"
            }
            
            ssh_ip = ''

            # Parse yaml from vm description
            if 'description' in vm_config:
                desc_config = yaml.safe_load(vm_config['description'])

                if type(desc_config) == dict:
                    # Groups
                    filtered = False
                    if 'groups' in desc_config:
                        for group in desc_config['groups']:
                            if group in FILTERED_GROUPS:
                                filtered = True
                                continue
                            if group not in inventory:
                                inventory[group] = {
                                    'hosts': [ ],
                                    'host_vars': {
                                        # group specific vars would be in here if specific
                                     }
                                }
                            inventory[group]['hosts'].append(con['name'])

                    # Host-specific vars
                    if filtered:
                        del inventory['_meta']['hostvars'][con['name']]
                    elif 'host_vars'in desc_config:
                        inventory['_meta']['hostvars'][con['name']] = {
                            **inventory['_meta']['hostvars'][con['name']],
                            **desc_config['host_vars']
                        }

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
                "ansible_host": ssh_ip,
                "ansible_ssh_user": "netsoc"
            }
            
            ssh_ip = ''

            # Parse yaml from vm description
            if 'description' in vm_config:
                desc_config = yaml.safe_load(vm_config['description'])

                if type(desc_config) == dict:
                    # Groups
                    filtered = False
                    if 'groups' in desc_config:
                        for group in desc_config['groups']:
                            if group in FILTERED_GROUPS:
                                filtered = True
                                continue
                            if group not in inventory:
                                inventory[group] = {
                                    'hosts': [ ],
                                    'host_vars': {
                                        # group specific vars would be in here if specific
                                     }
                                }
                            inventory[group]['hosts'].append(vm['name'])

                    # Host-specific vars
                    if filtered:
                        del inventory['_meta']['hostvars'][vm['name']]
                    elif 'host_vars'in desc_config:
                        inventory['_meta']['hostvars'][vm['name']] = {
                            **inventory['_meta']['hostvars'][vm['name']],
                            **desc_config['host_vars']
                        }


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
