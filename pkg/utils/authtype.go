package utils

import (
	"errors"
	"strings"

	sf "github.com/snowflakedb/gosnowflake"
)

// available auth type (case insensitive) are
//   - snowflake
//   - oauth
//   - externalbrowser
//   - okta
//   - snowflake_jwt
//   - tokenaccessor
//   - username_password_mfa
func GetSnowflakeAuthType(snowflakeAuthenticator string) (sf.AuthType, error) {
	switch strings.ToLower(snowflakeAuthenticator) {
	case "snowflake":
		return sf.AuthTypeSnowflake, nil
	case "oauth":
		return sf.AuthTypeOAuth, nil
	case "externalbrowser":
		return sf.AuthTypeExternalBrowser, nil
	case "okta":
		return sf.AuthTypeOkta, nil
	case "snowflake_jwt":
		return sf.AuthTypeJwt, nil
	case "tokenaccessor":
		return sf.AuthTypeTokenAccessor, nil
	case "username_password_mfa":
		return sf.AuthTypeUsernamePasswordMFA, nil
	default:
		return -1, errors.New("Unkown snowflake authenticator")
	}
}
