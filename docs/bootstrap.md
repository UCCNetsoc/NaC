# Bootstrap Order

## How to build UCC Netsoc from scratch

# Production

1. Install Proxmox
1. Setup Networking to give the master a public IP/network connection
1. Modify `setup-proxmox-host-initial.yml` to match your NIC names
1. `setup-proxmox-host-initial.yml`
1. `setup-control-host.yml`
1. `template-netsoc-vyos-router.yml` in Packer (for one proxmox host)
1. `create-router.yml`
1. `provision-router.yml`
1. Join other hosts to the Cluster (and add to `./hosts`)
1. Run `setup-proxmox-host-initial.yml` on each host joined
1. `template-netsoc-ubuntu-server.yml`
1. `create-auth.yml`
1. `provision-auth.yml`
1. `create-nfs.yml`
1. `provision-nfs.yml`
1. `create-databases.yml`
1. `provision-databases.yml`
1. `create-games.yml`
1. `provision-games.yml`
1. `create-docker-swarm.yml`
1. `provision-docker-swarm.yml`
1. `create-web.yml`
1. `provision-web.yml`
1. `create-portal.yml`
1. `provision-portal.yml`
1. `setup-external-dns.yml`

# Development

Coming Soon