Hello,
Great news — your Matrix server order is confirmed! 🎉

Server: please, create an x86/amd64 VPS with any Debian-based distro. Minimal comfortable configuration for a Matrix server without any additional components: 2vCPU, 2GB RAM.
Open the required ports ([etke.cc/order/status/#ports-and-firewalls](https://etke.cc/order/status/#ports-and-firewalls)) send us your server's IP address, the username (with permissions to call sudo), and password (if set).

Please make sure [all mandatory ports are open](https://etke.cc/order/status/#ports-and-firewalls), and that the following SSH key is added to **/home/TODO/.ssh/authorized_keys**: `ssh-todo TODO etke.cc`. Here's how:

1.  Ensure the .ssh directory exists: `mkdir -p /home/TODO/.ssh`
2.  Ensure the authorized_keys file exists: `touch /home/TODO/.ssh/authorized_keys`
3.  Add the key to the authorized_keys file: `echo 'ssh-todo TODO etke.cc' >> /home/TODO/.ssh/authorized_keys`
4.  Set the correct permissions on the authorized_keys file: `chmod 600 /home/TODO/.ssh/authorized_keys`
5.  Set the correct permissions on the .ssh directory: `chmod 700 /home/TODO/.ssh`
6.  Set the correct ownership on the .ssh directory: `chown -hR TODO:TODO /home/TODO/.ssh`

Please, add the following DNS entries:

* @ A record server IP
* matrix A record server IP
* matrix TXT record v=spf1 ip4:server IP -all
* _dmarc.matrix TXT record v=DMARC1; p=quarantine;
* default._domainkey.matrix TXT record v=DKIM1; k=rsa; p=TODO

You can track your order status here: [Order Status Page](https://etke.cc/order/status/#a379a6f6eeafb9a55e378c118034e2751e682fab9f2d30ab13d2125586ce1947).

Questions? Feel free to [contact us](https://etke.cc/contacts/) — we're here to help.

We're excited to work with you and provide a top-notch Matrix server experience.

Warm regards,
the [etke.cc](https://etke.cc) team