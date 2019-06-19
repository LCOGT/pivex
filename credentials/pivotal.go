package credentials

import (
	"errors"
	"log"
	"os"
)

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
