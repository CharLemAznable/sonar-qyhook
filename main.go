package main

import (
    . "github.com/CharLemAznable/gokits"
    "github.com/kataras/golog"
    "net/http"
)

func main() {
    mux := http.NewServeMux()
    HandleFunc(mux, "/", EmptyHandler)
    HandleFunc(mux, "/qyhook", qyhook)
    HandleFunc(mux, "/badge", badge, GzipResponseDisabled)
    server := http.Server{Addr: ":" + StrFromInt(globalConfig.Port), Handler: mux}
    if err := server.ListenAndServe(); err != nil {
        golog.Errorf("Start server Error: %s", err.Error())
    }
}
