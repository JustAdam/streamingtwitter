package streamingtwitter

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/go-oauth/oauth"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	// File permissions for the token file.
	tokenFilePermission = os.FileMode(0600)
)

type Tokener interface {
	// Token returns a valid user access token to provide access to Twitter.
	// This method also needs to set the app token so valid requests can be made.
	Token(*oauth.Client) (*oauth.Credentials, error)
}

type ClientTokens struct {
	// Location to our token storage file (JSON format)
	TokenFile string `json:"-"`
	// Token for the actual application
	App *oauth.Credentials
	// Token for the user of the application
	User *oauth.Credentials
}

type ClientTokensError struct {
	Msg string
}

func (e ClientTokensError) Error() string {
	return e.Msg
}

// You get a token for your App from Twitter.  Put this within the App section
// of the  JSON token file.  The user's token will be requested, then written
// and saved to this file.
func (t *ClientTokens) Token(oc *oauth.Client) (*oauth.Credentials, error) {
	if t.TokenFile == "" {
		return nil, &ClientTokensError{
			Msg: "no token file supplied",
		}
	}

	cf, err := ioutil.ReadFile(t.TokenFile)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(cf, t); err != nil {
		return nil, err
	}

	if t.App == nil {
		return nil, &ClientTokensError{
			Msg: "missing \"App\" token",
		}
	}

	if t.App.Token == "" || t.App.Secret == "" {
		return nil, &ClientTokensError{
			Msg: "missing app's Token or Secret",
		}
	}
	oc.Credentials = *t.App

	var token *oauth.Credentials
	if t.User == nil {
		token = &oauth.Credentials{}
	} else {
		token = t.User
	}

	if token.Token == "" || token.Secret == "" {
		tempCredentials, err := oc.RequestTemporaryCredentials(http.DefaultClient, "oob", nil)
		if err != nil {
			return nil, err
		}

		url := oc.AuthorizationURL(tempCredentials, nil)
		fmt.Fprintf(os.Stdout, "Before we can continue ...\nGo to:\n\n\t%s\n\nAuthorize the application and enter in the verification code: ", url)

		var authCode string
		fmt.Scanln(&authCode)

		token, _, err = oc.RequestToken(http.DefaultClient, tempCredentials, authCode)
		if err != nil {
			return nil, err
		}

		// Save the user token within our token file
		t.User = token
		save, err := json.Marshal(t)
		if err != nil {
			return nil, err
		}

		if err := ioutil.WriteFile(t.TokenFile, save, tokenFilePermission); err != nil {
			return nil, err
		}
	}

	return token, nil
}
