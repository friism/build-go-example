package main

import (
  "flag"
  "fmt"
  "log"
  "time"

  "github.com/cyberdelia/heroku-go/v3"
)

var (
  username = flag.String("username", "", "api username")
  password = flag.String("password", "", "api password")
  appName = flag.String("app", "", "app")
  repo = flag.String("archive", "", "archive url")
)

func main() {
  log.SetFlags(0)
  flag.Parse()

  heroku.DefaultTransport.Username = *username
  heroku.DefaultTransport.Password = *password

  h := heroku.NewService(heroku.DefaultClient)

  build, err := h.BuildCreate(*appName, heroku.BuildCreateOpts{
    SourceBlob: &struct {
      URL     *string `json:"url,omitempty"`
      Version *string `json:"version,omitempty"`
    }{
      URL:     heroku.String(*repo),
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
    fmt.Println(line.Line)
  }
}