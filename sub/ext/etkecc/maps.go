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
