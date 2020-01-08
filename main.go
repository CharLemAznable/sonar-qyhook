package main

import (
    "github.com/CharLemAznable/gokits"
    "net/http"
)

func main() {
    mux := http.NewServeMux()
    gokits.HandleFunc(mux, "/qyhook", qyhook)
    server := http.Server{Addr: ":" + gokits.StrFromInt(appConfig.Port), Handler: mux}
    if err := server.ListenAndServe(); err != nil {
        gokits.LOG.Crashf("Start server Error: %s", err.Error())
    }
}
