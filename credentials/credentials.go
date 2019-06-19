package credentials

import (
	"bufio"
	"fmt"
	"golang.org/x/oauth2"
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
