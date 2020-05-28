package main

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
)

var runTest *bool = flag.Bool("run.test", false, "Should we run the tests")

var opt = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "progress", // can define default values
}

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opt)
}

func TestMain(m *testing.M) {
	flag.Parse()
	if !*runTest {
		return
	}

	format := "progress"
	for _, arg := range os.Args[1:] {
		fmt.Println(arg)
		if arg == "-test.v=true" { // go test transforms -v option
			format = "pretty"
			break
		}
	}

	status := godog.RunWithOptions("godog", func(s *godog.Suite) {
		FeatureContext(s)
	}, godog.Options{
		Format: format,
		Paths:  []string{"features"},
	})

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}

func theServerIsRunning() error {
	return godog.ErrPending
}

func iMakeARequestTo(arg1, arg2 string) error {
	return godog.ErrPending
}

func aListOfProductsShouldBeReturned() error {
	return godog.ErrPending
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^the server is running$`, theServerIsRunning)
	s.Step(`^I make a "([^"]*)" request to "([^"]*)"$`, iMakeARequestTo)
	s.Step(`^A list of products should be returned$`, aListOfProductsShouldBeReturned)
}
