// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package admin

import (
	"github.com/authgear/authgear-server/pkg/admin/graphql"
	"github.com/authgear/authgear-server/pkg/admin/loader"
	service3 "github.com/authgear/authgear-server/pkg/admin/service"
	"github.com/authgear/authgear-server/pkg/admin/transport"
	"github.com/authgear/authgear-server/pkg/lib/admin/authz"
	"github.com/authgear/authgear-server/pkg/lib/authn/authenticator/oob"
	"github.com/authgear/authgear-server/pkg/lib/authn/authenticator/password"
	service2 "github.com/authgear/authgear-server/pkg/lib/authn/authenticator/service"
	"github.com/authgear/authgear-server/pkg/lib/authn/authenticator/totp"
	"github.com/authgear/authgear-server/pkg/lib/authn/challenge"
	"github.com/authgear/authgear-server/pkg/lib/authn/identity/anonymous"
	"github.com/authgear/authgear-server/pkg/lib/authn/identity/loginid"
	"github.com/authgear/authgear-server/pkg/lib/authn/identity/oauth"
	"github.com/authgear/authgear-server/pkg/lib/authn/identity/service"
	"github.com/authgear/authgear-server/pkg/lib/authn/mfa"
	"github.com/authgear/authgear-server/pkg/lib/authn/otp"
	"github.com/authgear/authgear-server/pkg/lib/authn/sso"
	"github.com/authgear/authgear-server/pkg/lib/authn/user"
	"github.com/authgear/authgear-server/pkg/lib/config"
	"github.com/authgear/authgear-server/pkg/lib/deps"
	"github.com/authgear/authgear-server/pkg/lib/facade"
	"github.com/authgear/authgear-server/pkg/lib/feature/forgotpassword"
	"github.com/authgear/authgear-server/pkg/lib/feature/verification"
	"github.com/authgear/authgear-server/pkg/lib/feature/welcomemessage"
	"github.com/authgear/authgear-server/pkg/lib/hook"
	"github.com/authgear/authgear-server/pkg/lib/infra/db"
	"github.com/authgear/authgear-server/pkg/lib/infra/middleware"
	"github.com/authgear/authgear-server/pkg/lib/interaction"
	"github.com/authgear/authgear-server/pkg/lib/session/access"
	"github.com/authgear/authgear-server/pkg/lib/session/idpsession"
	"github.com/authgear/authgear-server/pkg/lib/translation"
	"github.com/authgear/authgear-server/pkg/util/clock"
	"github.com/authgear/authgear-server/pkg/util/httproute"
	"github.com/authgear/authgear-server/pkg/util/rand"
	"net/http"
)

// Injectors from wire.go:

func newSentryMiddleware(p *deps.RootProvider) httproute.Middleware {
	hub := p.SentryHub
	environmentConfig := p.EnvironmentConfig
	trustProxy := environmentConfig.TrustProxy
	sentryMiddleware := &middleware.SentryMiddleware{
		SentryHub:  hub,
		TrustProxy: trustProxy,
	}
	return sentryMiddleware
}

func newPanicEndMiddleware(p *deps.RootProvider) httproute.Middleware {
	panicEndMiddleware := &middleware.PanicEndMiddleware{}
	return panicEndMiddleware
}

func newPanicWriteEmptyResponseMiddleware(p *deps.RootProvider) httproute.Middleware {
	panicWriteEmptyResponseMiddleware := &middleware.PanicWriteEmptyResponseMiddleware{}
	return panicWriteEmptyResponseMiddleware
}

func newBodyLimitMiddleware(p *deps.RootProvider) httproute.Middleware {
	bodyLimitMiddleware := &middleware.BodyLimitMiddleware{}
	return bodyLimitMiddleware
}

func newPanicLogMiddleware(p *deps.RequestProvider) httproute.Middleware {
	appProvider := p.AppProvider
	factory := appProvider.LoggerFactory
	panicLogMiddlewareLogger := middleware.NewPanicLogMiddlewareLogger(factory)
	panicLogMiddleware := &middleware.PanicLogMiddleware{
		Logger: panicLogMiddlewareLogger,
	}
	return panicLogMiddleware
}

func newAuthorizationMiddleware(p *deps.RequestProvider, auth config.AdminAPIAuth) httproute.Middleware {
	appProvider := p.AppProvider
	factory := appProvider.LoggerFactory
	logger := authz.NewLogger(factory)
	configConfig := appProvider.Config
	appConfig := configConfig.AppConfig
	appID := appConfig.ID
	secretConfig := configConfig.SecretConfig
	adminAPIAuthKey := deps.ProvideAdminAPIAuthKeyMaterials(secretConfig)
	clock := _wireSystemClockValue
	authzMiddleware := &authz.Middleware{
		Logger:  logger,
		Auth:    auth,
		AppID:   appID,
		AuthKey: adminAPIAuthKey,
		Clock:   clock,
	}
	return authzMiddleware
}

var (
	_wireSystemClockValue = clock.NewSystemClock()
)

func newGraphQLHandler(p *deps.RequestProvider) http.Handler {
	appProvider := p.AppProvider
	factory := appProvider.LoggerFactory
	logger := graphql.NewLogger(factory)
	configConfig := appProvider.Config
	secretConfig := configConfig.SecretConfig
	databaseCredentials := deps.ProvideDatabaseCredentials(secretConfig)
	appConfig := configConfig.AppConfig
	appID := appConfig.ID
	sqlBuilder := db.ProvideSQLBuilder(databaseCredentials, appID)
	request := p.Request
	context := deps.ProvideRequestContext(request)
	handle := appProvider.Database
	sqlExecutor := db.SQLExecutor{
		Context:  context,
		Database: handle,
	}
	store := &user.Store{
		SQLBuilder:  sqlBuilder,
		SQLExecutor: sqlExecutor,
	}
	authenticationConfig := appConfig.Authentication
	identityConfig := appConfig.Identity
	serviceStore := &service.Store{
		SQLBuilder:  sqlBuilder,
		SQLExecutor: sqlExecutor,
	}
	loginidStore := &loginid.Store{
		SQLBuilder:  sqlBuilder,
		SQLExecutor: sqlExecutor,
	}
	loginIDConfig := identityConfig.LoginID
	rootProvider := appProvider.RootProvider
	reservedNameChecker := rootProvider.ReservedNameChecker
	typeCheckerFactory := &loginid.TypeCheckerFactory{
		Config:              loginIDConfig,
		ReservedNameChecker: reservedNameChecker,
	}
	checker := &loginid.Checker{
		Config:             loginIDConfig,
		TypeCheckerFactory: typeCheckerFactory,
	}
	normalizerFactory := &loginid.NormalizerFactory{
		Config: loginIDConfig,
	}
	clockClock := _wireSystemClockValue
	provider := &loginid.Provider{
		Store:             loginidStore,
		Config:            loginIDConfig,
		Checker:           checker,
		NormalizerFactory: normalizerFactory,
		Clock:             clockClock,
	}
	oauthStore := &oauth.Store{
		SQLBuilder:  sqlBuilder,
		SQLExecutor: sqlExecutor,
	}
	oauthProvider := &oauth.Provider{
		Store: oauthStore,
		Clock: clockClock,
	}
	anonymousStore := &anonymous.Store{
		SQLBuilder:  sqlBuilder,
		SQLExecutor: sqlExecutor,
	}
	anonymousProvider := &anonymous.Provider{
		Store: anonymousStore,
		Clock: clockClock,
	}
	serviceService := &service.Service{
		Authentication: authenticationConfig,
		Identity:       identityConfig,
		Store:          serviceStore,
		LoginID:        provider,
		OAuth:          oauthProvider,
		Anonymous:      anonymousProvider,
	}
	store2 := &service2.Store{
		SQLBuilder:  sqlBuilder,
		SQLExecutor: sqlExecutor,
	}
	passwordStore := &password.Store{
		SQLBuilder:  sqlBuilder,
		SQLExecutor: sqlExecutor,
	}
	authenticatorConfig := appConfig.Authenticator
	authenticatorPasswordConfig := authenticatorConfig.Password
	passwordLogger := password.NewLogger(factory)
	historyStore := &password.HistoryStore{
		Clock:       clockClock,
		SQLBuilder:  sqlBuilder,
		SQLExecutor: sqlExecutor,
	}
	passwordChecker := password.ProvideChecker(authenticatorPasswordConfig, historyStore)
	queue := appProvider.TaskQueue
	passwordProvider := &password.Provider{
		Store:           passwordStore,
		Config:          authenticatorPasswordConfig,
		Clock:           clockClock,
		Logger:          passwordLogger,
		PasswordHistory: historyStore,
		PasswordChecker: passwordChecker,
		TaskQueue:       queue,
	}
	totpStore := &totp.Store{
		SQLBuilder:  sqlBuilder,
		SQLExecutor: sqlExecutor,
	}
	authenticatorTOTPConfig := authenticatorConfig.TOTP
	totpProvider := &totp.Provider{
		Store:  totpStore,
		Config: authenticatorTOTPConfig,
		Clock:  clockClock,
	}
	authenticatorOOBConfig := authenticatorConfig.OOB
	oobStore := &oob.Store{
		SQLBuilder:  sqlBuilder,
		SQLExecutor: sqlExecutor,
	}
	oobProvider := &oob.Provider{
		Config: authenticatorOOBConfig,
		Store:  oobStore,
		Clock:  clockClock,
	}
	service4 := &service2.Service{
		Store:    store2,
		Password: passwordProvider,
		TOTP:     totpProvider,
		OOBOTP:   oobProvider,
	}
	verificationLogger := verification.NewLogger(factory)
	verificationConfig := appConfig.Verification
	redisHandle := appProvider.Redis
	storeRedis := &verification.StoreRedis{
		Redis: redisHandle,
		AppID: appID,
		Clock: clockClock,
	}
	storePQ := &verification.StorePQ{
		SQLBuilder:  sqlBuilder,
		SQLExecutor: sqlExecutor,
	}
	verificationService := &verification.Service{
		Logger:     verificationLogger,
		Config:     verificationConfig,
		Clock:      clockClock,
		CodeStore:  storeRedis,
		ClaimStore: storePQ,
	}
	coordinator := &facade.Coordinator{
		Identities:     serviceService,
		Authenticators: service4,
		Verification:   verificationService,
		IdentityConfig: identityConfig,
	}
	identityFacade := facade.IdentityFacade{
		Coordinator: coordinator,
	}
	queries := &user.Queries{
		Store:        store,
		Identities:   identityFacade,
		Verification: verificationService,
	}
	userLoader := &loader.UserLoader{
		Users: queries,
	}
	interactionLogger := interaction.NewLogger(factory)
	authenticatorFacade := facade.AuthenticatorFacade{
		Coordinator: coordinator,
	}
	environmentConfig := rootProvider.EnvironmentConfig
	staticAssetURLPrefix := environmentConfig.StaticAssetURLPrefix
	engine := appProvider.TemplateEngine
	translationService := &translation.Service{
		Context:           context,
		EnvironmentConfig: environmentConfig,
		TemplateEngine:    engine,
	}
	webEndpoints := &WebEndpoints{}
	messageSender := &otp.MessageSender{
		StaticAssetURLPrefix: staticAssetURLPrefix,
		Translation:          translationService,
		Endpoints:            webEndpoints,
		TaskQueue:            queue,
	}
	codeSender := &oob.CodeSender{
		OTPMessageSender: messageSender,
	}
	oAuthClientCredentials := deps.ProvideOAuthClientCredentials(secretConfig)
	userInfoDecoder := sso.UserInfoDecoder{
		LoginIDNormalizerFactory: normalizerFactory,
	}
	oAuthProviderFactory := &sso.OAuthProviderFactory{
		Endpoints:                webEndpoints,
		IdentityConfig:           identityConfig,
		Credentials:              oAuthClientCredentials,
		RedirectURL:              webEndpoints,
		Clock:                    clockClock,
		UserInfoDecoder:          userInfoDecoder,
		LoginIDNormalizerFactory: normalizerFactory,
	}
	storeDeviceTokenRedis := &mfa.StoreDeviceTokenRedis{
		Redis: redisHandle,
		AppID: appID,
		Clock: clockClock,
	}
	storeRecoveryCodePQ := &mfa.StoreRecoveryCodePQ{
		SQLBuilder:  sqlBuilder,
		SQLExecutor: sqlExecutor,
	}
	mfaService := &mfa.Service{
		DeviceTokens:  storeDeviceTokenRedis,
		RecoveryCodes: storeRecoveryCodePQ,
		Clock:         clockClock,
		Config:        authenticationConfig,
	}
	forgotPasswordConfig := appConfig.ForgotPassword
	forgotpasswordStore := &forgotpassword.Store{
		Redis: redisHandle,
	}
	providerLogger := forgotpassword.NewProviderLogger(factory)
	forgotpasswordProvider := &forgotpassword.Provider{
		StaticAssetURLPrefix: staticAssetURLPrefix,
		Translation:          translationService,
		Config:               forgotPasswordConfig,
		Store:                forgotpasswordStore,
		Clock:                clockClock,
		URLs:                 webEndpoints,
		TaskQueue:            queue,
		Logger:               providerLogger,
		Identities:           identityFacade,
		Authenticators:       authenticatorFacade,
	}
	verificationCodeSender := &verification.CodeSender{
		OTPMessageSender: messageSender,
		WebAppURLs:       webEndpoints,
	}
	challengeProvider := &challenge.Provider{
		Redis: redisHandle,
		AppID: appID,
		Clock: clockClock,
	}
	welcomeMessageConfig := appConfig.WelcomeMessage
	welcomemessageProvider := &welcomemessage.Provider{
		Translation:          translationService,
		WelcomeMessageConfig: welcomeMessageConfig,
		TaskQueue:            queue,
	}
	rawCommands := &user.RawCommands{
		Store:                  store,
		Clock:                  clockClock,
		WelcomeMessageProvider: welcomemessageProvider,
		Queries:                queries,
	}
	hookLogger := hook.NewLogger(factory)
	rawProvider := &user.RawProvider{
		RawCommands: rawCommands,
		Queries:     queries,
	}
	hookStore := &hook.Store{
		SQLBuilder:  sqlBuilder,
		SQLExecutor: sqlExecutor,
	}
	hookConfig := appConfig.Hook
	webhookKeyMaterials := deps.ProvideWebhookKeyMaterials(secretConfig)
	syncHTTPClient := hook.NewSyncHTTPClient(hookConfig)
	asyncHTTPClient := hook.NewAsyncHTTPClient()
	deliverer := &hook.Deliverer{
		Config:    hookConfig,
		Secret:    webhookKeyMaterials,
		Clock:     clockClock,
		SyncHTTP:  syncHTTPClient,
		AsyncHTTP: asyncHTTPClient,
	}
	hookProvider := &hook.Provider{
		Context:   context,
		Logger:    hookLogger,
		Database:  handle,
		Clock:     clockClock,
		Users:     rawProvider,
		Store:     hookStore,
		Deliverer: deliverer,
	}
	commands := &user.Commands{
		Raw:          rawCommands,
		Hooks:        hookProvider,
		Verification: verificationService,
	}
	userProvider := &user.Provider{
		Commands: commands,
		Queries:  queries,
	}
	trustProxy := environmentConfig.TrustProxy
	cookieFactory := deps.NewCookieFactory(request, trustProxy)
	storeRedisLogger := idpsession.NewStoreRedisLogger(factory)
	idpsessionStoreRedis := &idpsession.StoreRedis{
		Redis:  redisHandle,
		AppID:  appID,
		Clock:  clockClock,
		Logger: storeRedisLogger,
	}
	eventStoreRedis := &access.EventStoreRedis{
		Redis: redisHandle,
		AppID: appID,
	}
	eventProvider := &access.EventProvider{
		Store: eventStoreRedis,
	}
	sessionConfig := appConfig.Session
	rand := _wireRandValue
	idpsessionProvider := &idpsession.Provider{
		Request:      request,
		Store:        idpsessionStoreRedis,
		AccessEvents: eventProvider,
		TrustProxy:   trustProxy,
		Config:       sessionConfig,
		Clock:        clockClock,
		Random:       rand,
	}
	httpConfig := appConfig.HTTP
	cookieDef := idpsession.NewSessionCookieDef(httpConfig, sessionConfig)
	mfaCookieDef := mfa.NewDeviceTokenCookieDef(httpConfig, authenticationConfig)
	interactionContext := &interaction.Context{
		Database:                 sqlExecutor,
		Clock:                    clockClock,
		Config:                   appConfig,
		Identities:               identityFacade,
		Authenticators:           authenticatorFacade,
		AnonymousIdentities:      anonymousProvider,
		OOBAuthenticators:        oobProvider,
		OOBCodeSender:            codeSender,
		OAuthProviderFactory:     oAuthProviderFactory,
		MFA:                      mfaService,
		ForgotPassword:           forgotpasswordProvider,
		ResetPassword:            forgotpasswordProvider,
		LoginIDNormalizerFactory: normalizerFactory,
		Verification:             verificationService,
		VerificationCodeSender:   verificationCodeSender,
		Challenges:               challengeProvider,
		Users:                    userProvider,
		Hooks:                    hookProvider,
		CookieFactory:            cookieFactory,
		Sessions:                 idpsessionProvider,
		SessionCookie:            cookieDef,
		MFADeviceTokenCookie:     mfaCookieDef,
	}
	interactionStoreRedis := &interaction.StoreRedis{
		Redis: redisHandle,
		AppID: appID,
	}
	interactionService := &interaction.Service{
		Logger:  interactionLogger,
		Context: interactionContext,
		Store:   interactionStoreRedis,
	}
	serviceInteractionService := &service3.InteractionService{
		Graph: interactionService,
	}
	identityLoader := &loader.IdentityLoader{
		Identities:  serviceService,
		Interaction: serviceInteractionService,
	}
	authenticatorLoader := &loader.AuthenticatorLoader{
		Authenticators: service4,
	}
	graphqlContext := &graphql.Context{
		GQLLogger:      logger,
		Users:          userLoader,
		Identities:     identityLoader,
		Authenticators: authenticatorLoader,
	}
	devMode := environmentConfig.DevMode
	graphQLHandler := &transport.GraphQLHandler{
		GraphQLContext: graphqlContext,
		DevMode:        devMode,
		Database:       handle,
	}
	return graphQLHandler
}

var (
	_wireRandValue = idpsession.Rand(rand.SecureRand)
)
