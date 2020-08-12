datacenter = "cix-main"
primary_datacenter = "cix-main"
server = true
ui = false

leave_on_terminate = false
bootstrap_expect = 4

client_addr = "0.0.0.0"

acl {
    enabled = true
    default_policy = "deny"
    tokens {
        master = "{{ consul.master_token }}"
        agent = "{{ consul.agent_token }}"
    }
}