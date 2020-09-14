package main

import (
    . "github.com/CharLemAznable/gokits"
    "net/http"
)

func main() {
    mux := http.NewServeMux()
    HandleFunc(mux, "/", EmptyHandler)
    HandleFunc(mux, "/qyhook", qyhook)
    HandleFunc(mux, "/badge", badge, GzipResponseDisabled)
    server := http.Server{Addr: ":" + StrFromInt(appConfig.Port), Handler: mux}
    if err := server.ListenAndServe(); err != nil {
        LOG.Crashf("Start server Error: %s", err.Error())
    }
}
