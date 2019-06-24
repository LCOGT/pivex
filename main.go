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
	"strconv"
)

var (
	logger = log.New(os.Stdout, "logger: ", log.Lshortfile)
)

func main() {
	box := packr.New("version", "./static")

	pivCreds := credentials.NewPivotal(logger)
	pivApiTokenFile := flag.StringP(
		"pivotal-api-token-file",
		"p",
		"",
		fmt.Sprintf(
			"exports the Pivotal API token from the given file, this token will be stored under %s and overwrite any existing token",
			pivCreds.Path))

	gSlideCreds := credentials.NewGoogleSlides(logger)
	oauthClientIdFile := flag.StringP(
		"google-client-id-file",
		"g",
		"",
		fmt.Sprintf("exports the given Google API token JSON, this token will be stored under %s and overwrite any existing token",
			gSlideCreds.Path))

	deckName := flag.StringP(
		"deck-name",
		"d",
		"",
		getDefaultDeckName("[current sprint]"))

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

	if *pivApiTokenFile != "" {
		err := pivCreds.CopyApiTokenFile(*pivApiTokenFile)

		if err != nil {
			logger.Panicln(err)
		}
	}

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

	gsOpts := export.Opts{
		ForceCreate: *fCreate,
	}

	if *deckName == "" {
		gsOpts.DeckName = getDefaultDeckName(strconv.Itoa(piv.Iterations[0].Number))
	} else {
		gsOpts.DeckName = *deckName
	}

	gs := export.New(gSlideCreds, &gsOpts, logger, piv.Iterations[0])

	gs.Export()
}

func getDefaultDeckName(sprint string) string {
	return fmt.Sprintf("Sprint Demo %s", sprint)
}
