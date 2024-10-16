Hello,
We're thrilled to share that your Matrix server order is confirmed! 🎉

Server: please, create an x86/amd64 VPS with any Debian-based distro. Minimal comfortable configuration for a Matrix server without any additional components: 2vCPU, 2GB RAM.
Open the required ports ([etke.cc/order/status/#ports-and-firewalls](https://etke.cc/order/status/#ports-and-firewalls)) send us your server's IP address, the username (with permissions to call sudo), and password (if set).

Please, ensure [all mandatory ports are open](https://etke.cc/order/status/#ports-and-firewalls), and the following ssh key added to the **/home/TODO/.ssh/authorized_keys**: `ssh-todo TODO etke.cc`. Here is how you can do that:

1.  ensure the .ssh directory exists: `mkdir -p /home/TODO/.ssh`
2.  ensure the authorized_keys file exists: `touch /home/TODO/.ssh/authorized_keys`
3.  add the key to the authorized_keys file: `echo 'ssh-todo TODO etke.cc' >> /home/TODO/.ssh/authorized_keys`
4.  ensure the correct permissions are set on the authorized_keys file: `chmod 600 /home/TODO/.ssh/authorized_keys`
5.  ensure the correct permissions are set on the .ssh directory: `chmod 700 /home/TODO/.ssh`
6.  ensure the correct ownership is set on the .ssh directory: `chown -hR TODO:TODO /home/TODO/.ssh`

To check the status of your order and stay updated, please keep an eye on your [Order Status Page](https://etke.cc/order/status/#749f066f31d6e795088f154897aba00b72bdbf951e4d5721caa37ee9d6eb31d9).

Got any questions? Feel free to [contact us](https://etke.cc/contacts/) - we're here to assist you!

We're genuinely excited to serve you and provide a top-notch Matrix server experience.

Warm regards,
the [etke.cc](https://etke.cc) team