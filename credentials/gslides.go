package credentials

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"log"
	"os"
)

func NewGoogleSlides(logger *log.Logger) *GoogleSlides {
	return &GoogleSlides{
		newCredentials(logger),
		nil,
		nil,
		"google-slides-oauth-2.0-client-id.json",
		"google-slides-oauth-2.0-token.json",
	}
}

func (gs *GoogleSlides) Init() error {
	if gs.Oauth2ClientIdFile == "" && !gs.doesOauth2ClientIdFilepathExist() {
		return errors.New("Google Slides credentials do not exist and have not been specified")
	}

	config, err := gs.getOauth2Config()

	if err != nil {
		return err
	}

	gs.Oauth2Config = config

	token, err := gs.getToken()
	gs.Oauth2Token = token

	return err
}

func (gs *GoogleSlides) getToken() (*oauth2.Token, error) {
	token, err := gs.getTokenFromFile()

	if err != nil {
		token, err = gs.generateToken()

		gs.saveToken(gs.getOauth2TokenFilepath(), token)
	}

	return token, err
}

func (gs *GoogleSlides) saveToken(path string, token *oauth2.Token) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer f.Close()

	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}

	json.NewEncoder(f).Encode(token)
}

func (gs *GoogleSlides) getTokenFromFile() (*oauth2.Token, error) {
	f, err := os.Open(gs.getOauth2TokenFilepath())

	defer f.Close()

	if err != nil {
		return nil, err
	}

	token := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(token)

	return token, err
}

func (gs *GoogleSlides) generateToken() (*oauth2.Token, error) {
	authUrl := gs.Oauth2Config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	openBrowser(authUrl)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	token, err := gs.Oauth2Config.Exchange(oauth2.NoContext, authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}

	return token, nil
}

func (gs *GoogleSlides) getOauth2Config() (*oauth2.Config, error) {
	b, err := ioutil.ReadFile(gs.getOauth2ClientIdFilePath())

	if err != nil {
		return nil, errors.New("Unable to read file " + gs.getOauth2TokenFilepath())
	}

	// If modifying these scopes, delete your previously saved client_secret.json
	config, err := google.ConfigFromJSON(
		b,
		"https://www.googleapis.com/auth/presentations.readonly",
		"https://www.googleapis.com/auth/drive",
	)

	if err != nil {
		return nil, errors.New("Unable to parse config from file " + gs.getOauth2TokenFilepath())
	}

	return config, nil
}

func (gs *GoogleSlides) getOauth2TokenFilepath() string {
	return gs.Path + string(os.PathSeparator) + gs.Oauth2TokenFile
}

func (gs *GoogleSlides) getOauth2ClientIdFilePath() string {
	return gs.Path + string(os.PathSeparator) + gs.Oauth2ClientIdFile
}

func (gs *GoogleSlides) doesOauth2ClientIdFilepathExist() bool {
	return doesFileExist(gs.getOauth2TokenFilepath())
}

func (gs *GoogleSlides) CopyOauth2ClientIdFile(src string) error {
	return copyFile(src, gs.getOauth2ClientIdFilePath())
}
