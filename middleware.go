package main

import (
    "log"
    "encoding/json"
    "net/http"
    "io/ioutil"
    "github.com/martini-contrib/oauth2"
    "github.com/martini-contrib/sessions"
)

/* Middleware to:
   - Retrieve data from github if user is logged in and data is not stored in sessions
   - Delete session data if user logged out */
func RequestData(tokens oauth2.Tokens, session sessions.Session) {
    if !tokens.Expired() && session.Get("username") == nil {
        access := tokens.Access()
        req, _ := http.NewRequest("GET", "https://api.github.com/user?access_token=" + access, nil)
        client := &http.Client {}
        resp, _ := client.Do(req)

        body, _ := ioutil.ReadAll(resp.Body)

        parse := map[string]interface{} {}
        json.Unmarshal([]byte(string(body)), &parse)
        session.Set("username", parse["login"])
        session.Set("user_id", parse["id"])
        session.Set("avatar", parse["avatar_url"])

        log.Println("retrieved data")
    } else if tokens.Expired() && session.Get("username") != nil {
        session.Clear()
        log.Println("removed data")
    }
}
