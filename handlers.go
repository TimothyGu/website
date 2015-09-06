package main

import (
    "os"
    "log"
    "strings"
    "math/rand"
    "os/exec"
    "strconv"
    "net"
    "sync"
    "bufio"
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
    m.Get("/socket", SocketPage)
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

var ws_transfer = map[int64][]string{}

//Rendering search page with template data
func QueryPage(tokens oauth2.Tokens, session sessions.Session, r render.Render, req *http.Request) {
    data := CreateData(tokens, session)

    req.ParseForm()
    id := rand.Int63()
    data["ws_id"] = id
    ws_transfer[id] = []string{req.FormValue("q"), req.FormValue("repo")}
    r.HTML(200, "query", data)
}

var ActiveClients = map[ClientConn]int {}
var ActiveClientsRWMutex sync.RWMutex

type ClientConn struct {
    websocket *websocket.Conn
    clientIP net.Addr
}

func SocketPage(tokens oauth2.Tokens, r *http.Request, w http.ResponseWriter) {
    ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
    if _, ok := err.(websocket.HandshakeError); (ok || err != nil) {
        log.Fatal(err)
        return
    }
    //Initial connection, store 
    client := ws.RemoteAddr()
    sockCli := ClientConn {ws, client}
    ActiveClientsRWMutex.Lock()
    ActiveClients[sockCli] = 0
    ActiveClientsRWMutex.Unlock()

    log.Print("Starting")
    _, msg, err := sockCli.websocket.ReadMessage()
    if err != nil {
        log.Print(err)
        return
    }

    id, err := strconv.ParseInt(string(msg), 10, 64)
    if err != nil {
        log.Print(err)
        return
    }

    arr := ws_transfer[id]
    query := arr[0]
    repo := arr[1]
    delete(ws_transfer, id)

    repo_parts := strings.Split(repo, "/")
    if len(repo_parts) < 2 {
        log.Print("Invalid repo")
        return
    }

    path := "_repos/" + repo

    if _, err := os.Stat(path); os.IsNotExist(err) {
        os.MkdirAll(path, os.ModeDir)
        c := exec.Command("git", "clone", "https://github.com/" + repo + ".git", path)
        c.Run()
        c.Wait()
    }

    log.Print("query: " + query)
    cmd := exec.Command("java", "-jar", "/Users/August/Code/projects/semquery/engine/target/engine-1.0-SNAPSHOT.jar", "index", "/Users/August/Documents/binnavi", repo)

    cmdReader, err := cmd.StdoutPipe()

    scanner := bufio.NewScanner(cmdReader)
    go func() {
        cmd.Start()

        for scanner.Scan() {
            sockCli.websocket.WriteMessage(1, []byte(scanner.Text()))
        }

        defer cmd.Wait()
    }()

}
