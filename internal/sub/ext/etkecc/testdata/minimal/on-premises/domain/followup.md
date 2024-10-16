Hello,
We're thrilled to share that your Matrix server order is confirmed! 🎉

Server: please, create an x86/amd64 VPS with any Debian-based distro. Minimal comfortable configuration for a Matrix server without any additional components: 2vCPU, 2GB RAM.
Open the required ports ([etke.cc/order/status/#ports-and-firewalls](https://etke.cc/order/status/#ports-and-firewalls)) send us your server's IP address, the username (with permissions to call sudo), and password (if set).

Please, ensure [all mandatory ports are open](https://etke.cc/order/status/#ports-and-firewalls), and the following ssh key added to the **/home/TODO/.ssh/authorized_keys**: `ssh-todo TODO etke.cc`. Here is how you can do that:

```bash
# ensure the .ssh directory exists
mkdir -p /home/TODO/.ssh
# ensure the authorized_keys file exists
touch /home/TODO/.ssh/authorized_keys
# add the key to the authorized_keys file
echo 'ssh-todo TODO etke.cc' >> /home/TODO/.ssh/authorized_keys
# ensure the correct permissions are set on the authorized_keys file
chmod 600 /home/TODO/.ssh/authorized_keys
# ensure the correct permissions are set on the .ssh directory
chmod 700 /home/TODO/.ssh
# ensure the correct ownership is set on the .ssh directory
chown -hR TODO:TODO /home/TODO/.ssh
```

Please, add the following DNS entries:

* @    A record    server IP
* matrix    A record    server IP
* matrix    TXT record    v=spf1 ip4:server IP -all
* _dmarc.matrix    TXT record    v=DMARC1; p=quarantine;
* default._domainkey.matrix    TXT record    v=DKIM1; k=rsa; p=TODO

To check the status of your order and stay updated, please keep an eye on your [Order Status Page](https://etke.cc/order/status/#a379a6f6eeafb9a55e378c118034e2751e682fab9f2d30ab13d2125586ce1947).

Got any questions? Feel free to [contact us](https://etke.cc/contacts/) - we're here to assist you!

We're genuinely excited to serve you and provide a top-notch Matrix server experience.

Warm regards,
the [etke.cc](https://etke.cc) team