Netsoc Ansible
=========

## Local Development

### Prerequisites

1. Hashicorp Packer
2. LXD up and running

### How 2

1. Run `./simulation.sh` and wait for everything to startup.
2. Change the IPs in `./host_vars/*` to the ones output in the console

The IPs only need to be updated once or after every time that you run `./simulation.sh teardown`

SSH Keys are generated and placed into the correct directories. You can SSH to each container/host using username `netsoc` and password `borger`

Make sure not to commit the host_vars with changed IPs.

## Misc Info

ssh priv key structure:
```/keys/<host>/id_rsa```


i.e for feynman
```/keys/feynman/id_rsa```

```/vars/cloudflare.yml``` format:
```
cf_api_email: <cloudflare email>
cf_api_key: <cloudflare api key>
```

```/vars/freeipa.yml``` format:
```
ipa_ds_password: <directory server pass>
ipa_admin_password: <admin pass>
```
