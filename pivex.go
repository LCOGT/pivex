package main

import (
	"pivex/pivotal"
	"pivex/export"
	"flag"
)

func main() {
	delAuth := flag.Bool("d", false, "delete the authentication files being used for pivex")
	flag.Parse()
	piv := pivotal.New()


	piv.GetStories()

	gs := gslides.New()

	if *delAuth {
		gs.DelAuth()
	}

	gs.Export(&piv.Intervals[0])
}
