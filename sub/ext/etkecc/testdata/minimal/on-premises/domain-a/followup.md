Hello,
We're thrilled to share that your Matrix server order is confirmed! 🎉

Server: please, create an x86/amd64 VPS with any Debian-based distro. Minimal comfortable configuration for a Matrix server without any additional components: 1vCPU, 2GB RAM.
Add our SSH keys ([etke.cc/keys.txt](https://etke.cc/keys.txt)) to your server, open the required ports ([etke.cc/help/faq#what-ports-should-be-open](https://etke.cc/help/faq#what-ports-should-be-open)) send us your server's IP address, the username (with permissions to call sudo), and password (if set).

Please, ensure [all mandatory ports are open](https://etke.cc/help/faq#what-ports-should-be-open).

Please, add the following DNS entries (ensure that the CloudFlare proxy is disabled!):

* matrix    A record    server IP

We see that you have something on your base domain.
**If** that's a domain registrar's (parking) page and/or you intend to serve base domain (example.com) from the matrix server, just add the `@` DNS record pointing to the server IP and tell us about that.
**If** that's your website and/or you intend to serve base domain from some other server, you should add the following HTTPS redirects (HTTP 301):

* [example.com/.well-known/matrix/server](https://example.com/.well-known/matrix/server) -> [matrix.example.com/.well-known/matrix/server](https://matrix.example.com/.well-known/matrix/server)
* [example.com/.well-known/matrix/client](https://example.com/.well-known/matrix/client) -> [matrix.example.com/.well-known/matrix/client](https://matrix.example.com/.well-known/matrix/client)
* [example.com/.well-known/matrix/support](https://example.com/.well-known/matrix/support) -> [matrix.example.com/.well-known/matrix/support](https://matrix.example.com/.well-known/matrix/support)
  To learn more about why these redirects are necessary and what the connection between the base domain (example.com) and the Matrix domain (matrix.example.com) is, read the following guide: [etke.cc/help/faq#why-are-well-known-redirects-on-the-base-domain-important](https://etke.cc/help/faq#why-are-well-known-redirects-on-the-base-domain-important)

Got any questions? Feel free to reply to this email - we're here to assist you!

We're genuinely excited to serve you and provide a top-notch Matrix server experience.

Best regards,

etke.cc