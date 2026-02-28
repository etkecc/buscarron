Hello,
Great news — your Matrix server order is confirmed! 🎉

We don't need any additional details right now — here are the next steps:

**Payment Instructions**:

1. Visit our [membership page](https://etke.cc/membership).
2. Select the "By Complexity" tier.
3. Set the custom price to **€14**.
4. Subscribe on Ko-fi with the same email address you used for this order (test@test.com).

Once your payment is confirmed, we'll start setting up your Matrix server right away. You'll receive a separate onboarding email with all the details you need.

Please make sure [all mandatory ports are open](https://etke.cc/order/status/#ports-and-firewalls), and that the following SSH key is added to **/home/matrix/.ssh/authorized_keys**: `ssh-todo TODO etke.cc`. Here's how:

1.  Ensure the .ssh directory exists: `mkdir -p /home/matrix/.ssh`
2.  Ensure the authorized_keys file exists: `touch /home/matrix/.ssh/authorized_keys`
3.  Add the key to the authorized_keys file: `echo 'ssh-todo TODO etke.cc' >> /home/matrix/.ssh/authorized_keys`
4.  Set the correct permissions on the authorized_keys file: `chmod 600 /home/matrix/.ssh/authorized_keys`
5.  Set the correct permissions on the .ssh directory: `chmod 700 /home/matrix/.ssh`
6.  Set the correct ownership on the .ssh directory: `chown -hR matrix:matrix /home/matrix/.ssh`

Please, add the following DNS entries:

* @ A record 1.2.3.4
* matrix A record 1.2.3.4
* matrix TXT record v=spf1 ip4:1.2.3.4 -all
* _dmarc.matrix TXT record v=DMARC1; p=quarantine;
* default._domainkey.matrix TXT record v=DKIM1; k=rsa; p=TODO

You can track your order status here: [Order Status Page](https://etke.cc/order/status/#a379a6f6eeafb9a55e378c118034e2751e682fab9f2d30ab13d2125586ce1947).

Questions? Feel free to [contact us](https://etke.cc/contacts/) — we're here to help.

We're excited to work with you and provide a top-notch Matrix server experience.

Warm regards,
the [etke.cc](https://etke.cc) team