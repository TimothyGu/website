package main

import (
    "github.com/go-martini/martini"
    "github.com/martini-contrib/sessions"
    "github.com/martini-contrib/oauth2"
    "github.com/martini-contrib/render"
    goauth2 "golang.org/x/oauth2"

    "log"
    "os"
    "encoding/json"

    "gopkg.in/mgo.v2"
)

var config struct {
    WebAddr string `json:"web_addr"`

    DBAddr string `json:"db_addr"`
    DBName string `json:"db_name"`
    DBUser string `json:"db_user"`
    DBPass string `json:"db_pass"`
}

var database *mgo.Database

func main() {
    cfg, err := os.Open("config.json")
    if err != nil {
        log.Fatal(err)
    }
    parser := json.NewDecoder(cfg)
    if err = parser.Decode(&config); err != nil {
        log.Fatal("Bad json")
    }

    session, err := mgo.DialWithInfo(&mgo.DialInfo{
        Addrs: []string{config.DBAddr},
        Database: config.DBName,
        Username: config.DBUser,
        Password: config.DBPass,
    })
    if err != nil {
        log.Fatal(err)
    }

    database = session.DB(config.DBName)
    log.Print("Database online")

    m := martini.Classic()
    m.Use(sessions.Sessions("semquery", sessions.NewCookieStore([]byte("secret"))))
    m.Use(oauth2.Github(
        &goauth2.Config{
            ClientID: "f918501c6b895e21252f",
            ClientSecret: "7850b01cee636e7449e29e9d425afaa912cf40b4",
            Scopes: []string {},
            RedirectURL: "",
        },
    ))
    m.Use(render.Renderer(render.Options {
        Layout: "layout",
    }))
    m.Use(RequestData)

    RegisterHandlers(m)

    m.Run()
}

