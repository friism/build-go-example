package main

import (
	"flag"
	"fmt"
	"log"
	"time"
	"os"
	"path/filepath"
	"net/url"
	"code.google.com/p/go-netrc/netrc"
	"github.com/cyberdelia/heroku-go/v3"
)

var (
	password = flag.String("apikey", "", "api key")
	appName  = flag.String("app", "", "app")
	repo     = flag.String("archive", "", "archive url")

	apiURL    = "https://api.heroku.com"
	netrcPath = filepath.Join(os.Getenv("HOME"), ".netrc")
)

func main() {
	log.SetFlags(0)
	flag.Parse()

  if *password == "" {
  	u, _ := url.Parse(apiURL)
  	_,netrcpass := getCreds(u)
  	heroku.DefaultTransport.Password = netrcpass
  } else {
  	heroku.DefaultTransport.Password = *password
  }

	h := heroku.NewService(heroku.DefaultClient)

	build, err := h.BuildCreate(*appName, heroku.BuildCreateOpts{
		SourceBlob: &struct {
			URL     *string `json:"url,omitempty"`
			Version *string `json:"version,omitempty"`
		}{
			URL: heroku.String(*repo),
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	for build.Status == "pending" {
		build, err = h.BuildInfo(*appName, build.ID)
		fmt.Print(".")
		time.Sleep(time.Second)
	}

	r, err := h.BuildResultInfo(*appName, build.ID)
	if err != nil {
		log.Fatal(err)
	}

	for _, line := range r.Lines {
		fmt.Print(line.Line)
	}
}

func getCreds(u *url.URL) (user, pass string) {

	m, err := netrc.FindMachine(netrcPath, u.Host)
	if err != nil {
		fmt.Printf("netrc error (%s): %v", u.Host, err)
	}

	return m.Login, m.Password
}
