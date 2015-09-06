package main

import (
    "os"
    "log"
    "strings"
    "io/ioutil"
    "math/rand"
    "os/exec"
    "strconv"
    "net"
    "sync"
    "encoding/json"
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

    file1 := []string{"foo", "bar", "baz"}
    file2 := []string{"a", "b", "c"}
    files := [][]string{file1, file2}

    data["files"] = files

    path := "_repos/" + req.FormValue("repo")

    if _, err := os.Stat(path); os.IsNotExist(err) {
        data["indexed"] = false
    } else {
        data["indexed"] = true
    }

    data["query"] = req.FormValue("q")

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
        os.MkdirAll(path, 0777)
        c := exec.Command("git", "clone", "https://github.com/" + repo + ".git", path)
        c.Run()
        c.Wait()

        cmd := exec.Command("java", "-jar", "/Users/August/Code/projects/semquery/engine/target/engine-1.0-SNAPSHOT.jar", "index", path, repo)

        cmdReader, _ := cmd.StdoutPipe()

        scanner := bufio.NewScanner(cmdReader)

        go func() {
            cmd.Start()
            for scanner.Scan() {
                sockCli.websocket.WriteMessage(1, []byte(scanner.Text()))
            }
        }()
        cmd.Wait()
    }


    cmd := exec.Command("java", "-jar", "/Users/August/Code/projects/semquery/engine/target/engine-1.0-SNAPSHOT.jar", "query", query, repo)

    cmdReader, _ := cmd.StdoutPipe()

    scanner := bufio.NewScanner(cmdReader)

    go func() {
        cmd.Start()
        for scanner.Scan() {
            text := scanner.Text()
            parts := strings.Split(text, ",")
            if len(parts) == 1 {
                sockCli.websocket.WriteMessage(1, []byte("#" + parts[0]))
                continue
            }
            file := parts[0]
            src, _ := ioutil.ReadFile(file)
            start, _ := strconv.Atoi(parts[1])
            end, _ := strconv.Atoi(parts[2])
            lines := extractLines(string(src), start, end)
            j := map[string]interface{}{}
            for k, v := range lines {
                j[strconv.Itoa(k)] = v
            }
            jstr, _ := json.Marshal(j)
            sockCli.websocket.WriteMessage(1, []byte(jstr))
        }
    }()
    cmd.Wait()

    log.Print("DONE WITH INDEXING!")
}

func extractLines(src string, start int, end int) map[int]string {
    lines := map[int]string{}

    currentLine := 1
    lineStartPos := 0
    relativeStartPos := 0

    for i := 0; i < start; i++ {
        if src[i] == '\n' {
            currentLine++;
            lineStartPos = i + 1
        }
        if (i == start - 1) {
            relativeStartPos = i - lineStartPos + 1
        }
    }

    relativeEndPos := 0

    for i := start; i < len(src); i++ {
        if i == end {
            relativeEndPos = i - lineStartPos
        }
        if src[i] == '\n' || i == len(src) - 1 {
            sub := src[lineStartPos : i]
            lines[currentLine] = sub

            if len(lines) == 15 {
                return lines
            }

            currentLine += 1
            lineStartPos = i + 1

            if i >= end {
                break
            }
        }
    }

    log.Print(relativeStartPos, relativeEndPos)

    return lines
}
