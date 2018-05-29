package main

import (
	"pivex/pivex/pivotal"
	"pivex/pivex/export"
	"flag"
	"fmt"
	"pivex/pivex"
)


func main() {
	delAuth := flag.Bool("d", false, "delete the authentication files being used for pivex")
	pivApiTok := flag.String("p", "", fmt.Sprintf("the Pivotal API token to be used, this token only needs to specified when the token is set for the firts time and will be stored under %s after the first time it is set", pivex.ApiCreds))
	flag.Parse()
	piv := pivotal.New()


	piv.GetStories()

	gs := export.New()

	if *delAuth {
		gs.DelAuth()
	}

	gs.Export(&piv.Intervals[0])
}
