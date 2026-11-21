package matrixhelper

import (
	"errors"
	"fmt"
	"regexp"
)

var clientSecretRegex = regexp.MustCompile(`\bclient_secret:[ \t]*(?:"([^"]*)"|([^\s"#]+))`)

func ensureOidcClientSecret(masConfig, sourceClientSecret string) (string, bool, error) {
	if err := validateClientSecret(sourceClientSecret); err != nil {
		return "", false, err
	}
	switch n := len(clientSecretRegex.FindAllString(masConfig, -1)); n {
	case 0:
		return "", false, errors.New("client_secret not found in user-keycloak config")
	case 1:
	default:
		return "", false, fmt.Errorf("client_secret appears %d times in user-keycloak config", n)
	}
	m := clientSecretRegex.FindStringSubmatch(masConfig)
	current := m[1]
	if current == "" {
		current = m[2]
	}
	if current == sourceClientSecret {
		return masConfig, false, nil
	}
	desired := clientSecretRegex.ReplaceAllStringFunc(masConfig, func(string) string {
		return `client_secret: "` + sourceClientSecret + `"`
	})
	return desired, true, nil
}

func validateClientSecret(clientSecret string) error {
	var validNewValueRegex = regexp.MustCompile(`^[^\s"]+$`)
	if clientSecret == "" {
		return errors.New("client secret is empty")
	}
	if !validNewValueRegex.MatchString(clientSecret) {
		return fmt.Errorf("client secret is outside the supported character set (no whitespace, no double quotes)")
	}
	return nil
}
