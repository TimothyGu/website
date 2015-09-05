package main

import (
    "os"
    "strings"
    "os/exec"
    "strconv"
    "net/http"
    "github.com/go-martini/martini"
    "github.com/gorilla/websocket"
    "github.com/martini-contrib/sessions"
    "github.com/martini-contrib/oauth2"
    "github.com/martini-contrib/render"
)

func RegisterHandlers(m *martini.ClassicMartini) {
    m.Get("/", RootPage)
    m.Get("/repo", CacheRepository)

    m.Post("/query", QueryPage)
}

func CreateData(tokens oauth2.Tokens, session sessions.Session) map[string]interface{} {
    data := map[string]interface{} {
        "loggedin": strconv.FormatBool(!tokens.Expired()),
    }
    if !tokens.Expired() {
        data["username"] = session.Get("username").(string)
        data["avatar"] = session.Get("avatar").(string)
    }
    return data
}

//Rendering home page with template data
func RootPage(tokens oauth2.Tokens, session sessions.Session, r render.Render) {
    data := CreateData(tokens, session)
    r.HTML(200, "index", data)
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

//Rendering search page with template data
func QueryPage(tokens oauth2.Tokens, session sessions.Session, r render.Render, req *http.Request) {
    data := CreateData(tokens, session)
    r.HTML(200, "query", data)
}

var ActiveClients = map[ClientConn]int {}
var ActiveClientsRWMutex sync.RWMutex

type ClientConn struct {
    websocket *websocket.Conn
    clientIP net.Addr
}

func SocketPage(tokens oauth2.Tokens, r *http.Request, w http.ResponseWriter) {
    ws, err := websocket.Upgrade(w, r, 1024, 1024)
    if _, ok := err.(websocketHandshakeError; (ok || err != nil) {
        return
    }
    //Initial connection, store 
    client := ws.RemoteAddr()
    sockCli := ClientConn {ws, client}
    ActiveClientsRWMutex.Lock()
    ActiveClients[sockCli] = 0
    ActiveClientsRWMutex.Unlock()

    // Read incoming messages
    for {
        _, p, err := ws.ReadMessage()
        if err != nil {
            ActiveClientsRWMutex.Lock()
            delete(ActiveClients, sockCli)
            ActiveClientsRWMutex.Unlock()
            return
        }
    }
}
