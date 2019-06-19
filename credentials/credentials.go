package credentials

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
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
	OauthClientIdFile string
	OauthTokenFile    string
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
	if p.ApiToken == "" && !p.DoesApiTokenFileExist() {
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

func (p *Pivotal) DoesApiTokenFileExist() bool {
	return doesFileExist(p.getApiTokenFilepath())
}

func (p *Pivotal) CreateApiTokenFile() error {
	return nil
}

func NewGoogleSlides(logger *log.Logger) *GoogleSlides {
	return &GoogleSlides{
		newCredentials(logger),
		"google-slides-oauth-2.0-client-id.json",
		"google-slides-oauth-2.0-token.json",
	}
}

func (gs *GoogleSlides) GetOauthClientIdFilepath() string {
	return gs.Path + string(os.PathSeparator) + gs.OauthClientIdFile
}

func (gs *GoogleSlides) DoesOauthClientIdFilepathExist() bool {
	return doesFileExist(gs.GetOauthClientIdFilepath())
}
