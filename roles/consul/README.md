Consul Swarm Cluster
===================

### How it works:

1. Deploys a bootstrap server instance
2. Tries to get its IP address (up to 50s/10 times)
3. Deploys two server replica leaders and a client (non leader)