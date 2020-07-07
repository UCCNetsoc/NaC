Netsoc Ansible
=========

## Local Development

### Prerequisites

1. Hashicorp Packer
2. LXD up and running

### How 2

1. Run `./simulation.sh` and wait for everything to startup.
2. Go ham and run some playbooks!

The simulation may need to be torn down and brought back up after PC restarts.

SSH Keys are generated and placed into the correct directories. You can SSH to each container/host using username `netsoc` and password `borger`

Make sure not to commit the host_vars with changed IPs.

## Misc Info

ssh priv key structure:
`/keys/<host>/id_rsa`

i.e for feynman
`/keys/feynman/id_rsa`

`/vars/cloudflare.yml` format:

```yaml
cf_api_email: <cloudflare email>
cf_api_key: <cloudflare api key>
```

`/vars/freeipa.yml` format:

```yaml
ipa_ds_password: <directory server pass>
ipa_admin_password: <admin pass>
```
