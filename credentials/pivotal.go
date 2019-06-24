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
	if !p.doesApiTokenFileExist() {
		return errors.New("Pivotal credentials do not exist and have not been specified")
	} else {
		p.ApiToken = p.getApiTokenFromFile()

		return nil
	}
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

func (p *Pivotal) CopyApiTokenFile(src string) error {
	return copyFile(src, p.getApiTokenFilepath())
}
