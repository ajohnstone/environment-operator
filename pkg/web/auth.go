package web

import (
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/coreos/go-oidc/jose"
	oidc "github.com/coreos/go-oidc/oidc"
	"github.com/pearsontechnology/environment-operator/pkg/config"
)

type AuthClient struct {
	Client        *oidc.Client
	AllowedGroups []string
	Token         string
}

func NewAuthClient() (*AuthClient, error) {

	retval := &AuthClient{}

	// Handle AUTH_TOKEN_FILE
	if config.Env.TokenFile != "" {
		b, err := ioutil.ReadFile(config.Env.TokenFile)
		if err != nil {
			return nil, err
		}
		retval.Token = string(b)
		return retval, nil
	}

	provider, err := oidc.FetchProviderConfig(http.DefaultClient, config.Env.OIDCIssuerURL)
	if err != nil {
		return nil, err
	}

	clientCredentials := oidc.ClientCredentials{ID: config.Env.OIDCClientID}

	clientConfig := oidc.ClientConfig{
		ProviderConfig: provider,
		Credentials:    clientCredentials,
	}

	client, err := oidc.NewClient(clientConfig)
	if err != nil {
		return nil, err
	}

	retval.AllowedGroups = strings.Split(config.Env.OIDCAllowedGroups, ",")
	retval.Client = client

	return retval, nil

}

func (a *AuthClient) Authenticate(token string) bool {
	if a.Token != "" {
		return a.Token == token
	}

	jwt, err := jose.ParseJWT(token)
	if err != nil {
		log.Errorf("Error parsing JWT: %s", err.Error())
		return false
	}

	if err = a.Client.VerifyJWT(jwt); err != nil {
		log.Errorf("Error verifying JWT: %s", err.Error())
		return false
	}

	claims, err := jwt.Claims()
	if err != nil {
		log.Errorf("Error getting claims from JWT: %s", err.Error())
		return false
	}

	log.Debugf("Token claims: %+v", claims)

	groups := claims["groups"].([]interface{})
	if len(groups) == 0 {
		log.Errorf("Error getting groups from JWT")
		return false
	}

	return a.allowsGroup(groups)
}

func (a *AuthClient) allowsGroup(groups []interface{}) bool {

	for _, g1 := range a.AllowedGroups {
		for _, g2 := range groups {
			log.Debugf("allowsGroup g1: %s, g2: %s", g1, g2.(string))
			if g1 == g2.(string) {
				return true
			}
		}
	}
	return false
}
