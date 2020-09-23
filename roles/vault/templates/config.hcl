ui = true
log_format = "json"
cluster_name = "cix-main"
disable_mlock = true

api_addr = "{{ hashivault_url }}"

# 10 years.
max_lease_ttl = "87600h"
default_lease_ttl = "87600h"

listener "tcp" {
    address = "0.0.0.0:8200"
    tls_disable = 1
}

storage "mysql" {
  username = "{{ hashivault.mysql.user }}"
  password = "{{ hashivault.mysql.password }}"
  database = "{{ hashivault.mysql.database }}"
  address =  "{{ mysql_server }}:3306"
}