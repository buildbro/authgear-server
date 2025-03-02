package declarative

import (
	"context"
	"errors"
	"fmt"

	"github.com/authgear/authgear-server/pkg/api"
	"github.com/authgear/authgear-server/pkg/api/apierrors"
	"github.com/authgear/authgear-server/pkg/api/model"
	authflow "github.com/authgear/authgear-server/pkg/lib/authenticationflow"
	"github.com/authgear/authgear-server/pkg/lib/authn/authenticator"
	"github.com/authgear/authgear-server/pkg/lib/authn/identity"
	"github.com/authgear/authgear-server/pkg/lib/authn/otp"
	"github.com/authgear/authgear-server/pkg/lib/authn/sso"
	"github.com/authgear/authgear-server/pkg/lib/config"
	"github.com/authgear/authgear-server/pkg/lib/infra/mail"
	"github.com/authgear/authgear-server/pkg/lib/uiparam"
	"github.com/authgear/authgear-server/pkg/util/errorutil"
	"github.com/authgear/authgear-server/pkg/util/phone"
	"github.com/authgear/authgear-server/pkg/util/uuid"
	"github.com/iawaknahc/jsonschema/pkg/jsonpointer"
)

func authenticatorIsDefault(deps *authflow.Dependencies, userID string, authenticatorKind model.AuthenticatorKind) (isDefault bool, err error) {
	ais, err := deps.Authenticators.List(
		userID,
		authenticator.KeepKind(authenticatorKind),
		authenticator.KeepDefault,
	)
	if err != nil {
		return
	}

	isDefault = len(ais) == 0
	return
}

func flowRootObject(deps *authflow.Dependencies, flowReference authflow.FlowReference) (config.AuthenticationFlowObject, error) {
	switch flowReference.Type {
	case authflow.FlowTypeSignup:
		return flowRootObjectForSignupFlow(deps, flowReference)
	case authflow.FlowTypePromote:
		return flowRootObjectForPromoteFlow(deps, flowReference)
	case authflow.FlowTypeLogin:
		return flowRootObjectForLoginFlow(deps, flowReference)
	case authflow.FlowTypeSignupLogin:
		return flowRootObjectForSignupLoginFlow(deps, flowReference)
	case authflow.FlowTypeReauth:
		return flowRootObjectForReauthFlow(deps, flowReference)
	case authflow.FlowTypeAccountRecovery:
		return flowRootObjectForAccountRecoveryFlow(deps, flowReference)
	default:
		panic(fmt.Errorf("unexpected flow type: %v", flowReference.Type))
	}
}

func flowRootObjectForSignupFlow(deps *authflow.Dependencies, flowReference authflow.FlowReference) (config.AuthenticationFlowObject, error) {
	var root config.AuthenticationFlowObject

	for _, f := range deps.Config.AuthenticationFlow.SignupFlows {
		f := f
		if f.Name == flowReference.Name {
			root = f
			break
		}
	}
	if root == nil && flowReference.Name == nameGeneratedFlow {
		root = GenerateSignupFlowConfig(deps.Config)
	}

	if root == nil {
		return nil, ErrFlowNotFound
	}

	return root, nil
}

func flowRootObjectForPromoteFlow(deps *authflow.Dependencies, flowReference authflow.FlowReference) (config.AuthenticationFlowObject, error) {
	var root config.AuthenticationFlowObject

	for _, f := range deps.Config.AuthenticationFlow.PromoteFlows {
		f := f
		if f.Name == flowReference.Name {
			root = f
			break
		}
	}
	if root == nil && flowReference.Name == nameGeneratedFlow {
		root = GeneratePromoteFlowConfig(deps.Config)
	}

	if root == nil {
		return nil, ErrFlowNotFound
	}

	return root, nil
}

func flowRootObjectForLoginFlow(deps *authflow.Dependencies, flowReference authflow.FlowReference) (config.AuthenticationFlowObject, error) {
	var root config.AuthenticationFlowObject

	for _, f := range deps.Config.AuthenticationFlow.LoginFlows {
		f := f
		if f.Name == flowReference.Name {
			root = f
			break
		}
	}
	if root == nil && flowReference.Name == nameGeneratedFlow {
		root = GenerateLoginFlowConfig(deps.Config)
	}

	if root == nil {
		return nil, ErrFlowNotFound
	}

	return root, nil
}

func flowRootObjectForSignupLoginFlow(deps *authflow.Dependencies, flowReference authflow.FlowReference) (config.AuthenticationFlowObject, error) {
	var root config.AuthenticationFlowObject

	for _, f := range deps.Config.AuthenticationFlow.SignupLoginFlows {
		f := f
		if f.Name == flowReference.Name {
			root = f
			break
		}
	}
	if root == nil && flowReference.Name == nameGeneratedFlow {
		root = GenerateSignupLoginFlowConfig(deps.Config)
	}

	if root == nil {
		return nil, ErrFlowNotFound
	}

	return root, nil
}

func flowRootObjectForReauthFlow(deps *authflow.Dependencies, flowReference authflow.FlowReference) (config.AuthenticationFlowObject, error) {
	var root config.AuthenticationFlowObject

	for _, f := range deps.Config.AuthenticationFlow.ReauthFlows {
		f := f
		if f.Name == flowReference.Name {
			root = f
			break
		}
	}
	if root == nil && flowReference.Name == nameGeneratedFlow {
		root = GenerateReauthFlowConfig(deps.Config)
	}

	if root == nil {
		return nil, ErrFlowNotFound
	}

	return root, nil
}

func flowRootObjectForAccountRecoveryFlow(deps *authflow.Dependencies, flowReference authflow.FlowReference) (config.AuthenticationFlowObject, error) {
	var root config.AuthenticationFlowObject

	for _, f := range deps.Config.AuthenticationFlow.AccountRecoveryFlows {
		f := f
		if f.Name == flowReference.Name {
			root = f
			break
		}
	}
	if root == nil && flowReference.Name == nameGeneratedFlow {
		root = GenerateAccountRecoveryFlowConfig(deps.Config)
	}

	if root == nil {
		return nil, ErrFlowNotFound
	}

	return root, nil
}

// nolint: gocyclo
func getAuthenticationOptionsForLogin(ctx context.Context, deps *authflow.Dependencies, flows authflow.Flows, userID string, step *config.AuthenticationFlowLoginFlowStep) ([]AuthenticateOption, error) {
	options := []AuthenticateOption{}

	identities, err := deps.Identities.ListByUser(userID)
	if err != nil {
		return nil, err
	}

	authenticators, err := deps.Authenticators.List(userID)
	if err != nil {
		return nil, err
	}

	secondaryAuthenticators := authenticator.ApplyFilters(authenticators, authenticator.KeepKind(model.AuthenticatorKindSecondary))

	isOptional := step.IsOptional()
	userHasSomeSecondaryAuthenticators := len(secondaryAuthenticators) > 0

	findIdentity := func(targetStepName string) (*identity.Info, error) {
		// Find the target step from the root.
		targetStepFlow, err := authflow.FindTargetStep(flows.Root, targetStepName)
		if err != nil {
			return nil, err
		}

		target, ok := targetStepFlow.Intent.(IntentLoginFlowStepAuthenticateTarget)
		if !ok {
			return nil, InvalidTargetStep.NewWithInfo("invalid target_step", apierrors.Details{
				"target_step": targetStepName,
			})
		}

		info := target.GetIdentity(ctx, deps, flows.Replace(targetStepFlow))

		return info, nil
	}

	useAuthenticationOptionAddRecoveryCodes := func(options []AuthenticateOption, isOptional bool, userHasSomeSecondaryAuthenticators bool) []AuthenticateOption {
		shouldAdd := false
		switch {
		case !isOptional:
			// We always add recovery_code even though the end-user does not actually has any.
			// One case that this situation will happen:
			// 1. The project makes 2FA required, disable recovery code, and enable TOTP.
			// 2. Alice enrolls TOTP. She does not have recovery code.
			// 3. The project enables recovery code.
			//
			// Alice can still use her TOTP.
			shouldAdd = true
		case isOptional && userHasSomeSecondaryAuthenticators:
			shouldAdd = true
		}

		if shouldAdd {
			options = append(options, NewAuthenticateOptionRecoveryCode())
		}

		return options
	}

	useAuthenticationOptionAddPrimaryPassword := func(options []AuthenticateOption) []AuthenticateOption {
		// We always add primary_password even though the end-user does not actually has one.
		// Showing this branch is necessary to convince the frontend to show a primary password page, where
		// the end-user can trigger account recovery flow and create a new password.
		options = append(options, NewAuthenticateOptionPassword(
			config.AuthenticationFlowAuthenticationPrimaryPassword),
		)
		return options
	}

	useAuthenticationOptionAddSecondaryPassword := func(options []AuthenticateOption, isOptional bool, userHasSomeSecondaryAuthenticators bool) []AuthenticateOption {
		shouldAdd := false
		switch {
		case !isOptional:
			// We always add secondary_password even though the end-user does not actually has one.
			// One case that this situation will happen:
			// 1. The project makes 2FA required, and enable TOTP.
			// 2. Alice enrolls TOTP.
			// 3. The project enables secondary password, and disable TOTP.
			// 4. Alice does not have secondary password.
			//
			// If recovery code is also disabled, Alice is locked out.
			shouldAdd = true
		case isOptional && userHasSomeSecondaryAuthenticators:
			shouldAdd = true
		}

		if shouldAdd {
			options = append(options, NewAuthenticateOptionPassword(
				config.AuthenticationFlowAuthenticationSecondaryPassword,
			))
		}

		return options
	}

	useAuthenticationOptionAddTOTP := func(options []AuthenticateOption, isOptional bool, userHasSomeSecondaryAuthenticators bool) []AuthenticateOption {
		shouldAdd := false
		switch {
		case !isOptional:
			// We always add secondary_totp even though the end-user does not actually has one.
			// One case that this situation will happen:
			// 1. The project makes 2FA required, and enable OOBOTP.
			// 2. Alice enrolls OOBOTP with her phone number.
			// 3. The project enables TOTP.
			// 4. Alice does not have TOTP.
			//
			// Alice can still use her OOBOTP with phone number.
			shouldAdd = true
		case isOptional && userHasSomeSecondaryAuthenticators:
			shouldAdd = true
		}

		if shouldAdd {
			options = append(options, NewAuthenticateOptionTOTP())
		}

		return options
	}

	useAuthenticationOptionAddPasskey := func(options []AuthenticateOption, deps *authflow.Dependencies, userID string) ([]AuthenticateOption, error) {
		requestOptions, err := deps.PasskeyRequestOptionsService.MakeModalRequestOptionsWithUser(userID)
		if err != nil {
			return nil, err
		}

		options = append(options, NewAuthenticateOptionPasskey(requestOptions))
		return options, nil
	}

	useAuthenticationOptionAddPrimaryOOBOTPOfIdentity := func(options []AuthenticateOption, deps *authflow.Dependencies, authentication config.AuthenticationFlowAuthentication, info *identity.Info) []AuthenticateOption {
		option, ok := NewAuthenticateOptionOOBOTPFromIdentity(deps.Config.Authenticator.OOB, info)
		if !ok {
			return options
		}

		if option.Authentication != authentication {
			return options
		}

		options = append(options, *option)
		return options
	}

	useAuthenticationOptionAddPrimaryOOBOTPOfAllIdentities := func(options []AuthenticateOption, deps *authflow.Dependencies, authentication config.AuthenticationFlowAuthentication, infos []*identity.Info) []AuthenticateOption {
		for _, info := range infos {
			options = useAuthenticationOptionAddPrimaryOOBOTPOfIdentity(options, deps, authentication, info)
		}

		return options
	}

	useAuthenticationOptionAddSecondaryOOBOTP := func(options []AuthenticateOption, deps *authflow.Dependencies, authentication config.AuthenticationFlowAuthentication, infos []*authenticator.Info) []AuthenticateOption {
		for _, info := range infos {
			if option, ok := NewAuthenticateOptionOOBOTPFromAuthenticator(deps.Config.Authenticator.OOB, info); ok {
				if option.Authentication == authentication {
					options = append(options, *option)
				}
			}
		}
		return options
	}

	for _, branch := range step.OneOf {
		switch branch.Authentication {
		case config.AuthenticationFlowAuthenticationDeviceToken:
			// Device token is handled transparently.
			break
		case config.AuthenticationFlowAuthenticationRecoveryCode:
			options = useAuthenticationOptionAddRecoveryCodes(options, isOptional, userHasSomeSecondaryAuthenticators)
		case config.AuthenticationFlowAuthenticationPrimaryPassword:
			options = useAuthenticationOptionAddPrimaryPassword(options)
		case config.AuthenticationFlowAuthenticationPrimaryPasskey:
			options, err = useAuthenticationOptionAddPasskey(options, deps, userID)
			if err != nil {
				return nil, err
			}
		case config.AuthenticationFlowAuthenticationSecondaryPassword:
			options = useAuthenticationOptionAddSecondaryPassword(options, isOptional, userHasSomeSecondaryAuthenticators)
		case config.AuthenticationFlowAuthenticationSecondaryTOTP:
			options = useAuthenticationOptionAddTOTP(options, isOptional, userHasSomeSecondaryAuthenticators)
		case config.AuthenticationFlowAuthenticationPrimaryOOBOTPEmail:
			fallthrough
		case config.AuthenticationFlowAuthenticationPrimaryOOBOTPSMS:
			if targetStepName := branch.TargetStep; targetStepName != "" {
				info, err := findIdentity(targetStepName)
				if err != nil {
					return nil, err
				}

				options = useAuthenticationOptionAddPrimaryOOBOTPOfIdentity(options, deps, branch.Authentication, info)
			} else {
				options = useAuthenticationOptionAddPrimaryOOBOTPOfAllIdentities(options, deps, branch.Authentication, identities)
			}
		case config.AuthenticationFlowAuthenticationSecondaryOOBOTPEmail:
			fallthrough
		case config.AuthenticationFlowAuthenticationSecondaryOOBOTPSMS:
			options = useAuthenticationOptionAddSecondaryOOBOTP(options, deps, branch.Authentication, authenticators)
		}
	}

	return options, nil
}

func getAuthenticationOptionsForReauth(ctx context.Context, deps *authflow.Dependencies, flows authflow.Flows, userID string, step *config.AuthenticationFlowReauthFlowStep) ([]AuthenticateOption, error) {
	options := []AuthenticateOption{}

	authenticators, err := deps.Authenticators.List(userID)
	if err != nil {
		return nil, err
	}

	checkHasAuthenticator := func(kind model.AuthenticatorKind, typ model.AuthenticatorType) bool {
		as := authenticator.ApplyFilters(
			authenticators,
			authenticator.KeepKind(kind),
			authenticator.KeepType(typ),
		)
		return len(as) > 0
	}

	useAuthenticationOptionAddPrimaryPassword := func(options []AuthenticateOption) []AuthenticateOption {
		if checkHasAuthenticator(
			model.AuthenticatorKindPrimary,
			model.AuthenticatorTypePassword,
		) {
			options = append(options, NewAuthenticateOptionPassword(
				config.AuthenticationFlowAuthenticationPrimaryPassword),
			)
		}
		return options
	}

	useAuthenticationOptionAddSecondaryPassword := func(options []AuthenticateOption) []AuthenticateOption {
		if checkHasAuthenticator(
			model.AuthenticatorKindSecondary,
			model.AuthenticatorTypePassword,
		) {
			options = append(options, NewAuthenticateOptionPassword(
				config.AuthenticationFlowAuthenticationSecondaryPassword),
			)
		}
		return options
	}

	useAuthenticationOptionAddTOTP := func(options []AuthenticateOption) []AuthenticateOption {
		if checkHasAuthenticator(
			model.AuthenticatorKindSecondary,
			model.AuthenticatorTypeTOTP,
		) {
			options = append(options, NewAuthenticateOptionTOTP())
		}
		return options
	}

	useAuthenticationOptionAddPasskey := func(options []AuthenticateOption) ([]AuthenticateOption, error) {
		if checkHasAuthenticator(
			model.AuthenticatorKindPrimary,
			model.AuthenticatorTypePasskey,
		) {
			requestOptions, err := deps.PasskeyRequestOptionsService.MakeModalRequestOptionsWithUser(userID)
			if err != nil {
				return nil, err
			}

			options = append(options, NewAuthenticateOptionPasskey(requestOptions))
		}

		return options, nil
	}

	useAuthenticationOptionAddOOBOTP := func(options []AuthenticateOption, authentication config.AuthenticationFlowAuthentication, kind model.AuthenticatorKind, typ model.AuthenticatorType) []AuthenticateOption {
		as := authenticator.ApplyFilters(
			authenticators,
			authenticator.KeepKind(kind),
			authenticator.KeepType(typ),
		)
		for _, info := range as {
			option, ok := NewAuthenticateOptionOOBOTPFromAuthenticator(deps.Config.Authenticator.OOB, info)
			if ok && option.Authentication == authentication {
				options = append(options, *option)
			}
		}
		return options
	}

	for _, branch := range step.OneOf {
		switch branch.Authentication {
		case config.AuthenticationFlowAuthenticationPrimaryPassword:
			options = useAuthenticationOptionAddPrimaryPassword(options)
		case config.AuthenticationFlowAuthenticationPrimaryPasskey:
			options, err = useAuthenticationOptionAddPasskey(options)
			if err != nil {
				return nil, err
			}
		case config.AuthenticationFlowAuthenticationSecondaryPassword:
			options = useAuthenticationOptionAddSecondaryPassword(options)
		case config.AuthenticationFlowAuthenticationSecondaryTOTP:
			options = useAuthenticationOptionAddTOTP(options)
		case config.AuthenticationFlowAuthenticationPrimaryOOBOTPEmail:
			options = useAuthenticationOptionAddOOBOTP(options, branch.Authentication, authenticator.KindPrimary, model.AuthenticatorTypeOOBEmail)
		case config.AuthenticationFlowAuthenticationPrimaryOOBOTPSMS:
			options = useAuthenticationOptionAddOOBOTP(options, branch.Authentication, authenticator.KindPrimary, model.AuthenticatorTypeOOBSMS)
		case config.AuthenticationFlowAuthenticationSecondaryOOBOTPEmail:
			options = useAuthenticationOptionAddOOBOTP(options, branch.Authentication, authenticator.KindSecondary, model.AuthenticatorTypeOOBEmail)
		case config.AuthenticationFlowAuthenticationSecondaryOOBOTPSMS:
			options = useAuthenticationOptionAddOOBOTP(options, branch.Authentication, authenticator.KindSecondary, model.AuthenticatorTypeOOBSMS)
		}
	}

	return options, nil
}

func identityFillDetails(err error, spec *identity.Spec, otherSpec *identity.Spec) error {
	details := errorutil.Details{}

	if spec != nil {
		details["IdentityTypeIncoming"] = apierrors.APIErrorDetail.Value(spec.Type)
		switch spec.Type {
		case model.IdentityTypeLoginID:
			details["LoginIDTypeIncoming"] = apierrors.APIErrorDetail.Value(spec.LoginID.Type)
		case model.IdentityTypeOAuth:
			details["OAuthProviderTypeIncoming"] = apierrors.APIErrorDetail.Value(spec.OAuth.ProviderID.Type)
		}
	}

	if otherSpec != nil {
		details["IdentityTypeExisting"] = apierrors.APIErrorDetail.Value(otherSpec.Type)
		switch otherSpec.Type {
		case model.IdentityTypeLoginID:
			details["LoginIDTypeExisting"] = apierrors.APIErrorDetail.Value(otherSpec.LoginID.Type)
		case model.IdentityTypeOAuth:
			details["OAuthProviderTypeExisting"] = apierrors.APIErrorDetail.Value(otherSpec.OAuth.ProviderID.Type)
		}
	}

	return errorutil.WithDetails(err, details)
}

func getChannels(claimName model.ClaimName, oobConfig *config.AuthenticatorOOBConfig) []model.AuthenticatorOOBChannel {
	channels := []model.AuthenticatorOOBChannel{}

	switch claimName {
	case model.ClaimEmail:
		channels = append(channels, model.AuthenticatorOOBChannelEmail)
	case model.ClaimPhoneNumber:
		switch oobConfig.SMS.PhoneOTPMode {
		case config.AuthenticatorPhoneOTPModeSMSOnly:
			channels = append(channels, model.AuthenticatorOOBChannelSMS)
		case config.AuthenticatorPhoneOTPModeWhatsappOnly:
			channels = append(channels, model.AuthenticatorOOBChannelWhatsapp)
		case config.AuthenticatorPhoneOTPModeWhatsappSMS:
			channels = append(channels, model.AuthenticatorOOBChannelWhatsapp)
			channels = append(channels, model.AuthenticatorOOBChannelSMS)
		}
	}

	return channels
}

func getOTPForm(purpose otp.Purpose, claimName model.ClaimName, cfg *config.AuthenticatorOOBEmailConfig) otp.Form {
	switch purpose {
	case otp.PurposeVerification:
		// Always use code.
		return otp.FormCode
	case otp.PurposeForgotPassword:
		// Always use link.
		return otp.FormLink
	case otp.PurposeOOBOTP:
		switch claimName {
		case model.ClaimEmail:
			if cfg.EmailOTPMode == config.AuthenticatorEmailOTPModeLoginLinkOnly {
				return otp.FormLink
			}
			return otp.FormCode
		case model.ClaimPhoneNumber:
			return otp.FormCode
		default:
			panic(fmt.Errorf("unexpected claim name: %v", claimName))
		}
	default:
		panic(fmt.Errorf("unexpected purpose: %v", purpose))
	}
}

func newIdentityInfo(deps *authflow.Dependencies, newUserID string, spec *identity.Spec) (*identity.Info, error) {
	// FIXME(authflow): allow bypassing email blocklist for Admin API.
	info, err := deps.Identities.New(newUserID, spec, identity.NewIdentityOptions{})
	if err != nil {
		return nil, err
	}

	duplicate, err := deps.Identities.CheckDuplicated(info)
	if err != nil && !errors.Is(err, identity.ErrIdentityAlreadyExists) {
		return nil, err
	}

	if err != nil {
		spec := info.ToSpec()
		otherSpec := duplicate.ToSpec()
		return nil, identityFillDetails(api.ErrDuplicatedIdentity, &spec, &otherSpec)
	}

	return info, nil
}

func findExactOneIdentityInfo(deps *authflow.Dependencies, spec *identity.Spec) (*identity.Info, error) {
	bucketSpec := AccountEnumerationPerIPRateLimitBucketSpec(
		deps.Config.Authentication,
		string(deps.RemoteIP),
	)

	reservation := deps.RateLimiter.Reserve(bucketSpec)
	err := reservation.Error()
	if err != nil {
		return nil, err
	}
	defer deps.RateLimiter.Cancel(reservation)

	exactMatch, otherMatches, err := deps.Identities.SearchBySpec(spec)
	if err != nil {
		return nil, err
	}

	if exactMatch == nil {
		// Consume the reservation if exact match is not found.
		reservation.Consume()

		var otherSpec *identity.Spec
		if len(otherMatches) > 0 {
			s := otherMatches[0].ToSpec()
			otherSpec = &s
		}
		return nil, identityFillDetails(api.ErrUserNotFound, spec, otherSpec)
	}

	return exactMatch, nil
}

type HandleOAuthAuthorizationResponseOptions struct {
	Alias       string
	RedirectURI string
}

func handleOAuthAuthorizationResponse(deps *authflow.Dependencies, opts HandleOAuthAuthorizationResponseOptions, inputOAuth inputTakeOAuthAuthorizationResponse) (*identity.Spec, error) {
	if oauthError := inputOAuth.GetOAuthError(); oauthError != "" {
		errorDescription := inputOAuth.GetOAuthErrorDescription()
		errorURI := inputOAuth.GetOAuthErrorURI()

		return nil, sso.NewOAuthError(oauthError, errorDescription, errorURI)
	}

	oauthProvider := deps.OAuthProviderFactory.NewOAuthProvider(opts.Alias)
	if oauthProvider == nil {
		return nil, api.ErrOAuthProviderNotFound
	}

	code := inputOAuth.GetOAuthAuthorizationCode()

	// TODO(authflow): support nonce but do not save nonce in cookies.
	// Nonce in the current implementation is stored in cookies.
	// In the Authentication Flow API, cookies are not sent in Safari in third-party context.
	emptyNonce := ""
	authInfo, err := oauthProvider.GetAuthInfo(
		sso.OAuthAuthorizationResponse{
			Code: code,
		},
		sso.GetAuthInfoParam{
			RedirectURI: opts.RedirectURI,
			Nonce:       emptyNonce,
		},
	)
	if err != nil {
		return nil, err
	}

	providerConfig := oauthProvider.Config()
	providerID := providerConfig.ProviderID()
	identitySpec := &identity.Spec{
		Type: model.IdentityTypeOAuth,
		OAuth: &identity.OAuthSpec{
			ProviderID:     providerID,
			SubjectID:      authInfo.ProviderUserID,
			RawProfile:     authInfo.ProviderRawProfile,
			StandardClaims: authInfo.StandardAttributes.ToClaims(),
		},
	}

	return identitySpec, nil
}

type GetOAuthDataOptions struct {
	RedirectURI  string
	Alias        string
	ResponseMode sso.ResponseMode
}

func getOAuthData(ctx context.Context, deps *authflow.Dependencies, opts GetOAuthDataOptions) (data OAuthData, err error) {
	oauthProvider := deps.OAuthProviderFactory.NewOAuthProvider(opts.Alias)
	if oauthProvider == nil {
		err = api.ErrOAuthProviderNotFound
		return
	}

	uiParam := uiparam.GetUIParam(ctx)

	param := sso.GetAuthURLParam{
		RedirectURI:  opts.RedirectURI,
		ResponseMode: opts.ResponseMode,
		Prompt:       uiParam.Prompt,
	}

	authorizationURL, err := oauthProvider.GetAuthURL(param)
	if err != nil {
		return
	}

	data = OAuthData{
		Alias:                 opts.Alias,
		OAuthProviderType:     oauthProvider.Config().Type,
		OAuthAuthorizationURL: authorizationURL,
		WechatAppType:         oauthProvider.Config().AppType,
	}
	return
}

func getMaskedOTPTarget(claimName model.ClaimName, claimValue string) string {
	switch claimName {
	case model.ClaimEmail:
		return mail.MaskAddress(claimValue)
	case model.ClaimPhoneNumber:
		return phone.Mask(claimValue)
	default:
		panic(fmt.Errorf("unexpected claim name: %v", claimName))
	}
}

func createAuthenticator(deps *authflow.Dependencies, userID string, authentication config.AuthenticationFlowAuthentication, target string) (*authenticator.Info, error) {
	spec := &authenticator.Spec{
		UserID: userID,
		OOBOTP: &authenticator.OOBOTPSpec{},
	}

	switch authentication {
	case config.AuthenticationFlowAuthenticationPrimaryOOBOTPEmail:
		spec.Kind = model.AuthenticatorKindPrimary
		spec.Type = model.AuthenticatorTypeOOBEmail
		spec.OOBOTP.Email = target

	case config.AuthenticationFlowAuthenticationPrimaryOOBOTPSMS:
		spec.Kind = model.AuthenticatorKindPrimary
		spec.Type = model.AuthenticatorTypeOOBSMS
		spec.OOBOTP.Phone = target

	case config.AuthenticationFlowAuthenticationSecondaryOOBOTPEmail:
		spec.Kind = model.AuthenticatorKindSecondary
		spec.Type = model.AuthenticatorTypeOOBEmail
		spec.OOBOTP.Email = target

	case config.AuthenticationFlowAuthenticationSecondaryOOBOTPSMS:
		spec.Kind = model.AuthenticatorKindSecondary
		spec.Type = model.AuthenticatorTypeOOBSMS
		spec.OOBOTP.Phone = target

	default:
		panic(fmt.Errorf("unexpected authentication method: %v", authentication))
	}

	isDefault, err := authenticatorIsDefault(deps, userID, spec.Kind)
	if err != nil {
		return nil, err
	}
	spec.IsDefault = isDefault

	authenticatorID := uuid.New()
	info, err := deps.Authenticators.NewWithAuthenticatorID(authenticatorID, spec)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func isNodeRestored(nodePointer jsonpointer.T, restoreTo jsonpointer.T) bool {
	return !authflow.JSONPointerSubtract(nodePointer, restoreTo).More()
}
