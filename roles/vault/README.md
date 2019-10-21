Hashicorp Vault role
====================

Installs a single Hashicorp Vault in Swarm, using the existing Consul cluster as a backing store for HA storage.

If the hostnames IP does not start with 84, mlock is disabled. Haven't figured out how to get mlock to work in LXD yet.