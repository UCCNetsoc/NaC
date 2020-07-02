useradd -s /bin/bash -G sudo netsoc
echo "netsoc:borger" | chpasswd
mkdir -p /home/netsoc/.ssh
chown -R netsoc:netsoc /home/netsoc

echo 'netsoc  ALL=(ALL) NOPASSWD:ALL' | sudo EDITOR='tee -a' visudo

sed -i 's/PasswordAuthentication no/PasswordAuthentication yes/' /etc/ssh/sshd_config