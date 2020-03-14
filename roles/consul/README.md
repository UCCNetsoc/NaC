Consul Swarm Cluster
===================

## How it works:

Deploys a server instance to each of the manager nodes along with a bootstrap instance.
This bootstrap instance is where services should register themselves with. (IN TESTING STAGE)

There is a fault tolerance of 1 node before quorum is lost. Dont know what a quorum is? Read [this](https://www.consul.io/docs/internals/consensus.html) and make sure you understand how it works. Pay particular attention to the *Deployment Table*.

Adapted from [here](https://github.com/hashicorp/docker-consul/issues/66#issuecomment-353696910). Bless this man.