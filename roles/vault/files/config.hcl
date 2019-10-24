ui = true
log_format = "json"
cluster_name = "cix-main"
{% if inventory_hostname|ipaddr('1') != '84' %}
disable_mlock = true
{% endif %}

listener "tcp" {
    address = "0.0.0.0:8200"
    tls_disable = "true"
}

storage "consul" {
    address = "consul:8500"
    path = "hashicorp-vault/"
    token = "{{ consul_vault_token }}"
}