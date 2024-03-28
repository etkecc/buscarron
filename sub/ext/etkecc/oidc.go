package etkecc

import (
	"fmt"
	"net/url"
	"strings"
)

// excluded by purpose: apple, dex, keycloak
var oidcConfigs = map[string]string{
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

func (o *order) getOIDCConfig() string {
	id := strings.ToLower(o.get("sso-idp-id"))
	name := o.get("sso-idp-name")
	brand := strings.ToLower(o.get("sso-idp-brand"))
	issuer := o.get("sso-issuer")
	clientID := o.get("sso-client-id")
	clientSecret := o.get("sso-client-secret")

	provider := "default"
	keys := []string{id, name, brand}
	for _, key := range keys {
		if _, ok := oidcConfigs[key]; ok {
			provider = key
			break
		}
	}
	config := fmt.Sprintf(oidcConfigs[provider], id, name, brand, issuer, clientID, clientSecret)

	// special case: OIDC providers that require to use specific endpoints instead of auto-discovery
	issuerHost := "example.com"
	issuerURL, err := url.Parse(issuer)
	if err == nil {
		issuerHost = issuerURL.Host
		if provider == "microsoft" { // tenant id from issuer url: https://login.microsoftonline.com/<tenant id>/v2.0
			issuerHost = strings.Split(strings.TrimPrefix(issuerURL.Path, "/"), "/")[0]
		}
	}

	return strings.ReplaceAll(config, "example.com", issuerHost)
}
