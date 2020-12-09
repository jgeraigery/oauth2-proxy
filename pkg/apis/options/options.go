package options

import (
	"crypto"
	"net/url"
	"regexp"

	oidc "github.com/coreos/go-oidc"
	ipapi "github.com/oauth2-proxy/oauth2-proxy/pkg/apis/ip"
	"github.com/oauth2-proxy/oauth2-proxy/providers"
	"github.com/spf13/pflag"
)

// SignatureData holds hmacauth signature hash and key
type SignatureData struct {
	Hash crypto.Hash
	Key  string
}

// Options holds Configuration Options that can be set by Command Line Flag,
// or Config File
type Options struct {
	ProxyPrefix        string   `flag:"proxy-prefix" cfg:"proxy_prefix"`
	PingPath           string   `flag:"ping-path" cfg:"ping_path"`
	PingUserAgent      string   `flag:"ping-user-agent" cfg:"ping_user_agent"`
	HTTPAddress        string   `flag:"http-address" cfg:"http_address"`
	HTTPSAddress       string   `flag:"https-address" cfg:"https_address"`
	ReverseProxy       bool     `flag:"reverse-proxy" cfg:"reverse_proxy"`
	RealClientIPHeader string   `flag:"real-client-ip-header" cfg:"real_client_ip_header"`
	TrustedIPs         []string `flag:"trusted-ip" cfg:"trusted_ips"`
	ForceHTTPS         bool     `flag:"force-https" cfg:"force_https"`
	RawRedirectURL     string   `flag:"redirect-url" cfg:"redirect_url"`
	TLSCertFile        string   `flag:"tls-cert-file" cfg:"tls_cert_file"`
	TLSKeyFile         string   `flag:"tls-key-file" cfg:"tls_key_file"`

	AuthenticatedEmailsFile string   `flag:"authenticated-emails-file" cfg:"authenticated_emails_file"`
	EmailDomains            []string `flag:"email-domain" cfg:"email_domains"`
	WhitelistDomains        []string `flag:"whitelist-domain" cfg:"whitelist_domains"`

	HtpasswdFile        string `flag:"htpasswd-file" cfg:"htpasswd_file"`
	DisplayHtpasswdForm bool   `flag:"display-htpasswd-form" cfg:"display_htpasswd_form"`
	CustomTemplatesDir  string `flag:"custom-templates-dir" cfg:"custom_templates_dir"`
	Banner              string `flag:"banner" cfg:"banner"`
	Footer              string `flag:"footer" cfg:"footer"`

	Cookie  Cookie         `cfg:",squash"`
	Session SessionOptions `cfg:",squash"`
	Logging Logging        `cfg:",squash"`

	// Not used in the legacy config, name not allowed to match an external key (upstreams)
	// TODO(JoelSpeed): Rename when legacy config is removed
	UpstreamServers Upstreams `cfg:",internal"`

	Providers Providers `cfg:",internal"`

	SkipAuthRegex         []string `flag:"skip-auth-regex" cfg:"skip_auth_regex"`
	SkipAuthStripHeaders  bool     `flag:"skip-auth-strip-headers" cfg:"skip_auth_strip_headers"`
	SkipJwtBearerTokens   bool     `flag:"skip-jwt-bearer-tokens" cfg:"skip_jwt_bearer_tokens"`
	ExtraJwtIssuers       []string `flag:"extra-jwt-issuers" cfg:"extra_jwt_issuers"`
	PassBasicAuth         bool     `flag:"pass-basic-auth" cfg:"pass_basic_auth"`
	SetBasicAuth          bool     `flag:"set-basic-auth" cfg:"set_basic_auth"`
	PreferEmailToUser     bool     `flag:"prefer-email-to-user" cfg:"prefer_email_to_user"`
	BasicAuthPassword     string   `flag:"basic-auth-password" cfg:"basic_auth_password"`
	PassAccessToken       bool     `flag:"pass-access-token" cfg:"pass_access_token"`
	SkipProviderButton    bool     `flag:"skip-provider-button" cfg:"skip_provider_button"`
	PassUserHeaders       bool     `flag:"pass-user-headers" cfg:"pass_user_headers"`
	SSLInsecureSkipVerify bool     `flag:"ssl-insecure-skip-verify" cfg:"ssl_insecure_skip_verify"`
	SetXAuthRequest       bool     `flag:"set-xauthrequest" cfg:"set_xauthrequest"`
	SetAuthorization      bool     `flag:"set-authorization-header" cfg:"set_authorization_header"`
	PassAuthorization     bool     `flag:"pass-authorization-header" cfg:"pass_authorization_header"`
	SkipAuthPreflight     bool     `flag:"skip-auth-preflight" cfg:"skip_auth_preflight"`

	SignatureKey    string `flag:"signature-key" cfg:"signature_key"`
	GCPHealthChecks bool   `flag:"gcp-healthchecks" cfg:"gcp_healthchecks"`

	// internal values that are set after config validation
	redirectURL        *url.URL
	compiledRegex      []*regexp.Regexp
	providers          map[string]providers.Provider
	signatureData      *SignatureData
	oidcVerifiers      map[string]*oidc.IDTokenVerifier
	jwtBearerVerifiers []*oidc.IDTokenVerifier
	realClientIPParser ipapi.RealClientIPParser
}

// Options for Getting internal values
func (o *Options) GetRedirectURL() *url.URL                           { return o.redirectURL }
func (o *Options) GetCompiledRegex() []*regexp.Regexp                 { return o.compiledRegex }
func (o *Options) GetProviders() map[string]providers.Provider        { return o.providers }
func (o *Options) GetSignatureData() *SignatureData                   { return o.signatureData }
func (o *Options) GetOIDCVerifiers() map[string]*oidc.IDTokenVerifier { return o.oidcVerifiers }
func (o *Options) GetJWTBearerVerifiers() []*oidc.IDTokenVerifier     { return o.jwtBearerVerifiers }
func (o *Options) GetRealClientIPParser() ipapi.RealClientIPParser    { return o.realClientIPParser }

// Options for Setting internal values
func (o *Options) SetRedirectURL(s *url.URL)                         { o.redirectURL = s }
func (o *Options) SetCompiledRegex(s []*regexp.Regexp)               { o.compiledRegex = s }
func (o *Options) SetProvider(p string, s providers.Provider)        { o.providers[p] = s }
func (o *Options) SetSignatureData(s *SignatureData)                 { o.signatureData = s }
func (o *Options) SetOIDCVerifier(p string, s *oidc.IDTokenVerifier) { o.oidcVerifiers[p] = s }
func (o *Options) SetJWTBearerVerifiers(s []*oidc.IDTokenVerifier)   { o.jwtBearerVerifiers = s }
func (o *Options) SetRealClientIPParser(s ipapi.RealClientIPParser)  { o.realClientIPParser = s }

// NewOptions constructs a new Options with defaulted values
func NewOptions() *Options {
	return &Options{
		ProxyPrefix:         "/oauth2",
		Providers:           providerDefaults(),
		PingPath:            "/ping",
		HTTPAddress:         "127.0.0.1:4180",
		HTTPSAddress:        ":443",
		RealClientIPHeader:  "X-Real-IP",
		ForceHTTPS:          false,
		DisplayHtpasswdForm: true,
		Cookie:              cookieDefaults(),
		Session:             sessionOptionsDefaults(),
		UpstreamServers:     Upstreams{},
		// AzureTenant:                      "common",
		SetXAuthRequest:   false,
		SkipAuthPreflight: false,
		PassBasicAuth:     true,
		SetBasicAuth:      false,
		PassUserHeaders:   true,
		PassAccessToken:   false,
		SetAuthorization:  false,
		PassAuthorization: false,
		PreferEmailToUser: false,

		// Prompt:                           "", // Change to "login" when ApprovalPrompt officially deprecated
		// ApprovalPrompt:                   "force",
		// UserIDClaim:                      "email",
		// InsecureOIDCAllowUnverifiedEmail: false,
		// SkipOIDCDiscovery:                false,
		Logging:       loggingDefaults(),
		providers:     make(map[string]providers.Provider),
		oidcVerifiers: make(map[string]*oidc.IDTokenVerifier),
	}
}

// NewFlagSet creates a new FlagSet with all of the flags required by Options
func NewFlagSet() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("oauth2-proxy", pflag.ExitOnError)

	flagSet.String("http-address", "127.0.0.1:4180", "[http://]<addr>:<port> or unix://<path> to listen on for HTTP clients")
	flagSet.String("https-address", ":443", "<addr>:<port> to listen on for HTTPS clients")
	flagSet.Bool("reverse-proxy", false, "are we running behind a reverse proxy, controls whether headers like X-Real-Ip are accepted")
	flagSet.String("real-client-ip-header", "X-Real-IP", "Header used to determine the real IP of the client (one of: X-Forwarded-For, X-Real-IP, or X-ProxyUser-IP)")
	flagSet.StringSlice("trusted-ip", []string{}, "list of IPs or CIDR ranges to allow to bypass authentication. WARNING: trusting by IP has inherent security flaws, read the configuration documentation for more information.")
	flagSet.Bool("force-https", false, "force HTTPS redirect for HTTP requests")
	flagSet.String("tls-cert-file", "", "path to certificate file")
	flagSet.String("tls-key-file", "", "path to private key file")
	flagSet.String("redirect-url", "", "the OAuth Redirect URL. ie: \"https://internalapp.yourcompany.com/oauth2/callback\"")
	flagSet.Bool("set-xauthrequest", false, "set X-Auth-Request-User and X-Auth-Request-Email response headers (useful in Nginx auth_request mode)")
	flagSet.Bool("pass-basic-auth", true, "pass HTTP Basic Auth, X-Forwarded-User and X-Forwarded-Email information to upstream")
	flagSet.Bool("set-basic-auth", false, "set HTTP Basic Auth information in response (useful in Nginx auth_request mode)")
	flagSet.Bool("prefer-email-to-user", false, "Prefer to use the Email address as the Username when passing information to upstream. Will only use Username if Email is unavailable, eg. htaccess authentication. Used in conjunction with -pass-basic-auth and -pass-user-headers")
	flagSet.Bool("pass-user-headers", true, "pass X-Forwarded-User and X-Forwarded-Email information to upstream")
	flagSet.String("basic-auth-password", "", "the password to set when passing the HTTP Basic Auth header")
	flagSet.Bool("pass-access-token", false, "pass OAuth access_token to upstream via X-Forwarded-Access-Token header")
	flagSet.Bool("pass-authorization-header", false, "pass the Authorization Header to upstream")
	flagSet.Bool("set-authorization-header", false, "set Authorization response headers (useful in Nginx auth_request mode)")
	flagSet.StringSlice("skip-auth-regex", []string{}, "bypass authentication for requests path's that match (may be given multiple times)")
	flagSet.Bool("skip-auth-strip-headers", false, "strips X-Forwarded-* style authentication headers & Authorization header if they would be set by oauth2-proxy for request paths in --skip-auth-regex")
	flagSet.Bool("skip-provider-button", false, "will skip sign-in-page to directly reach the next step: oauth/start")
	flagSet.Bool("skip-auth-preflight", false, "will skip authentication for OPTIONS requests")
	flagSet.Bool("ssl-insecure-skip-verify", false, "skip validation of certificates presented when using HTTPS providers")
	flagSet.Bool("skip-jwt-bearer-tokens", false, "will skip requests that have verified JWT bearer tokens (default false)")
	flagSet.StringSlice("extra-jwt-issuers", []string{}, "if skip-jwt-bearer-tokens is set, a list of extra JWT issuer=audience pairs (where the issuer URL has a .well-known/openid-configuration or a .well-known/jwks.json)")

	flagSet.StringSlice("email-domain", []string{}, "authenticate emails with the specified domain (may be given multiple times). Use * to authenticate any email")
	flagSet.StringSlice("whitelist-domain", []string{}, "allowed domains for redirection after authentication. Prefix domain with a . to allow subdomains (eg .example.com)")
	flagSet.String("authenticated-emails-file", "", "authenticate against emails via file (one per line)")
	flagSet.String("htpasswd-file", "", "additionally authenticate against a htpasswd file. Entries must be created with \"htpasswd -s\" for SHA encryption or \"htpasswd -B\" for bcrypt encryption")
	flagSet.Bool("display-htpasswd-form", true, "display username / password login form if an htpasswd file is provided")
	flagSet.String("custom-templates-dir", "", "path to custom html templates")
	flagSet.String("banner", "", "custom banner string. Use \"-\" to disable default banner.")
	flagSet.String("footer", "", "custom footer string. Use \"-\" to disable default footer.")
	flagSet.String("proxy-prefix", "/oauth2", "the url root path that this proxy should be nested under (e.g. /<oauth2>/sign_in)")
	flagSet.String("ping-path", "/ping", "the ping endpoint that can be used for basic health checks")
	flagSet.String("ping-user-agent", "", "special User-Agent that will be used for basic health checks")
	flagSet.String("session-store-type", "cookie", "the session storage provider to use")
	flagSet.Bool("session-cookie-minimal", false, "strip OAuth tokens from cookie session stores if they aren't needed (cookie session store only)")
	flagSet.String("redis-connection-url", "", "URL of redis server for redis session storage (eg: redis://HOST[:PORT])")
	flagSet.String("redis-password", "", "Redis password. Applicable for all Redis configurations. Will override any password set in `--redis-connection-url`")
	flagSet.Bool("redis-use-sentinel", false, "Connect to redis via sentinels. Must set --redis-sentinel-master-name and --redis-sentinel-connection-urls to use this feature")
	flagSet.String("redis-sentinel-password", "", "Redis sentinel password. Used only for sentinel connection; any redis node passwords need to use `--redis-password`")
	flagSet.String("redis-sentinel-master-name", "", "Redis sentinel master name. Used in conjunction with --redis-use-sentinel")
	flagSet.String("redis-ca-path", "", "Redis custom CA path")
	flagSet.Bool("redis-insecure-skip-tls-verify", false, "Use insecure TLS connection to redis")
	flagSet.StringSlice("redis-sentinel-connection-urls", []string{}, "List of Redis sentinel connection URLs (eg redis://HOST[:PORT]). Used in conjunction with --redis-use-sentinel")
	flagSet.Bool("redis-use-cluster", false, "Connect to redis cluster. Must set --redis-cluster-connection-urls to use this feature")
	flagSet.StringSlice("redis-cluster-connection-urls", []string{}, "List of Redis cluster connection URLs (eg redis://HOST[:PORT]). Used in conjunction with --redis-use-cluster")

	flagSet.String("signature-key", "", "GAP-Signature request signature key (algorithm:secretkey)")
	flagSet.Bool("gcp-healthchecks", false, "Enable GCP/GKE healthcheck endpoints")

	flagSet.AddFlagSet(cookieFlagSet())
	flagSet.AddFlagSet(loggingFlagSet())
	flagSet.AddFlagSet(legacyUpstreamsFlagSet())

	return flagSet
}
