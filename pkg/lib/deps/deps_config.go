package deps

import (
	"github.com/google/wire"

	"github.com/authgear/authgear-server/pkg/lib/config"
	"github.com/authgear/authgear-server/pkg/util/template"
)

var ConfigDeps = wire.NewSet(
	wire.FieldsOf(new(*config.Config), "AppConfig", "SecretConfig", "FeatureConfig"),
	wire.FieldsOf(new(*config.AppConfig),
		"ID",
		"HTTP",
		"Hook",
		"UI",
		"Localization",
		"Messaging",
		"Authentication",
		"Session",
		"OAuth",
		"Identity",
		"Authenticator",
		"UserProfile",
		"ForgotPassword",
		"WelcomeMessage",
		"Verification",
		"OTP",
		"AccountDeletion",
		"AccountAnonymization",
		"Web3",
		"GoogleTagManager",
		"AccountMigration",
		"Captcha",
	),
	wire.FieldsOf(new(*config.AuthenticationConfig),
		"Lockout",
	),
	wire.FieldsOf(new(*config.IdentityConfig),
		"LoginID",
		"OAuth",
		"Biometric",
		"OnConflict",
	),
	wire.FieldsOf(new(*config.MessagingConfig),
		"SMS",
		"Email",
		"Whatsapp",
		"RateLimits",
	),
	wire.FieldsOf(new(*config.AuthenticatorConfig),
		"Password",
		"TOTP",
		"OOB",
	),
	wire.FieldsOf(new(*config.AccountMigrationConfig),
		"Hook",
	),
	wire.FieldsOf(new(*config.FeatureConfig),
		"Identity",
		"Authenticator",
		"UI",
		"OAuth",
		"AuditLog",
		"Collaborator",
		"RateLimit",
		"RateLimits",
		"Messaging",
		"AdminAPI",
		"TestMode",
	),
	ProvideDefaultLanguageTag,
	ProvideSupportedLanguageTags,
	ProvideTestModeEmailSuppressed,
	ProvideTestModeSMSSuppressed,
	ProvideTestModeWhatsappSuppressed,
	secretDeps,
)

func ProvideDefaultLanguageTag(c *config.Config) template.DefaultLanguageTag {
	return template.DefaultLanguageTag(*c.AppConfig.Localization.FallbackLanguage)
}

func ProvideSupportedLanguageTags(c *config.Config) template.SupportedLanguageTags {
	return template.SupportedLanguageTags(c.AppConfig.Localization.SupportedLanguages)
}

func ProvideTestModeEmailSuppressed(c *config.TestModeFeatureConfig) config.TestModeEmailSuppressed {
	return config.TestModeEmailSuppressed(c.Email.Suppressed)
}

func ProvideTestModeSMSSuppressed(c *config.TestModeFeatureConfig) config.TestModeSMSSuppressed {
	return config.TestModeSMSSuppressed(c.SMS.Suppressed)
}

func ProvideTestModeWhatsappSuppressed(c *config.TestModeFeatureConfig) config.TestModeWhatsappSuppressed {
	return config.TestModeWhatsappSuppressed(c.Whatsapp.Suppressed)
}

var secretDeps = wire.NewSet(
	ProvideDatabaseCredentials,
	ProvideAuditDatabaseCredentials,
	ProvideElasticsearchCredentials,
	ProvideRedisCredentials,
	ProvideAnalyticRedisCredentials,
	ProvideAdminAPIAuthKeyMaterials,
	ProvideOAuthSSOProviderCredentials,
	ProvideSMTPServerCredentials,
	ProvideTwilioCredentials,
	ProvideNexmoCredentials,
	ProvideCustomSMSProviderConfig,
	ProvideOAuthKeyMaterials,
	ProvideCSRFKeyMaterials,
	ProvideWebhookKeyMaterials,
	ProvideImagesKeyMaterials,
	ProvideWATICredentials,
	ProvideOAuthClientCredentials,
	ProvideCaptchaCloudflareCredentials,
	ProvideWhatsappOnPremisesCredentials,
)

func ProvideDatabaseCredentials(c *config.SecretConfig) *config.DatabaseCredentials {
	s, _ := c.LookupData(config.DatabaseCredentialsKey).(*config.DatabaseCredentials)
	return s
}

func ProvideAuditDatabaseCredentials(c *config.SecretConfig) *config.AuditDatabaseCredentials {
	s, _ := c.LookupData(config.AuditDatabaseCredentialsKey).(*config.AuditDatabaseCredentials)
	return s
}

func ProvideElasticsearchCredentials(c *config.SecretConfig) *config.ElasticsearchCredentials {
	s, _ := c.LookupData(config.ElasticsearchCredentialsKey).(*config.ElasticsearchCredentials)
	return s
}

func ProvideRedisCredentials(c *config.SecretConfig) *config.RedisCredentials {
	s, _ := c.LookupData(config.RedisCredentialsKey).(*config.RedisCredentials)
	return s
}

func ProvideAnalyticRedisCredentials(c *config.SecretConfig) *config.AnalyticRedisCredentials {
	s, _ := c.LookupData(config.AnalyticRedisCredentialsKey).(*config.AnalyticRedisCredentials)
	return s
}

func ProvideAdminAPIAuthKeyMaterials(c *config.SecretConfig) *config.AdminAPIAuthKey {
	s, _ := c.LookupData(config.AdminAPIAuthKeyKey).(*config.AdminAPIAuthKey)
	return s
}

func ProvideOAuthSSOProviderCredentials(c *config.SecretConfig) *config.OAuthSSOProviderCredentials {
	s, _ := c.LookupData(config.OAuthSSOProviderCredentialsKey).(*config.OAuthSSOProviderCredentials)
	return s
}

func ProvideSMTPServerCredentials(c *config.SecretConfig) *config.SMTPServerCredentials {
	s, _ := c.LookupData(config.SMTPServerCredentialsKey).(*config.SMTPServerCredentials)
	return s
}

func ProvideTwilioCredentials(c *config.SecretConfig) *config.TwilioCredentials {
	s, _ := c.LookupData(config.TwilioCredentialsKey).(*config.TwilioCredentials)
	return s
}

func ProvideNexmoCredentials(c *config.SecretConfig) *config.NexmoCredentials {
	s, _ := c.LookupData(config.NexmoCredentialsKey).(*config.NexmoCredentials)
	return s
}

func ProvideCustomSMSProviderConfig(c *config.SecretConfig) *config.CustomSMSProviderConfig {
	return c.GetCustomSMSProviderConfig()
}

func ProvideOAuthKeyMaterials(c *config.SecretConfig) *config.OAuthKeyMaterials {
	s, _ := c.LookupData(config.OAuthKeyMaterialsKey).(*config.OAuthKeyMaterials)
	return s
}

func ProvideCSRFKeyMaterials(c *config.SecretConfig) *config.CSRFKeyMaterials {
	s, _ := c.LookupData(config.CSRFKeyMaterialsKey).(*config.CSRFKeyMaterials)
	return s
}

func ProvideWebhookKeyMaterials(c *config.SecretConfig) *config.WebhookKeyMaterials {
	s, _ := c.LookupData(config.WebhookKeyMaterialsKey).(*config.WebhookKeyMaterials)
	return s
}

func ProvideImagesKeyMaterials(c *config.SecretConfig) *config.ImagesKeyMaterials {
	s, _ := c.LookupData(config.ImagesKeyMaterialsKey).(*config.ImagesKeyMaterials)
	return s
}

func ProvideWATICredentials(c *config.SecretConfig) *config.WATICredentials {
	s, _ := c.LookupData(config.WATICredentialsKey).(*config.WATICredentials)
	return s
}

func ProvideOAuthClientCredentials(c *config.SecretConfig) *config.OAuthClientCredentials {
	s, _ := c.LookupData(config.OAuthClientCredentialsKey).(*config.OAuthClientCredentials)
	return s
}

func ProvideCaptchaCloudflareCredentials(c *config.SecretConfig) *config.CaptchaCloudflareCredentials {
	s, _ := c.LookupData(config.CaptchaCloudflareCredentialsKey).(*config.CaptchaCloudflareCredentials)
	return s
}

func ProvideWhatsappOnPremisesCredentials(c *config.SecretConfig) *config.WhatsappOnPremisesCredentials {
	s, _ := c.LookupData(config.WhatsappOnPremisesCredentialsKey).(*config.WhatsappOnPremisesCredentials)
	return s
}
