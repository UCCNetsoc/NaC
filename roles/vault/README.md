Hashicorp Vault role
====================

Installs a single Hashicorp Vault in Swarm, using the existing Consul cluster as a backing store for HA storage.

Requires `netaddr` Python library to be installed locally.

# Convention

To keep a nice order to the various files that will hopefully accumulate here, we'll introduce a simple naming scheme and structure.

## Roles

In general, there are two kinds of roles that we'll be needing:

1. Secrets (mounting a secret engine and pushing the secrets)
2. Policies (creating/pushing a token for an associated policy and pushing the policy)

As such, secrets roles should be prefixed with `secrets_` and policy roles should be prefixed with `policy_`.

Then, in the top level playbook, we'll create a task for the new secrets role, and include any policy roles that have access to these secrets, eg:
Discord Webhook secrets can be referenced by the CI policy, therefore the CI policy role will be run followed by the Discord Webhook secrets role.

---

**NOTE**
Using fork of `consul_acl` until prefix policies are merged upstream. ([PR 1 Closed.](https://github.com/ansible/ansible/pull/64809)) ([PR 2 Open.](https://github.com/ansible/ansible/pull/62925))

**NOTE 2**
Until Swarm supports setting capabilities, mlock is disabled. ([GH Issue](https://github.com/moby/moby/issues/25885))  
If the hostnames IP does not start with 84, `mlock` is disabled. Haven't figured out how to get mlock to work in LXD yet.

Add the following to config.hcl once caps are supported in swarm. This will enable `mlock` in prod and disable it in the simulation:

```hcl
{% if inventory_hostname|ipaddr('1') != '84' %}
disable_mlock = true
{% endif %}
```
