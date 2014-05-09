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
)

func main() {
  log.SetFlags(0)
  flag.Parse()

  heroku.DefaultTransport.Username = *username
  heroku.DefaultTransport.Password = *password

  h := heroku.NewService(heroku.DefaultClient)

  a := "testing-build-api"
  u := "https://github.com/heroku/node-js-sample/archive/master.tar.gz"
  v := "foo"

  s := struct {
    URL     *string `json:"url,omitempty"`
    Version *string `json:"version,omitempty"`
  }{
    &u,
    &v,
  }

  build, err := h.BuildCreate(a, heroku.BuildCreateOpts{
    &s,
  })
  if err != nil {
    log.Fatal(err)
  }

  for build.Status == "pending" {
    build, err = h.BuildInfo(a, build.ID)
    fmt.Print(".")
    time.Sleep(time.Second)
  }

  r, err := h.BuildResultInfo(a, build.ID)
  if err != nil {
    log.Fatal(err)
  }

  for _, line := range r.Lines {
    fmt.Println(line.Line)
  }

  fmt.Print(r.Build.ID)
}