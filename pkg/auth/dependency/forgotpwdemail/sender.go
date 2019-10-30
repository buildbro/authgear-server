package forgotpwdemail

import (
	"net/url"
	"path"
	"strconv"
	"time"

	"github.com/skygeario/skygear-server/pkg/auth/model"
	authTemplate "github.com/skygeario/skygear-server/pkg/auth/template"
	"github.com/skygeario/skygear-server/pkg/core/auth/authinfo"
	"github.com/skygeario/skygear-server/pkg/core/config"
	"github.com/skygeario/skygear-server/pkg/core/errors"
	"github.com/skygeario/skygear-server/pkg/core/mail"
	"github.com/skygeario/skygear-server/pkg/core/template"
)

type Sender interface {
	Send(
		email string,
		authInfo authinfo.AuthInfo,
		user model.User,
		hashedPassword []byte,
	) error
}

type DefaultSender struct {
	Config         config.ForgotPasswordConfiguration
	URLPrefix      *url.URL
	Sender         mail.Sender
	CodeGenerator  *CodeGenerator
	TemplateEngine *template.Engine
}

func NewDefaultSender(
	config config.TenantConfiguration,
	urlPrefix *url.URL,
	sender mail.Sender,
	templateEngine *template.Engine,
) Sender {
	return &DefaultSender{
		Config:         config.UserConfig.ForgotPassword,
		URLPrefix:      urlPrefix,
		Sender:         sender,
		CodeGenerator:  &CodeGenerator{config.UserConfig.MasterKey},
		TemplateEngine: templateEngine,
	}
}

func (d *DefaultSender) Send(
	email string,
	authInfo authinfo.AuthInfo,
	user model.User,
	hashedPassword []byte,
) (err error) {
	expireAt :=
		time.Now().UTC().
			Truncate(time.Second * 1).
			Add(time.Second * time.Duration(d.Config.ResetURLLifetime))
	code := d.CodeGenerator.Generate(authInfo, hashedPassword, expireAt)
	link := *d.URLPrefix
	link.Path = path.Join(link.Path, "_auth/forgot_password/reset_password_form")
	link.RawQuery = url.Values{
		"code":      []string{code},
		"user_id":   []string{authInfo.ID},
		"expire_at": []string{strconv.FormatInt(expireAt.UTC().Unix(), 10)},
	}.Encode()
	context := map[string]interface{}{
		"appname":    d.Config.AppName,
		"link":       link.String(),
		"email":      email,
		"user":       user,
		"url_prefix": d.URLPrefix.String(),
		"code":       code,
		"expire_at":  expireAt,
	}

	var textBody string
	if textBody, err = d.TemplateEngine.ParseTextTemplate(
		string(authTemplate.TemplateItemTypeForgotPasswordEmailTXT),
		context,
		template.ParseOption{Required: true},
	); err != nil {
		err = errors.Newf("failed to render forgot password text email: %w", err)
		return
	}

	var htmlBody string
	if htmlBody, err = d.TemplateEngine.ParseHTMLTemplate(
		string(authTemplate.TemplateItemTypeForgotPasswordEmailHTML),
		context,
		template.ParseOption{Required: false},
	); err != nil {
		err = errors.Newf("failed to render forgot password HTML email: %w", err)
		return
	}

	err = d.Sender.Send(mail.SendOptions{
		Sender:    d.Config.Sender,
		Recipient: email,
		Subject:   d.Config.Subject,
		ReplyTo:   d.Config.ReplyTo,
		TextBody:  textBody,
		HTMLBody:  htmlBody,
	})
	if err != nil {
		err = errors.Newf("failed to send forgot password email: %w", err)
	}

	return
}
