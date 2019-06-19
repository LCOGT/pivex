package credentials

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"log"
	"os"
	"os/user"
)

var (
	credentialsPath = func() (credsPath string) {
		usr, err := user.Current()
		credsPath = fmt.Sprintf("%s/.pivex", usr.HomeDir)

		if _, err := os.Stat(credsPath); os.IsNotExist(err) {
			os.Mkdir(credsPath, 0700)
		}

		if err != nil {
			// logger.Fatal(err)
		}

		return
	}()
)

type credentials struct {
	Path   string
	logger *log.Logger
}

type Pivotal struct {
	*credentials
	ApiToken     string
	ApiTokenFile string
}

type GoogleSlides struct {
	*credentials
	Oauth2Config       *oauth2.Config
	Oauth2Token        *oauth2.Token
	Oauth2ClientIdFile string
	Oauth2TokenFile    string
}

func doesFileExist(filepath string) bool {
	info, err := os.Stat(filepath)

	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}
func readFile(filePath string) (apiToken string) {
	f, err := os.Open(filePath)

	defer f.Close()

	if err != nil {
		log.Fatalf("Unable to read token file %s: %v", filePath, err)
	}

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		apiToken = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return
}

func copyFile(sourceFile string, destinationFile string) error {
	input, err := ioutil.ReadFile(sourceFile)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(destinationFile, input, 0644)

	if err != nil {
		return err
	}

	return nil
}

func writeFile(content string, filepath string) error {
	b := []byte(content + "\n")

	return ioutil.WriteFile(filepath, b, 0644)
}

func getCredentialsPath() *string {
	return &credentialsPath
}

func newCredentials(logger *log.Logger) *credentials {
	return &credentials{
		Path: *getCredentialsPath(),
	}
}

func NewPivotal(logger *log.Logger) *Pivotal {
	return &Pivotal{
		newCredentials(logger),
		"",
		"pivotal-api-token",
	}
}

func (p *Pivotal) Init() error {
	if p.ApiToken == "" && !p.doesApiTokenFileExist() {
		return errors.New("Pivotal credentials do not exist and have not been specified")
	} else if p.ApiToken != "" {
		return p.writeApiTokenToFile()
	} else {
		p.ApiToken = p.getApiTokenFromFile()

		return nil
	}
}

func (p *Pivotal) writeApiTokenToFile() error {
	return writeFile(p.ApiToken, p.getApiTokenFilepath())
}

func (p *Pivotal) getApiTokenFromFile() string {
	return readFile(p.getApiTokenFilepath())
}

func (p *Pivotal) getApiTokenFilepath() string {
	return p.Path + string(os.PathSeparator) + p.ApiTokenFile
}

func (p *Pivotal) doesApiTokenFileExist() bool {
	return doesFileExist(p.getApiTokenFilepath())
}

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
	authURL := gs.Oauth2Config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf(
		"Go to the following link in your browser and copy the authorization code, then paste it in the line "+
			"below\n%v\nAuthorization code: ",
		authURL)

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
