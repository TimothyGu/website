package main

import (
    "github.com/go-martini/martini"
    "github.com/martini-contrib/sessions"
    "github.com/martini-contrib/oauth2"
    "github.com/martini-contrib/render"
    goauth2 "golang.org/x/oauth2"
)

func main() {
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

