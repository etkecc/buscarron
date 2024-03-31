package etkecc

var defaults = map[string]string{
	"domain":   "TODO.com",
	"username": "TODO",
	"email":    "email@TODO.com",
}

var dnsmap = map[string]string{
	"buscarron":    "buscarron",
	"cinny":        "cinny",
	"element-web":  "element",
	"etherpad":     "etherpad",
	"firezone":     "firezone",
	"funkwhale":    "funkwhale",
	"gotosocial":   "social",
	"hydrogen":     "hydrogen",
	"jitsi":        "jitsi",
	"languagetool": "languagetool",
	"linkding":     "linkding",
	"miniflux":     "miniflux",
	"ntfy":         "ntfy",
	"radicale":     "radicale",
	"schildichat":  "schildichat",
	"stats":        "stats",
	"sygnal":       "sygnal",
	"uptime-kuma":  "kuma",
	"vaultwarden":  "vault",
}

var botmap = map[string]string{
	"buscarron":    "@buscarron",
	"chatgpt":      "@chatgpt",
	"honoroit":     "@honoroit",
	"reminder-bot": "@reminder",
}

var bridgemap = map[string]string{
	"discord":    "@discordbot",
	"facebook":   "@messengerbot",
	"gmessages":  "@gmessagesbot",
	"googlechat": "@googlechatbot",
	"groupme":    "@_groupmepuppet_bot",
	"instagram":  "@instagrambot",
	"irc":        "@heisenbridge",
	"linkedin":   "@linkedinbot",
	"postmoogle": "@emailbot",
	"signal":     "@signalbot",
	"skype":      "@skypebridgebot",
	"slack":      "@slackbot",
	"steam":      "@_steampuppet_bot",
	"telegram":   "@telegrambot",
	"twitter":    "@twitterbot",
	"webhooks":   "@hookshot",
	"whatsapp":   "@whatsappbot",
}

var helpmap = map[string]string{
	"buscarron":    "etke.cc/help/bots/buscarron",
	"chatgpt":      "etke.cc/help/bots/chatgpt",
	"cinny":        "",
	"discord":      "etke.cc/help/bridges/mautrix-discord",
	"element-web":  "",
	"etherpad":     "",
	"facebook":     "etke.cc/help/bridges/mautrix-meta-messenger",
	"firezone":     "",
	"funkwhale":    "etke.cc/help/extras/funkwhale",
	"gmessages":    "etke.cc/help/bridges/mautrix-gmessages",
	"googlechat":   "etke.cc/help/bridges/mautrix-googlechat",
	"gotosocial":   "",
	"groupme":      "etke.cc/help/bridges/mx-puppet-groupme",
	"honoroit":     "etke.cc/help/bots/honoroit",
	"hydrogen":     "",
	"instagram":    "etke.cc/help/bridges/mautrix-meta-instagram",
	"irc":          "etke.cc/help/bridges/heisenbridge",
	"jitsi":        "etke.cc/help/extras/jitsi",
	"languagetool": "etke.cc/help/extras/languagetool",
	"linkding":     "",
	"linkedin":     "etke.cc/help/bridges/beeper-linkedin",
	"miniflux":     "etke.cc/help/extras/miniflux",
	"ntfy":         "etke.cc/help/extras/ntfy",
	"postmoogle":   "etke.cc/help/bridges/postmoogle",
	"radicale":     "",
	"reminder-bot": "etke.cc/help/bots/reminder",
	"schildichat":  "",
	"signal":       "etke.cc/help/bridges/mautrix-signal",
	"skype":        "etke.cc/help/bridges/go-skype-bridge",
	"slack":        "etke.cc/help/bridges/mautrix-slack",
	"stats":        "",
	"steam":        "etke.cc/help/bridges/mx-puppet-steam",
	"sygnal":       "etke.cc/help/extras/sygnal",
	"telegram":     "etke.cc/help/bridges/mautrix-telegram",
	"twitter":      "etke.cc/help/bridges/mautrix-twitter",
	"uptime-kuma":  "etke.cc/help/extras/uptime-kuma",
	"vaultwarden":  "etke.cc/help/extras/vaultwarden",
	"webhooks":     "etke.cc/help/bridges/hookshot",
	"whatsapp":     "etke.cc/help/bridges/mautrix-whatsapp",
}

// excluded by purpose: apple, dex, keycloak
var oidcmap = map[string]string{
	"auth0": "" +
		"  - idp_id: %s\n" +
		"    idp_name: %q\n" +
		"    idp_brand: %q\n" +
		"    issuer: %q\n" +
		"    client_id: %q\n" +
		"    client_secret: %q\n" +
		"    scopes: [\"openid\", \"profile\"]\n" +
		"    allow_existing_users: yes\n" +
		"    user_mapping_provider:\n" +
		"      config:\n" +
		"        localpart_template: \"{%% raw %%}{{ user.preferred_username|lower }}{%% endraw %%}\"\n" +
		"        display_name_template: \"{%% raw %%}{{ user.name|capitalize }}{%% endraw %%}\"\n",
	"authentik": "" +
		"  - idp_id: %s\n" +
		"    idp_name: %q\n" +
		"    idp_brand: %q\n" +
		"    issuer: %q\n" +
		"    discover: yes\n" +
		"    client_id: %q\n" +
		"    client_secret: %q\n" +
		"    scopes: [\"openid\", \"profile\", \"email\"]\n" +
		"    allow_existing_users: yes\n" +
		"    user_mapping_provider:\n" +
		"      config:\n" +
		"        localpart_template: \"{%% raw %%}{{ user.preferred_username|lower }}{%% endraw %%}\"\n" +
		"        display_name_template: \"{%% raw %%}{{ user.preferred_username|capitalize }}{%% endraw %%}\"\n",
	"django": "" +
		"  - idp_id: %s\n" +
		"    idp_name: %q\n" +
		"    idp_brand: %q\n" +
		"    issuer: %q\n" +
		"    client_id: %q\n" +
		"    client_secret: %q\n" +
		"    scopes: [\"openid\"]\n" +
		"    user_profile_method: \"userinfo_endpoint\"\n" +
		"    allow_existing_users: yes\n" +
		"    user_mapping_provider:\n" +
		"      config:\n" +
		"        localpart_template: \"{%% raw %%}{{ user.email.split('@')[0] }}{%% endraw %%}\"\n" +
		"        display_name_template: \"{%% raw %%}{{ user.first_name }} {{ user.last_name }}{%% endraw %%}\"\n" +
		"        email_template: \"{%% raw %%}{{ user.email }}{%% endraw %%}\"\n",
	"facebook": "" +
		"  - idp_id: %s\n" +
		"    idp_name: %q\n" +
		"    idp_brand: %q\n" +
		"    issuer: %q\n" +
		"    discover: no\n" +
		"    client_id: %q\n" +
		"    client_secret: %q\n" +
		"    authorization_endpoint: \"https://facebook.com/dialog/oauth\"\n" +
		"    token_endpoint: \"https://graph.facebook.com/v9.0/oauth/access_token\"\n" +
		"    jwks_uri: \"https://www.facebook.com/.well-known/oauth/openid/jwks/\"\n" +
		"    scopes: [\"openid\", \"email\"]\n" +
		"    allow_existing_users: yes\n" +
		"    user_mapping_provider:\n" +
		"      config:\n" +
		"        localpart_template: \"{%% raw %%}{{ user.name|lower }}{%% endraw %%}\"\n" +
		"        display_name_template: \"{%% raw %%}{{ user.name|capitalize }}{%% endraw %%}\"\n" +
		"        email_template: \"{%% raw %%}{{ user.email }}{%% endraw %%}\"\n",
	"gitea": "" +
		"  - idp_id: %s\n" +
		"    idp_name: %q\n" +
		"    idp_brand: %q\n" +
		"    issuer: %q\n" +
		"    discover: no\n" +
		"    client_id: %q\n" +
		"    client_secret: %q\n" +
		"    client_auth_method: \"client_secret_post\"\n" +
		"    authorization_endpoint: \"https://example.com/login/oauth/authorize\"\n" +
		"    token_endpoint: \"https://example.com/login/oauth/access_token\"\n" +
		"    userinfo_endpoint: \"https://example.com/api/v1/user\"\n" +
		"    scopes: []\n" +
		"    allow_existing_users: yes\n" +
		"    user_mapping_provider:\n" +
		"      config:\n" +
		"        subject_claim: \"id\"\n" +
		"        localpart_template: \"{%% raw %%}{{ user.login }}{%% endraw %%}\"\n" +
		"        display_name_template: \"{%% raw %%}{{ user.full_name }}{%% endraw %%}\"\n",
	"github": "" +
		"  - idp_id: %s\n" +
		"    idp_name: %q\n" +
		"    idp_brand: %q\n" +
		"    issuer: %q\n" +
		"    discover: no\n" +
		"    client_id: %q\n" +
		"    client_secret: %q\n" +
		"    authorization_endpoint: \"https://github.com/login/oauth/authorize\"\n" +
		"    token_endpoint: \"https://github.com/login/oauth/access_token\"\n" +
		"    userinfo_endpoint: \"https://api.github.com/user\"\n" +
		"    scopes: [\"read:user\"]\n" +
		"    allow_existing_users: yes\n" +
		"    user_mapping_provider:\n" +
		"      config:\n" +
		"        subject_claim: \"id\"\n" +
		"        localpart_template: \"{%% raw %%}{{ user.login }}{%% endraw %%}\"\n" +
		"        display_name_template: \"{%% raw %%}{{ user.name }}{%% endraw %%}\"\n",
	"gitlab": "" +
		"  - idp_id: %s\n" +
		"    idp_name: %q\n" +
		"    idp_brand: %q\n" +
		"    issuer: %q\n" +
		"    client_id: %q\n" +
		"    client_secret: %q\n" +
		"    client_auth_method: \"client_secret_post\"\n" +
		"    scopes: [\"openid\", \"read_user\"]\n" +
		"    user_profile_method: \"userinfo_endpoint\"\n" +
		"    allow_existing_users: yes\n" +
		"    user_mapping_provider:\n" +
		"      config:\n" +
		"        localpart_template: \"{%% raw %%}{{ user.nickname|lower }}{%% endraw %%}\"\n" +
		"        display_name_template: \"{%% raw %%}{{ user.name|capitalize }}{%% endraw %%}\"\n",
	"google": "" +
		"  - idp_id: %s\n" +
		"    idp_name: %q\n" +
		"    idp_brand: %q\n" +
		"    issuer: %q\n" +
		"    client_id: %q\n" +
		"    client_secret: %q\n" +
		"    scopes: [\"openid\", \"profile\", \"email\"]\n" +
		"    allow_existing_users: yes\n" +
		"    user_mapping_provider:\n" +
		"      config:\n" +
		"        localpart_template: \"{%% raw %%}{{ user.given_name|lower }}{%% endraw %%}\"\n" +
		"        display_name_template: \"{%% raw %%}{{ user.name|capitalize }}{%% endraw %%}\"\n" +
		"        email_template: \"{%% raw %%}{{ user.email }}{%% endraw %%}\"\n",
	"lemonldap": "" +
		"  - idp_id: %s\n" +
		"    idp_name: %q\n" +
		"    idp_brand: %q\n" +
		"    issuer: %q\n" +
		"    discover: yes\n" +
		"    client_id: %q\n" +
		"    client_secret: %q\n" +
		"    scopes: [\"openid\", \"profile\", \"email\"]\n" +
		"    allow_existing_users: yes\n" +
		"    user_mapping_provider:\n" +
		"      config:\n" +
		"        localpart_template: \"{%% raw %%}{{ user.preferred_username|lower }}{%% endraw %%}\"\n" +
		"        display_name_template: \"{%% raw %%}{{ user.preferred_username|capitalize }}{%% endraw %%}\"\n",
	"mastodon": "" +
		"  - idp_id: %s\n" +
		"    idp_name: %q\n" +
		"    idp_brand: %q\n" +
		"    issuer: %q\n" +
		"    discover: no\n" +
		"    client_id: %q\n" +
		"    client_secret: %q\n" +
		"    authorization_endpoint: \"https://example.com/oauth/authorize\"\n" +
		"    token_endpoint: \"https://example.com/oauth/token\"\n" +
		"    userinfo_endpoint: \"https://example.com/api/v1/accounts/verify_credentials\"\n" +
		"    scopes: [\"read\"]\n" +
		"    allow_existing_users: yes\n" +
		"    user_mapping_provider:\n" +
		"      config:\n" +
		"        subject_template: \"{%% raw %%}{{ user.id }}{%% endraw %%}\"\n" +
		"        localpart_template: \"{%% raw %%}{{ user.username }}{%% endraw %%}\"\n" +
		"        display_name_template: \"{%% raw %%}{{ user.display_name }}{%% endraw %%}\"\n",
	"microsoft": "" +
		"  - idp_id: %s\n" +
		"    idp_name: %q\n" +
		"    idp_brand: %q\n" +
		"    issuer: %q\n" +
		"    client_id: %q\n" +
		"    client_secret: %q\n" +
		"    scopes: [\"openid\", \"profile\"]\n" +
		"    authorization_endpoint: \"https://login.microsoftonline.com/example.com/oauth2/v2.0/authorize\"\n" +
		"    token_endpoint: \"https://login.microsoftonline.com/example.com/oauth2/v2.0/token\"\n" +
		"    userinfo_endpoint: \"https://graph.microsoft.com/oidc/userinfo\"\n" +
		"    allow_existing_users: yes\n" +
		"    user_mapping_provider:\n" +
		"      config:\n" +
		"        localpart_template: \"{%% raw %%}{{ user.preferred_username.split('@')[0] }}{%% endraw %%}\"\n" +
		"        display_name_template: \"{%% raw %%}{{ user.name|capitalize }}{%% endraw %%}\"\n",
	"shibboleth": "" +
		"  - idp_id: %s\n" +
		"    idp_name: %q\n" +
		"    idp_brand: %q\n" +
		"    issuer: %q\n" +
		"    discover: yes\n" +
		"    client_id: %q\n" +
		"    client_secret: %q\n" +
		"    scopes: [\"openid\", \"profile\", \"email\"]\n" +
		"    user_profile_method: \"userinfo_endpoint\"\n" +
		"    allow_existing_users: yes\n" +
		"    user_mapping_provider:\n" +
		"      config:\n" +
		"        subject_claim: \"sub\"\n" +
		"        localpart_template: \"{%% raw %%}{{ user.sub.split('@')[0] }}{%% endraw %%}\"\n" +
		"        display_name_template: \"{%% raw %%}{{ user.name|capitalize }}{%% endraw %%}\"\n" +
		"        email_template: \"{%% raw %%}{{ user.email }}{%% endraw %%}\"\n",
	"twitch": "" +
		"  - idp_id: %s\n" +
		"    idp_name: %q\n" +
		"    idp_brand: %q\n" +
		"    issuer: %q\n" +
		"    client_id: %q\n" +
		"    client_secret: %q\n" +
		"    client_auth_method: \"client_secret_post\"\n" +
		"    allow_existing_users: yes\n" +
		"    user_mapping_provider:\n" +
		"      config:\n" +
		"        localpart_template: \"{%% raw %%}{{ user.preferred_username }}{%% endraw %%}\"\n" +
		"        display_name_template: \"{%% raw %%}{{ user.name|capitalize }}{%% endraw %%}\"\n",
	"twitter": "" +
		"  - idp_id: %s\n" +
		"    idp_name: %q\n" +
		"    idp_brand: %q\n" +
		"    issuer: %q\n" +
		"    discover: no\n" +
		"    client_id: %q\n" +
		"    client_secret: %q\n" +
		"    pkce_method: \"always\"\n" +
		"    scopes: [\"offline.access\", \"tweet.read\", \"users.read\"]\n" +
		"    authorization_endpoint: \"https://twitter.com/i/oauth2/authorize\"\n" +
		"    token_endpoint: \"https://api.twitter.com/2/oauth2/token\"\n" +
		"    userinfo_endpoint: \"https://api.twitter.com/2/users/me?user.fields=profile_image_url\"\n" +
		"    allow_existing_users: yes\n" +
		"    user_mapping_provider:\n" +
		"      config:\n" +
		"        subject_template: \"{%% raw %%}{{ user.data.id }}{%% endraw %%}\"\n" +
		"        localpart_template: \"{%% raw %%}{{ user.data.username }}{%% endraw %%}\"\n" +
		"        display_name_template: \"{%% raw %%}{{ user.data.name }}{%% endraw %%}\"\n" +
		"        picture_template: \"{%% raw %%}{{ user.data.profile_image_url }}{%% endraw %%}\"\n",
	"xwiki": "" +
		"  - idp_id: %s\n" +
		"    idp_name: %q\n" +
		"    idp_brand: %q\n" +
		"    issuer: %q\n" +
		"    client_id: %q\n" +
		"    client_secret: %q\n" +
		"    client_auth_method: none\n" +
		"    user_profile_method: \"userinfo_endpoint\"\n" +
		"    scopes: [\"openid\", \"profile\"]\n" +
		"    allow_existing_users: yes\n" +
		"    user_mapping_provider:\n" +
		"      config:\n" +
		"        localpart_template: \"{%% raw %%}{{ user.preferred_username }}{%% endraw %%}\"\n" +
		"        display_name_template: \"{%% raw %%}{{ user.name }}{%% endraw %%}\"\n",
	"default": "" +
		"  - idp_id: %s\n" +
		"    idp_name: %q\n" +
		"    idp_brand: %q\n" +
		"    issuer: %q\n" +
		"    client_id: %q\n" +
		"    client_secret: %q\n" +
		"    scopes: [\"openid\", \"profile\"]\n" +
		"    allow_existing_users: yes\n" +
		"    user_mapping_provider:\n" +
		"      config:\n" +
		"        localpart_template: \"{%% raw %%}{{ user.given_name|lower }}{%% endraw %%}\"\n" +
		"        display_name_template: \"{%% raw %%}{{ user.name|capitalize }}{%% endraw %%}\"\n",
}

var domains = map[string]string{
	".etke.host":     "enTDpM8y67STAZcQMpmqr7",
	".kupo.email":    "X2dqkUMuGzqXenSx3Huz9T",
	".ma3x.chat":     "4Ys5JTRJ8Hyoip3UPDMVgj",
	".matrix.fan":    "x8RAPicWWX3kPurkL9ntVo",
	".matrix.town":   "Gn9RYjWvznHkLfzVdY6Gsa",
	".onmatrix.chat": "zVNMf3dur7oHP8dcGETZs",
}

var locations = map[string]string{
	"ashburn":     "ash",
	"falkenstein": "fsn1",
	"helsinki":    "hel1",
	"hillsboro":   "hil",
	"nuremberg":   "nbg1",
}

var defaultFirewall = map[string]int{
	"firewall": 124003,
}

var openFirewall = map[string]int{
	"firewall": 394512,
}
