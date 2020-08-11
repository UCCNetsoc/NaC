datacenter = "cix-main"
primary_datacenter = "cix-main"
server = true
ui = true

leave_on_terminate = false
retry_join = ["consul_bootstrap"]
bootstrap_expect = 1

client_addr = "0.0.0.0"

acl {
    enabled = true
    default_policy = "deny"
    tokens {
        master = "{{ consul.master_token }}"
        agent = "{{ consul.agent_token }}"
    }
}