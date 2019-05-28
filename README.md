Netsoc Ansible
=========

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