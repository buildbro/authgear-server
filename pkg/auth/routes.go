package auth

import (
	apihandler "github.com/authgear/authgear-server/pkg/auth/handler/api"
	oauthhandler "github.com/authgear/authgear-server/pkg/auth/handler/oauth"
	siwehandler "github.com/authgear/authgear-server/pkg/auth/handler/siwe"
	webapphandler "github.com/authgear/authgear-server/pkg/auth/handler/webapp"
	"github.com/authgear/authgear-server/pkg/auth/webapp"
	"github.com/authgear/authgear-server/pkg/lib/config/configsource"
	"github.com/authgear/authgear-server/pkg/lib/deps"
	"github.com/authgear/authgear-server/pkg/lib/infra/middleware"
	"github.com/authgear/authgear-server/pkg/lib/oauth"
	"github.com/authgear/authgear-server/pkg/util/httproute"
	"github.com/authgear/authgear-server/pkg/util/httputil"
)

func newUnsafeDynamicCSPMiddleware(deps *deps.RequestProvider) httproute.Middleware {
	return newDynamicCSPMiddleware(deps, webapp.AllowInlineScript(true), webapp.AllowFrameAncestors(false))
}

func newSafeDynamicCSPMiddleware(deps *deps.RequestProvider) httproute.Middleware {
	return newDynamicCSPMiddleware(deps, webapp.AllowInlineScript(false), webapp.AllowFrameAncestors(false))
}

func newConsentPageDynamicCSPMiddleware(deps *deps.RequestProvider) httproute.Middleware {
	return newDynamicCSPMiddleware(deps, webapp.AllowInlineScript(false), webapp.AllowFrameAncestors(true))
}

func newAllSessionMiddleware(deps *deps.RequestProvider) httproute.Middleware {
	return newSessionMiddleware(deps, false)
}

func newIDPSessionOnlySessionMiddleware(deps *deps.RequestProvider) httproute.Middleware {
	return newSessionMiddleware(deps, true)
}

func NewRouter(p *deps.RootProvider, configSource *configsource.ConfigSource) *httproute.Router {

	newSessionMiddleware := func(idpSessionOnly bool) httproute.Middleware {
		if idpSessionOnly {
			return p.Middleware(newIDPSessionOnlySessionMiddleware)
		}
		return p.Middleware(newAllSessionMiddleware)
	}

	router := httproute.NewRouter()

	router.Add(httproute.Route{
		Methods:     []string{"GET"},
		PathPattern: "/healthz",
	}, p.RootHandler(newHealthzHandler))

	rootChain := httproute.Chain(
		p.RootMiddleware(newPanicMiddleware),
		p.RootMiddleware(newBodyLimitMiddleware),
		p.RootMiddleware(newSentryMiddleware),
		httproute.MiddlewareFunc(httputil.StaticSecurityHeaders),
		RequestMiddleware(p, configSource, newRequestMiddleware),
	)

	// This route is intentionally simple.
	// This does not check Host and allow any origin.
	generatedStaticChain := httproute.Chain(
		httproute.MiddlewareFunc(httputil.StaticSecurityHeaders),
		httproute.MiddlewareFunc(middleware.CORSStar),
	)

	appStaticChain := httproute.Chain(
		rootChain,
		p.Middleware(newCORSMiddleware),
		p.Middleware(newPublicOriginMiddleware),
	)

	oauthStaticChain := httproute.Chain(
		rootChain,
		p.Middleware(newCORSMiddleware),
		p.Middleware(newPublicOriginMiddleware),
	)

	newOAuthAPIChain := func(idpSessionOnly bool) httproute.Middleware {
		return httproute.Chain(
			rootChain,
			p.Middleware(newCORSMiddleware),
			p.Middleware(newPublicOriginMiddleware),
			newSessionMiddleware(idpSessionOnly),
			httproute.MiddlewareFunc(httputil.NoStore),
			p.Middleware(newWebAppWeChatRedirectURIMiddleware),
		)
	}

	oauthAPIChain := newOAuthAPIChain(false)
	// authz endpoint only accepts idp session
	oauthAuthzAPIChain := newOAuthAPIChain(true)
	siweAPIChain := httproute.Chain(
		rootChain,
		p.Middleware(newCORSMiddleware),
		p.Middleware(newPublicOriginMiddleware),
		httproute.MiddlewareFunc(httputil.NoCache),
	)

	apiChain := httproute.Chain(
		rootChain,
		p.Middleware(newCORSMiddleware),
		p.Middleware(newPublicOriginMiddleware),
		p.Middleware(newAllSessionMiddleware),
		httproute.MiddlewareFunc(httputil.NoStore),
	)
	workflowChain := httproute.Chain(
		apiChain,
		p.Middleware(newWorkflowUIParamMiddleware),
		p.Middleware(newWorkflowIntlMiddleware),
	)

	apiAuthenticatedChain := httproute.Chain(
		apiChain,
		p.Middleware(newAPIRRequireAuthenticatedMiddlewareMiddleware),
	)

	oauthAPIScopedChain := httproute.Chain(
		rootChain,
		p.Middleware(newCORSMiddleware),
		p.Middleware(newPublicOriginMiddleware),
		p.Middleware(newAllSessionMiddleware),
		httproute.MiddlewareFunc(httputil.NoStore),
		// Current we only require valid session and do not require any scope.
		httproute.MiddlewareFunc(oauth.RequireScope()),
	)

	newWebappChain := func(idpSessionOnly bool) httproute.Middleware {
		return httproute.Chain(
			rootChain,
			p.Middleware(newPublicOriginMiddleware),
			p.Middleware(newPanicWebAppMiddleware),
			newSessionMiddleware(idpSessionOnly),
			httproute.MiddlewareFunc(httputil.NoStore),
			httproute.MiddlewareFunc(webapp.IntlMiddleware),
			p.Middleware(newWebAppSessionMiddleware),
			p.Middleware(newWebAppUIParamMiddleware),
			p.Middleware(newWebAppColorSchemeMiddleware),
			p.Middleware(newWebAppWeChatRedirectURIMiddleware),
			p.Middleware(newTutorialMiddleware),
		)
	}
	webappChain := newWebappChain(false)
	webappSSOCallbackChain := httproute.Chain(
		webappChain,
	)
	webappWebsocketChain := httproute.Chain(
		webappChain,
	)
	webappAPIChain := httproute.Chain(
		webappChain,
	)

	newWebappPageChain := func(idpSessionOnly bool) httproute.Middleware {
		return httproute.Chain(
			newWebappChain(idpSessionOnly),
			p.Middleware(newCSRFMiddleware),
			// Turbo no longer requires us to tell the redirected location.
			// It can now determine redirection from the response.
			// https://github.com/hotwired/turbo/blob/daabebb0575fffbae1b2582dc458967cd638e899/src/core/drive/visit.ts#L316
			p.Middleware(newSafeDynamicCSPMiddleware),
		)
	}
	webappPageChain := newWebappPageChain(false)
	webappSIWEChain := httproute.Chain(
		webappChain,
		p.Middleware(newCSRFMiddleware),
		p.Middleware(newUnsafeDynamicCSPMiddleware),
	)
	webappAuthEntrypointChain := httproute.Chain(
		webappPageChain,
		p.Middleware(newAuthEntryPointMiddleware),
		// A unique visit is started when the user visit auth entry point
		p.Middleware(newWebAppVisitorIDMiddleware),
	)
	webappRequireAuthEnabledAuthEntrypointChain := httproute.Chain(
		webappPageChain,
		p.Middleware(newImplementationSwitcherMiddleware),
		p.Middleware(newRequireAuthenticationEnabledMiddleware),
		p.Middleware(newAuthEntryPointMiddleware),
		// A unique visit is started when the user visit auth entry point
		p.Middleware(newWebAppVisitorIDMiddleware),
	)
	webappPromoteChain := httproute.Chain(
		webappPageChain,
		p.Middleware(newImplementationSwitcherMiddleware),
		p.Middleware(newRequireAuthenticationEnabledMiddleware),
		p.Middleware(newAuthEntryPointMiddleware),
	)
	// select account page only accepts idp session
	webappSelectAccountChain := httproute.Chain(
		newWebappPageChain(true),
		p.Middleware(newImplementationSwitcherMiddleware),
		p.Middleware(newAuthEntryPointMiddleware),
	)
	// consent page only accepts idp session
	webappConsentPageChain := httproute.Chain(
		newWebappChain(true),
		p.Middleware(newCSRFMiddleware),
		p.Middleware(newConsentPageDynamicCSPMiddleware),
	)
	webappAuthenticatedChain := httproute.Chain(
		webappPageChain,
		webapp.RequireAuthenticatedMiddleware{},
	)
	webappSuccessPageChain := httproute.Chain(
		webappPageChain,
		// SuccessPageMiddleware check the cookie and see if it is valid to
		// visit the success page
		p.Middleware(newSuccessPageMiddleware),
	)
	webappSettingsChain := httproute.Chain(
		webappAuthenticatedChain,
		p.Middleware(newRequireSettingsEnabledMiddleware),
	)
	webappSettingsSubRoutesChain := httproute.Chain(
		webappAuthenticatedChain,
		p.Middleware(newRequireSettingsEnabledMiddleware),
		// SettingsSubRoutesMiddleware should be added to all the settings sub routes only
		// but no /settings itself
		// it redirects all sub routes to /settings if the current user is
		// anonymous user
		p.Middleware(newSettingsSubRoutesMiddleware),
	)

	appStaticRoute := httproute.Route{Middleware: appStaticChain}
	generatedStaticRoute := httproute.Route{Middleware: generatedStaticChain}
	oauthStaticRoute := httproute.Route{Middleware: oauthStaticChain}
	oauthAPIRoute := httproute.Route{Middleware: oauthAPIChain}
	oauthAuthzAPIRoute := httproute.Route{Middleware: oauthAuthzAPIChain}
	siweAPIRoute := httproute.Route{Middleware: siweAPIChain}
	apiRoute := httproute.Route{Middleware: apiChain}
	workflowRoute := httproute.Route{Middleware: workflowChain}
	apiAuthenticatedRoute := httproute.Route{Middleware: apiAuthenticatedChain}
	oauthAPIScopedRoute := httproute.Route{Middleware: oauthAPIScopedChain}
	webappPageRoute := httproute.Route{Middleware: webappPageChain}
	webappPromoteRoute := httproute.Route{Middleware: webappPromoteChain}
	webappSIWERoute := httproute.Route{Middleware: webappSIWEChain}
	webappAuthEntrypointRoute := httproute.Route{Middleware: webappAuthEntrypointChain}
	webappRequireAuthEnabledAuthEntrypointRoute := httproute.Route{Middleware: webappRequireAuthEnabledAuthEntrypointChain}
	webappSelectAccountRoute := httproute.Route{Middleware: webappSelectAccountChain}
	webappConsentPageRoute := httproute.Route{Middleware: webappConsentPageChain}
	webappAuthenticatedRoute := httproute.Route{Middleware: webappAuthenticatedChain}
	webappSuccessPageRoute := httproute.Route{Middleware: webappSuccessPageChain}
	webappSettingsRoute := httproute.Route{Middleware: webappSettingsChain}
	webappSettingsSubRoutesRoute := httproute.Route{Middleware: webappSettingsSubRoutesChain}
	webappTesterRouter := httproute.Route{Middleware: webappChain}
	webappSSOCallbackRoute := httproute.Route{Middleware: webappSSOCallbackChain}
	webappWebsocketRoute := httproute.Route{Middleware: webappWebsocketChain}
	webappAPIRoute := httproute.Route{Middleware: webappAPIChain}

	router.Add(webapphandler.ConfigureRootRoute(webappAuthEntrypointRoute), p.Handler(newWebAppRootHandler))
	router.Add(webapphandler.ConfigureOAuthEntrypointRoute(webappAuthEntrypointRoute), p.Handler(newWebAppOAuthEntrypointHandler))
	router.Add(webapphandler.ConfigureAuthflowLoginRoute(webappRequireAuthEnabledAuthEntrypointRoute), &webapphandler.ImplementationSwitcherHandler{
		Interaction: p.Handler(newWebAppLoginHandler),
		Authflow:    p.Handler(newWebAppAuthflowLoginHandler),
	})
	router.Add(webapphandler.ConfigureAuthflowSignupRoute(webappRequireAuthEnabledAuthEntrypointRoute), &webapphandler.ImplementationSwitcherHandler{
		Interaction: p.Handler(newWebAppSignupHandler),
		Authflow:    p.Handler(newWebAppAuthflowSignupHandler),
	})
	router.Add(webapphandler.ConfigureAuthflowPromoteRoute(webappPromoteRoute), &webapphandler.ImplementationSwitcherHandler{
		Interaction: p.Handler(newWebAppPromoteHandler),
		Authflow:    p.Handler(newWebAppAuthflowPromoteHandler),
	})
	router.Add(webapphandler.ConfigureAuthflowReauthRoute(webappSelectAccountRoute), &webapphandler.ImplementationSwitcherHandler{
		Interaction: p.Handler(newWebAppReauthHandler),
		Authflow:    p.Handler(newWebAppAuthflowReauthHandler),
	})

	router.Add(webapphandler.ConfigureSelectAccountRoute(webappSelectAccountRoute), p.Handler(newWebAppSelectAccountHandler))

	router.Add(webapphandler.ConfigureAuthflowEnterPasswordRoute(webappPageRoute), p.Handler(newWebAppAuthflowEnterPasswordHandler))
	router.Add(webapphandler.ConfigureAuthflowEnterOOBOTPRoute(webappPageRoute), p.Handler(newWebAppAuthflowEnterOOBOTPHandler))
	router.Add(webapphandler.ConfigureAuthflowCreatePasswordRoute(webappPageRoute), p.Handler(newWebAppAuthflowCreatePasswordHandler))
	router.Add(webapphandler.ConfigureAuthflowEnterTOTPRoute(webappPageRoute), p.Handler(newWebAppAuthflowEnterTOTPHandler))
	router.Add(webapphandler.ConfigureAuthflowSetupTOTPRoute(webappPageRoute), p.Handler(newWebAppAuthflowSetupTOTPHandler))
	router.Add(webapphandler.ConfigureAuthflowViewRecoveryCodeRoute(webappPageRoute), p.Handler(newWebAppAuthflowViewRecoveryCodeHandler))
	router.Add(webapphandler.ConfigureAuthflowWhatsappOTPRoute(webappPageRoute), p.Handler(newWebAppAuthflowWhatsappOTPHandler))
	router.Add(webapphandler.ConfigureAuthflowOOBOTPLinkRoute(webappPageRoute), p.Handler(newWebAppAuthflowOOBOTPLinkHandler))
	router.Add(webapphandler.ConfigureAuthflowChangePasswordRoute(webappPageRoute), p.Handler(newWebAppAuthflowChangePasswordHandler))
	router.Add(webapphandler.ConfigureAuthflowUsePasskeyRoute(webappPageRoute), p.Handler(newWebAppAuthflowUsePasskeyHandler))
	router.Add(webapphandler.ConfigureAuthflowPromptCreatePasskeyRoute(webappPageRoute), p.Handler(newWebAppAuthflowPromptCreatePasskeyHandler))
	router.Add(webapphandler.ConfigureAuthflowEnterRecoveryCodeRoute(webappPageRoute), p.Handler(newWebAppAuthflowEnterRecoveryCodeHandler))
	router.Add(webapphandler.ConfigureAuthflowSetupOOBOTPRoute(webappPageRoute), p.Handler(newWebAppAuthflowSetupOOBOTPHandler))
	router.Add(webapphandler.ConfigureAuthflowTerminateOtherSessionsRoute(webappPageRoute), p.Handler(newWebAppAuthflowTerminateOtherSessionsHandler))
	router.Add(webapphandler.ConfigureAuthflowAccountStatusRoute(webappPageRoute), p.Handler(newWebAppAuthflowAccountStatusHandler))
	router.Add(webapphandler.ConfigureAuthflowWechatRoute(webappPageRoute), p.Handler(newWebAppAuthflowWechatHandler))
	router.Add(webapphandler.ConfigureAuthflowForgotPasswordRoute(webappPageRoute), p.Handler(newWebAppAuthflowForgotPasswordHandler))
	router.Add(webapphandler.ConfigureAuthflowForgotPasswordSuccessRoute(webappPageRoute), p.Handler(newWebAppAuthflowForgotPasswordSuccessHandler))
	router.Add(webapphandler.ConfigureAuthflowResetPasswordRoute(webappPageRoute), p.Handler(newWebAppAuthflowResetPasswordHandler))
	router.Add(webapphandler.ConfigureAuthflowResetPasswordSuccessRoute(webappPageRoute), p.Handler(newWebAppAuthflowResetPasswordSuccessHandler))

	router.Add(webapphandler.ConfigureEnterPasswordRoute(webappPageRoute), p.Handler(newWebAppEnterPasswordHandler))
	router.Add(webapphandler.ConfigureConfirmTerminateOtherSessionsRoute(webappPageRoute), p.Handler(newWebConfirmTerminateOtherSessionsHandler))
	router.Add(webapphandler.ConfigureUsePasskeyRoute(webappPageRoute), p.Handler(newWebAppUsePasskeyHandler))
	router.Add(webapphandler.ConfigureSetupTOTPRoute(webappPageRoute), p.Handler(newWebAppSetupTOTPHandler))
	router.Add(webapphandler.ConfigureEnterTOTPRoute(webappPageRoute), p.Handler(newWebAppEnterTOTPHandler))
	router.Add(webapphandler.ConfigureSetupOOBOTPRoute(webappPageRoute), p.Handler(newWebAppSetupOOBOTPHandler))
	router.Add(webapphandler.ConfigureEnterOOBOTPRoute(webappPageRoute), p.Handler(newWebAppEnterOOBOTPHandler))
	router.Add(webapphandler.ConfigureSetupWhatsappOTPRoute(webappPageRoute), p.Handler(newWebAppSetupWhatsappOTPHandler))
	router.Add(webapphandler.ConfigureWhatsappOTPRoute(webappPageRoute), p.Handler(newWebAppWhatsappOTPHandler))
	router.Add(webapphandler.ConfigureSetupLoginLinkOTPRoute(webappPageRoute), p.Handler(newWebAppSetupLoginLinkOTPHandler))
	router.Add(webapphandler.ConfigureLoginLinkOTPRoute(webappPageRoute), p.Handler(newWebAppLoginLinkOTPHandler))
	router.Add(webapphandler.ConfigureVerifyLoginLinkOTPRoute(webappPageRoute), p.Handler(newWebAppVerifyLoginLinkOTPHandler))
	router.Add(webapphandler.ConfigureEnterRecoveryCodeRoute(webappPageRoute), p.Handler(newWebAppEnterRecoveryCodeHandler))
	router.Add(webapphandler.ConfigureSetupRecoveryCodeRoute(webappPageRoute), p.Handler(newWebAppSetupRecoveryCodeHandler))
	router.Add(webapphandler.ConfigureVerifyIdentityRoute(webappPageRoute), p.Handler(newWebAppVerifyIdentityHandler))
	router.Add(webapphandler.ConfigureVerifyIdentitySuccessRoute(webappPageRoute), p.Handler(newWebAppVerifyIdentitySuccessHandler))
	router.Add(webapphandler.ConfigureCreatePasswordRoute(webappPageRoute), p.Handler(newWebAppCreatePasswordHandler))
	router.Add(webapphandler.ConfigureCreatePasskeyRoute(webappPageRoute), p.Handler(newWebAppCreatePasskeyHandler))
	router.Add(webapphandler.ConfigurePromptCreatePasskeyRoute(webappPageRoute), p.Handler(newWebAppPromptCreatePasskeyHandler))
	router.Add(webapphandler.ConfigureForgotPasswordRoute(webappPageRoute), p.Handler(newWebAppForgotPasswordHandler))
	router.Add(webapphandler.ConfigureResetPasswordRoute(webappPageRoute), p.Handler(newWebAppResetPasswordHandler))
	router.Add(webapphandler.ConfigureAccountStatusRoute(webappPageRoute), p.Handler(newWebAppAccountStatusHandler))
	router.Add(webapphandler.ConfigureReturnRoute(webappPageRoute), p.Handler(newWebAppReturnHandler))
	router.Add(webapphandler.ConfigureErrorRoute(webappPageRoute), p.Handler(newWebAppErrorHandler))
	router.Add(webapphandler.ConfigureForceChangePasswordRoute(webappPageRoute), p.Handler(newWebAppForceChangePasswordHandler))
	router.Add(webapphandler.ConfigureForceChangeSecondaryPasswordRoute(webappPageRoute), p.Handler(newWebAppForceChangeSecondaryPasswordHandler))
	router.Add(webapphandler.ConfigureConnectWeb3AccountRoute(webappSIWERoute), p.Handler(newWebAppConnectWeb3AccountHandler))
	router.Add(webapphandler.ConfigureMissingWeb3WalletRoute(webappPageRoute), p.Handler(newWebAppMissingWeb3WalletHandler))
	router.Add(webapphandler.ConfigureFeatureDisabledRoute(webappPageRoute), p.Handler(newWebAppFeatureDisabledHandler))

	router.Add(webapphandler.ConfigureForgotPasswordSuccessRoute(webappSuccessPageRoute), p.Handler(newWebAppForgotPasswordSuccessHandler))
	router.Add(webapphandler.ConfigureResetPasswordSuccessRoute(webappSuccessPageRoute), p.Handler(newWebAppResetPasswordSuccessHandler))
	router.Add(webapphandler.ConfigureSettingsDeleteAccountSuccessRoute(webappSuccessPageRoute), p.Handler(newWebAppSettingsDeleteAccountSuccessHandler))

	router.Add(webapphandler.ConfigureLogoutRoute(webappAuthenticatedRoute), p.Handler(newWebAppLogoutHandler))
	router.Add(webapphandler.ConfigureEnterLoginIDRoute(webappAuthenticatedRoute), p.Handler(newWebAppEnterLoginIDHandler))
	router.Add(webapphandler.ConfigureSettingsRoute(webappSettingsRoute), p.Handler(newWebAppSettingsHandler))

	router.Add(webapphandler.ConfigureSettingsProfileRoute(webappSettingsSubRoutesRoute), p.Handler(newWebAppSettingsProfileHandler))
	router.Add(webapphandler.ConfigureSettingsProfileEditRoute(webappSettingsSubRoutesRoute), p.Handler(newWebAppSettingsProfileEditHandler))
	router.Add(webapphandler.ConfigureSettingsIdentityRoute(webappSettingsSubRoutesRoute), p.Handler(newWebAppSettingsIdentityHandler))
	router.Add(webapphandler.ConfigureSettingsBiometricRoute(webappSettingsSubRoutesRoute), p.Handler(newWebAppSettingsBiometricHandler))
	router.Add(webapphandler.ConfigureSettingsMFARoute(webappSettingsSubRoutesRoute), p.Handler(newWebAppSettingsMFAHandler))
	router.Add(webapphandler.ConfigureSettingsTOTPRoute(webappSettingsSubRoutesRoute), p.Handler(newWebAppSettingsTOTPHandler))
	router.Add(webapphandler.ConfigureSettingsPasskeyRoute(webappSettingsSubRoutesRoute), p.Handler(newWebAppSettingsPasskeyHandler))
	router.Add(webapphandler.ConfigureSettingsOOBOTPRoute(webappSettingsSubRoutesRoute), p.Handler(newWebAppSettingsOOBOTPHandler))
	router.Add(webapphandler.ConfigureSettingsRecoveryCodeRoute(webappSettingsSubRoutesRoute), p.Handler(newWebAppSettingsRecoveryCodeHandler))
	router.Add(webapphandler.ConfigureSettingsSessionsRoute(webappSettingsSubRoutesRoute), p.Handler(newWebAppSettingsSessionsHandler))
	router.Add(webapphandler.ConfigureSettingsChangePasswordRoute(webappSettingsSubRoutesRoute), p.Handler(newWebAppSettingsChangePasswordHandler))
	router.Add(webapphandler.ConfigureSettingsChangeSecondaryPasswordRoute(webappSettingsSubRoutesRoute), p.Handler(newWebAppSettingsChangeSecondaryPasswordHandler))
	router.Add(webapphandler.ConfigureSettingsDeleteAccountRoute(webappSettingsSubRoutesRoute), p.Handler(newWebAppSettingsDeleteAccountHandler))

	router.Add(webapphandler.ConfigureTesterRoute(webappTesterRouter), p.Handler(newWebAppTesterHandler))

	router.Add(webapphandler.ConfigureSSOCallbackRoute(webappSSOCallbackRoute), p.Handler(newWebAppSSOCallbackHandler))

	router.Add(webapphandler.ConfigureWechatAuthRoute(webappPageRoute), p.Handler(newWechatAuthHandler))
	router.Add(webapphandler.ConfigureWechatCallbackRoute(webappSSOCallbackRoute), p.Handler(newWechatCallbackHandler))

	router.Add(webapphandler.ConfigurePasskeyCreationOptionsRoute(webappAPIRoute), p.Handler(newWebAppPasskeyCreationOptionsHandler))
	router.Add(webapphandler.ConfigurePasskeyRequestOptionsRoute(webappAPIRoute), p.Handler(newWebAppPasskeyRequestOptionsHandler))

	router.Add(oauthhandler.ConfigureOIDCMetadataRoute(oauthStaticRoute), p.Handler(newOAuthMetadataHandler))
	router.Add(oauthhandler.ConfigureOAuthMetadataRoute(oauthStaticRoute), p.Handler(newOAuthMetadataHandler))
	router.Add(oauthhandler.ConfigureJWKSRoute(oauthStaticRoute), p.Handler(newOAuthJWKSHandler))

	router.Add(oauthhandler.ConfigureAuthorizeRoute(oauthAuthzAPIRoute), p.Handler(newOAuthAuthorizeHandler))
	router.Add(oauthhandler.ConfigureTokenRoute(oauthAPIRoute), p.Handler(newOAuthTokenHandler))
	router.Add(oauthhandler.ConfigureRevokeRoute(oauthAPIRoute), p.Handler(newOAuthRevokeHandler))
	router.Add(oauthhandler.ConfigureEndSessionRoute(oauthAPIRoute), p.Handler(newOAuthEndSessionHandler))

	router.Add(oauthhandler.ConfigureChallengeRoute(apiRoute), p.Handler(newOAuthChallengeHandler))
	router.Add(oauthhandler.ConfigureAppSessionTokenRoute(apiRoute), p.Handler(newOAuthAppSessionTokenHandler))
	router.Add(oauthhandler.ConfigureProxyRedirectRoute(apiRoute), p.Handler(newOAuthProxyRedirectHandler))

	router.Add(oauthhandler.ConfigureUserInfoRoute(oauthAPIScopedRoute), p.Handler(newOAuthUserInfoHandler))

	router.Add(oauthhandler.ConfigureConsentRoute(webappConsentPageRoute), p.Handler(newOAuthConsentHandler))

	router.Add(siwehandler.ConfigureNonceRoute(siweAPIRoute), p.Handler(newSIWENonceHandler))

	router.Add(apihandler.ConfigureAnonymousUserSignupRoute(apiRoute), p.Handler(newAPIAnonymousUserSignupHandler))
	router.Add(apihandler.ConfigureAnonymousUserPromotionCodeRoute(apiRoute), p.Handler(newAPIAnonymousUserPromotionCodeHandler))
	router.Add(apihandler.ConfigurePresignImagesUploadRoute(apiAuthenticatedRoute), p.Handler(newAPIPresignImagesUploadHandler))

	router.Add(webapphandler.ConfigureWebsocketRoute(webappWebsocketRoute), p.Handler(newWebAppWebsocketHandler))

	router.Add(webapphandler.ConfigureAppStaticAssetsRoute(appStaticRoute), p.Handler(newWebAppAppStaticAssetsHandler))

	router.Add(webapphandler.ConfigureGeneratedStaticAssetsRoute(generatedStaticRoute), p.RootHandler(newWebAppGeneratedStaticAssetsHandler))

	router.Add(apihandler.ConfigureWorkflowNewRoute(workflowRoute), p.Handler(newAPIWorkflowNewHandler))
	router.Add(apihandler.ConfigureWorkflowGetRoute(workflowRoute), p.Handler(newAPIWorkflowGetHandler))
	router.Add(apihandler.ConfigureWorkflowInputRoute(workflowRoute), p.Handler(newAPIWorkflowInputHandler))
	router.Add(apihandler.ConfigureWorkflowWebsocketRoute(workflowRoute), p.Handler(newAPIWorkflowWebsocketHandler))
	router.Add(apihandler.ConfigureWorkflowV2Route(workflowRoute), p.Handler(newAPIWorkflowV2Handler))

	router.Add(apihandler.ConfigureAuthenticationFlowV1CreateRoute(workflowRoute), p.Handler(newAPIAuthenticationFlowV1CreateHandler))
	router.Add(apihandler.ConfigureAuthenticationFlowV1InputRoute(workflowRoute), p.Handler(newAPIAuthenticationFlowV1InputHandler))
	router.Add(apihandler.ConfigureAuthenticationFlowV1GetRoute(workflowRoute), p.Handler(newAPIAuthenticationFlowV1GetHandler))
	router.Add(apihandler.ConfigureAuthenticationFlowV1WebsocketRoute(workflowRoute), p.Handler(newAPIAuthenticationFlowV1WebsocketHandler))

	router.NotFound(webappPageRoute, p.Handler(newWebAppNotFoundHandler))

	return router
}
