package main

import (
	"pivex/pivotal"
	"pivex/export"
)

func main() {
	piv := pivotal.New()
	piv.GetStories()

	gs := gslides.New()
	gs.Export(&piv.Intervals[0])
}
