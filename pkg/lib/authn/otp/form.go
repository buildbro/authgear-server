package otp

import (
	"github.com/authgear/authgear-server/pkg/lib/config"
	"github.com/authgear/authgear-server/pkg/util/secretcode"
)

type Form string

const (
	FormCode Form = "code"
	FormLink Form = "link"
)

func (f Form) AllowLookupByCode() bool {
	return f == FormLink
}

func (f Form) codeType() secretCode {
	switch f {
	case FormCode:
		return secretcode.OOBOTPSecretCode
	case FormLink:
		return secretcode.LinkOTPSecretCode
	default:
		panic("unexpected form: " + f)
	}
}

func (f Form) CodeLength() int {
	return f.codeType().Length()
}

func (f Form) GenerateCode(cfg *config.TestModeFeatureConfig, userID string) string {
	codeType := f.codeType()
	switch c := codeType.(type) {
	case secretcode.OOBOTPSecretCodeType:
		if cfg.FixedOOBOTP.Enabled {
			return c.GenerateFixed(cfg.FixedOOBOTP.Code)
		} else {
			return c.Generate()
		}
	case secretcode.LinkOTPSecretCodeType:
		if cfg.DeterministicLinkOTP.Enabled {
			return c.GenerateDeterministic(userID)
		} else {
			return c.Generate()
		}
	}
	panic("unknown otp form")
}

func (f Form) VerifyCode(input string, expected string) bool {
	return f.codeType().Compare(input, expected)
}

type secretCode interface {
	Length() int
	Compare(string, string) bool
}
