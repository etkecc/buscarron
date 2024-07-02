Hello,
We're thrilled to share that your Matrix server order is confirmed! 🎉

Reminder bot: What's your timezone (IANA)? Like America/Chicago, Asia/Seoul, or Europe/Berlin. [Full list](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones#List)

Telegram: please, go to [https://my.telegram.org/apps](https://https://my.telegram.org/apps) and create a new app. Share the API ID and Hash with us

Website: to deploy a static website you have to point your base domain (the @ DNS entry) to Matrix server IP and the website source has to be available in a public git repo. Supported generators: hugo, jekyll, plain html (no generator). Are you sure you want it? If so, please, provide the website repository URL, command(-s) to build it, and in what folder the build dist is saved (usually public or dist).

SSO: You didn't mention what OIDC/OAuth2 provider you want to integrate, so here is a list of common providers - [github.com/matrix-org/synapse/blob/develop/docs/openid.md#sample-configs](https://github.com/matrix-org/synapse/blob/develop/docs/openid.md#sample-configs). Please, send us the information required to configure it (usually it's provider name, issuer, client_id, client_secret, but that depends on the provider)

Sygnal: are you sure you want it? It's a push gateway, usable only for Matrix client app developers, so you can't use it if you don't develop your mobile Matrix app. If you want to add it, please, provide the following information: app ID(-s) (eg org.matrix.app), FCM api key, and/or APNS certificate (if used)

BorgBackup: please, provide the desired repository url (user@host:repo). We will generate an SSH key and encryption passphrase on your server. We will send you the public part of the generated SSH key. You will need to add that SSH key to your provider.

SSH: You are ordering a hosted/managed server. We will set up and manage the server on your behalf. Still, you can get full SSH access to this server. **If** you wish to have SSH access to this server, send us your public SSH key and a list of IP addresses from which you wish to access it.

To check the status of your order and stay updated, please keep an eye on your [Order Status Page](https://etke.cc/order/status/#749f066f31d6e795088f154897aba00b72bdbf951e4d5721caa37ee9d6eb31d9).

Got any questions? Feel free to reply to this email - we're here to assist you!

We're genuinely excited to serve you and provide a top-notch Matrix server experience.

Best regards,

etke.cc