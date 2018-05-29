package pivex

import (
	"os/user"
	"fmt"
	"os"
	"log"
)

var (
	Logger   = log.New(os.Stdout, "logger: ", log.Lshortfile)
	ApiCreds = func() string {
		usr, err := user.Current()
		if err != nil {
			Logger.Fatal(err)
			panic(err)
		}

		return fmt.Sprintf("%s/.google-api", usr.HomeDir)
	}()
)
