Hello,
We're thrilled to share that your Matrix server order is confirmed! 🎉

Server: please, create an x86/amd64 VPS with any Debian-based distro. Minimal comfortable configuration for a Matrix server without any additional components: 1vCPU, 2GB RAM.
Add our SSH keys ([etke.cc/keys.txt](https://etke.cc/keys.txt)) to your server, open the required ports ([etke.cc/help/faq#what-ports-should-be-open](https://etke.cc/help/faq#what-ports-should-be-open)) send us your server's IP address, the username (with permissions to call sudo), and password (if set).

Please, ensure [all mandatory ports are open](https://etke.cc/help/faq#what-ports-should-be-open).

Please, add the following DNS entries:

* @    A record    server IP
* matrix    A record    server IP
* matrix    TXT record    v=spf1 ip4:server IP -all
* _dmarc.matrix    TXT record    v=DMARC1; p=quarantine;
* default._domainkey.matrix    TXT record    v=DKIM1; k=rsa; p=MIIBojANBgkqhkiG9w0BAQEFAAOCAY8AMIIBigKCAYEAooAm4ZCkY9YxYN2ATlv2widBMf6hQavSCq3mzmxaw0FG6/X9zXtZR3h0OnJtPAFHZOt5uwNTqsUkwJk0yZqWqPDarSuN/+uvpwtN/qGN35TIIYmiEGiThOtVoesye8x1ZwTQCZxttL7GvzKOfBRVoJRPC3oM63/dRBcpf6GuTcdALF5GEGk+YDltb64PhyJ5TT57QFzgipgCjTwbugh/MrxwiXDpRRWPrSrSAOVKiITkvP2bedcXDliAiEVOwb+BrOUotiSiD3fr9Tb8bO33CDE62EvgXjFwXLHvDDCJBaI2Gqu2d9+jQX9ZRzWuNsc+HqO6nhlU+i025iZaawEcjC144ZExtqWm8gbqDfK2pJWKr0tyqxlj/Ujd7Mgat9YT4XZQr8YyKexB5GT26SYtj+fvRWxjZartLJbqs5PyluKgl2giz+MjKUvU4I1eeZsOuNG0fLK/FXJ20/1NIvO54Mt8mG3sioVW2UnrAkqmZLe6mQ5F8t5p9AzECEF2vRSfAgMBAAE=

Got any questions? Feel free to reply to this email - we're here to assist you!

We're genuinely excited to serve you and provide a top-notch Matrix server experience.

Best regards,

etke.cc