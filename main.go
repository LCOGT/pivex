package main

import (
	"fmt"
	"github.com/LCOGT/pivex/credentials"
	"github.com/LCOGT/pivex/export"
	"github.com/LCOGT/pivex/pivotal"
	"github.com/gobuffalo/packr/v2"
	flag "github.com/spf13/pflag"
	"log"
	"os"
)

var (
	logger = log.New(os.Stdout, "logger: ", log.Lshortfile)
)

func main() {
	box := packr.New("version", "./static")

	pivCreds := credentials.NewPivotal(logger)
	pivApiToken := flag.StringP(
		"pivotal-api-token",
		"p",
		"",
		fmt.Sprintf(
			"exports the given Pivotal API token, this token will be stored under %s and overwrite any existing token",
			pivCreds.Path))

	gSlideCreds := credentials.NewGoogleSlides(logger)
	oauthClientIdFile := flag.StringP(
		"google-token",
		"g",
		"",
		fmt.Sprintf("exports the given Google API token JSON, this token will be stored under %s and overwrite any existing token",
			gSlideCreds.Path))
	print(oauthClientIdFile)

	fCreate := flag.BoolP("force", "f", false, "overwrite an existing presentation")
	showVer := flag.BoolP("version", "v", false, "show the current version")

	flag.Parse()

	if *showVer {
		ver, err := box.FindString("version")

		if err != nil {
			print("Error getting version")
		}

		print("Version: " + ver)

		os.Exit(0)
	}

	pivCreds.ApiToken = *pivApiToken

	pivCredsErr := pivCreds.Init()

	if pivCredsErr != nil {
		logger.Panicln(pivCredsErr)
	}

	piv := pivotal.New(pivCreds, logger)
	piv.GetIterations()

	if *oauthClientIdFile != "" {
		err := gSlideCreds.CopyOauth2ClientIdFile(*oauthClientIdFile)

		if err != nil {
			logger.Panicln(err)
		}
	}

	gSlideCredsErr := gSlideCreds.Init()

	if gSlideCredsErr != nil {
		logger.Panicln(gSlideCredsErr)
	}

	gs := export.New(gSlideCreds, *fCreate, logger, piv.Iterations[0])

	gs.Export()
}
