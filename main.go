package main

import (
    "github.com/CharLemAznable/gokits"
    "net/http"
)

func main() {
    mux := http.NewServeMux()
    handleFunc(mux, "/qyhook", qyhook, true)
    server := http.Server{Addr: ":" + gokits.StrFromInt(appConfig.Port), Handler: mux}
    if err := server.ListenAndServe(); err != nil {
        gokits.LOG.Crashf("Start server Error: %s", err.Error())
    }
}

func handleFunc(mux *http.ServeMux, path string, handlerFunc http.HandlerFunc, requiredDump bool) {
    wrap := handlerFunc
    if requiredDump {
        wrap = dumpRequest(handlerFunc)
    }

    wrap = gzipHandlerFunc(wrap)
    handlePath := gokits.PathJoin(appConfig.ContextPath, path)
    mux.HandleFunc(handlePath, wrap)
}
