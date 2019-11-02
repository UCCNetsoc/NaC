ui = true
log_format = "json"
cluster_name = "cix-main"
disable_mlock = true

listener "tcp" {
    address = "0.0.0.0:8200"
    tls_disable = 1
}

storage "consul" {
    address = "consul_client:8500"
    path = "hashicorp-vault/"
    token = "{{ consul_vault_token }}"
}