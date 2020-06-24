package oauth

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/skygeario/skygear-server/pkg/auth/dependency/oauth/handler"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/oauth/protocol"
	"github.com/skygeario/skygear-server/pkg/db"
	"github.com/skygeario/skygear-server/pkg/log"
)

func ConfigureAuthorizeHandler(router *mux.Router, h http.Handler) {
	router.NewRoute().
		Path("/oauth2/authorize").
		Methods("GET", "POST").
		Handler(h)
}

type AuthorizeHandlerLogger struct{ *log.Logger }

func NewAuthorizeHandlerLogger(lf *log.Factory) AuthorizeHandlerLogger {
	return AuthorizeHandlerLogger{lf.New("handler-authz")}
}

type ProtocolAuthorizeHandler interface {
	Handle(r protocol.AuthorizationRequest) handler.AuthorizationResult
}

var errAuthzInternalError = errors.New("internal error")

type AuthorizeHandler struct {
	Logger       AuthorizeHandlerLogger
	DBContext    db.Context
	AuthzHandler ProtocolAuthorizeHandler
}

func (h *AuthorizeHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}

	req := protocol.AuthorizationRequest{}
	for name, values := range r.Form {
		req[name] = values[0]
	}

	var result handler.AuthorizationResult
	err = db.WithTx(h.DBContext, func() error {
		result = h.AuthzHandler.Handle(req)
		if result.IsInternalError() {
			return errAuthzInternalError
		}
		return nil
	})

	if err == nil || errors.Is(err, errAuthzInternalError) {
		result.WriteResponse(rw, r)
	} else {
		h.Logger.WithError(err).Error("oauth authz handler failed")
		http.Error(rw, "Internal Server Error", 500)
	}
}
