package main

import (
    "os"
    "strings"
    "os/exec"
    "strconv"
    "net/http"
    "github.com/go-martini/martini"
    "github.com/martini-contrib/sessions"
    "github.com/martini-contrib/oauth2"
    "github.com/martini-contrib/render"
)

func RegisterHandlers(m *martini.ClassicMartini) {
    m.Get("/", RootPage)
    m.Get("/repo", CacheRepository)
}

//Rendering home page with template data
func RootPage(tokens oauth2.Tokens, session sessions.Session, r render.Render) {
    data := map[string]string {
        "loggedin": strconv.FormatBool(!tokens.Expired()),
    }
    if !tokens.Expired() {
        data["username"] = session.Get("username").(string)
        data["avatar"] = session.Get("avatar").(string)
    }
    r.HTML(200, "index", data)
}

//Rendering search page with template data
func SearchPage(tokens oauth2.Tokens, session sessions.Session, r render.Render, req *http.Request) {
}

//Retrieves github repository to prepare to be indexed and searched 
func CacheRepository(tokens oauth2.Tokens, session sessions.Session, req *http.Request, w http.ResponseWriter) {
    if !tokens.Expired() && session.Get("username") != nil {
        query := req.URL.Query().Get("query")
        if query != "" {
            if _, err := os.Stat(strings.Split(query, "/")[1]); os.IsNotExist(err) {
                exec.Command("git", "clone", "https://github.com/" + query + ".git").Run()
            }
        }
    } else {
        http.Redirect(w, req, "/", 302)
    }
}

