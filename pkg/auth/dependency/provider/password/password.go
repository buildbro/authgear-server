package password

const providerPassword string = "password"

type Provider interface {
	ValidateLoginIDs(loginIDs []LoginID) error
	IsRealmValid(realm string) bool
	IsDefaultAllowedRealms() bool
	CreatePrincipalsByLoginID(authInfoID string, password string, loginIDs []LoginID, realm string) error
	CreatePrincipal(principal Principal) error
	GetPrincipalByLoginIDWithRealm(loginIDKey string, loginID string, realm string, principal *Principal) (err error)
	GetPrincipalsByUserID(userID string) ([]*Principal, error)
	GetPrincipalsByLoginID(loginIDKey string, loginID string) ([]*Principal, error)
	UpdatePrincipal(principal Principal) error
}
