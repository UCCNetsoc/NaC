Hashicorp Vault role
====================

Installs a single Hashicorp Vault in Swarm, using the existing Consul cluster as a backing store for HA storage.

Requires `netaddr` Python library to be installed locally.

If the hostnames IP does not start with 84, mlock is disabled. Haven't figured out how to get mlock to work in LXD yet.

**NOTE**
Using fork of `consul_acl` until prefix policies are merged upstream.

**NOTE 2**
Until Swarm supports setting capabilities, mlock is disabled.
Add the following to config.hcl once caps are supported:

```hcl
{% if inventory_#hostname|ipaddr('1') != '84' %}
disable_mlock = true
{% endif %}
```
