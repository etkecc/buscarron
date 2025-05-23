Hello,
We're thrilled to share that your Matrix server order is confirmed! 🎉

No need for additional details at this moment; we'll keep it simple:

**Payment Instructions**:

1. Visit our [membership page](https://etke.cc/membership).
2. Select the "By Complexity" tier.
3. Set the custom price to **$13**.
4. Subscribe on Ko-Fi with the same email address you used for this order (test@test.com).

Once your payment is confirmed, we'll promptly initiate the setup of your Matrix server. Look forward to a new email that will guide you through the onboarding process with all the necessary details.

Please, ensure [all mandatory ports are open](https://etke.cc/order/status/#ports-and-firewalls), and the following ssh key added to the **/home/user/.ssh/authorized_keys**: `ssh-todo TODO etke.cc`. Here is how you can do that:

1.  ensure the .ssh directory exists: `mkdir -p /home/user/.ssh`
2.  ensure the authorized_keys file exists: `touch /home/user/.ssh/authorized_keys`
3.  add the key to the authorized_keys file: `echo 'ssh-todo TODO etke.cc' >> /home/user/.ssh/authorized_keys`
4.  ensure the correct permissions are set on the authorized_keys file: `chmod 600 /home/user/.ssh/authorized_keys`
5.  ensure the correct permissions are set on the .ssh directory: `chmod 700 /home/user/.ssh`
6.  ensure the correct ownership is set on the .ssh directory: `chown -hR user:user /home/user/.ssh`

Please, add the following DNS entries:

* @    A record    1.2.3.4
* matrix    A record    1.2.3.4
* matrix    TXT record    v=spf1 ip4:1.2.3.4 -all
* _dmarc.matrix    TXT record    v=DMARC1; p=quarantine;
* default._domainkey.matrix    TXT record    v=DKIM1; k=rsa; p=TODO

To check the status of your order and stay updated, please keep an eye on your [Order Status Page](https://etke.cc/order/status/#a379a6f6eeafb9a55e378c118034e2751e682fab9f2d30ab13d2125586ce1947).

Got any questions? Feel free to [contact us](https://etke.cc/contacts/) - we're here to assist you!

We're genuinely excited to serve you and provide a top-notch Matrix server experience.

Warm regards,
the [etke.cc](https://etke.cc) team