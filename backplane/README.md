# Backplane Provisioning

All files in this directory correspond to setting up the Proxmox servers for VM creation & monitoring

## ISO downloads
ISOs are required one of the machines of the Proxmox cluster in order for Packer to build VM templates

ISOs are downloaded by the iso_download.yml playbook

### Ansible setup
See ```hosts``` and ```host_vars/feynman.yml``` for examples
Keys must be placed in ```/keys/<host>/id_rsa```

# Networking Assumptions
    insert mloc stuff here