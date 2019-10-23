datacenter = "cix-main"
primary_datacenter = "cix-main"
server = true
ui = true

leave_on_terminate = true
retry_join = ["consul"]
bootstrap_expect = 3

client_addr = "0.0.0.0"

acl {
    enabled = true
    default_policy = "deny"
    tokens {
        master = "{{ consul_master_token }}"
        agent = "{{ consul_agent_token }}"
    }
}