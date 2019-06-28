package password

import (
	"github.com/sirupsen/logrus"
	"github.com/skygeario/skygear-server/pkg/core/db"
)

type safeProviderImpl struct {
	impl      *providerImpl
	txContext db.SafeTxContext
}

func NewSafeProvider(
	builder db.SQLBuilder,
	executor db.SQLExecutor,
	logger *logrus.Entry,
	loginIDsKeyWhitelist []string,
	allowedRealms []string,
	passwordHistoryEnabled bool,
	txContext db.SafeTxContext,
) Provider {
	return &safeProviderImpl{
		impl:      newProvider(builder, executor, logger, loginIDsKeyWhitelist, allowedRealms, passwordHistoryEnabled),
		txContext: txContext,
	}
}

func (p *safeProviderImpl) IsLoginIDValid(loginID map[string]string) bool {
	p.txContext.EnsureTx()
	return p.impl.IsLoginIDValid(loginID)
}

func (p safeProviderImpl) IsRealmValid(realm string) bool {
	p.txContext.EnsureTx()
	return p.impl.IsRealmValid(realm)
}

func (p *safeProviderImpl) IsDefaultAllowedRealms() bool {
	p.txContext.EnsureTx()
	return p.impl.IsDefaultAllowedRealms()
}

func (p *safeProviderImpl) CreatePrincipalsByLoginID(authInfoID string, password string, loginID map[string]string, realm string) error {
	p.txContext.EnsureTx()
	return p.impl.CreatePrincipalsByLoginID(authInfoID, password, loginID, realm)
}

func (p *safeProviderImpl) CreatePrincipal(principal Principal) error {
	p.txContext.EnsureTx()
	return p.impl.CreatePrincipal(principal)
}

func (p *safeProviderImpl) GetPrincipalByLoginIDWithRealm(loginIDKey string, loginID string, realm string, principal *Principal) (err error) {
	p.txContext.EnsureTx()
	return p.impl.GetPrincipalByLoginIDWithRealm(loginIDKey, loginID, realm, principal)
}

func (p *safeProviderImpl) GetPrincipalsByUserID(userID string) ([]*Principal, error) {
	p.txContext.EnsureTx()
	return p.impl.GetPrincipalsByUserID(userID)
}

func (p *safeProviderImpl) GetPrincipalsByLoginID(loginIDKey string, loginID string) ([]*Principal, error) {
	p.txContext.EnsureTx()
	return p.impl.GetPrincipalsByLoginID(loginIDKey, loginID)
}

func (p *safeProviderImpl) UpdatePrincipal(principal Principal) error {
	p.txContext.EnsureTx()
	return p.impl.UpdatePrincipal(principal)
}
