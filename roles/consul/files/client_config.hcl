datacenter = "cix-main"
primary_datacenter = "cix-main"
server = false
ui = false
node_name = "{{ inventory_hostname }}.consul-client.netsoc.dev"

leave_on_terminate = true
retry_join = ["consul"]

client_addr = "0.0.0.0"

acl {
    enabled = true
    default_policy = "deny"
    tokens {
        agent = "{{ consul_agent_token }}"
    }
}