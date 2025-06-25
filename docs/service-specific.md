## Jellyfin

Note that while the config sets up the provider and other details for jellyfin, it does not automatically enable it.   The following plugin has to be installed and configured properly:

https://github.com/9p4/jellyfin-plugin-sso

See archived instructions here: https://web.archive.org/web/20250625053920/https://integrations.goauthentik.io/integrations/services/jellyfin/#oidc-configuration

Just another note, after setting up OIDC if there is still an error stating invalid redirect URI, ensure `Scheme Override` at the bottom of the plugin settings is set to `https`.