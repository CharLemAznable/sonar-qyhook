package main

import (
    "compress/gzip"
    "github.com/CharLemAznable/gokits"
    "io"
    "net/http"
    "net/http/httputil"
    "strings"
)

func dumpRequest(handlerFunc http.HandlerFunc) http.HandlerFunc {
    return func(writer http.ResponseWriter, request *http.Request) {
        // Save a copy of this request for debugging.
        requestDump, err := httputil.DumpRequest(request, true)
        if err != nil {
            _ = gokits.LOG.Error(err)
        }
        gokits.LOG.Debug(string(requestDump))
        handlerFunc(writer, request)
    }
}

type GzipResponseWriter struct {
    io.Writer
    http.ResponseWriter
}

func (w GzipResponseWriter) Write(b []byte) (int, error) {
    return w.Writer.Write(b)
}

func gzipHandlerFunc(handlerFunc http.HandlerFunc) http.HandlerFunc {
    return func(writer http.ResponseWriter, request *http.Request) {
        if !strings.Contains(request.Header.Get("Accept-Encoding"), "gzip") {
            handlerFunc(writer, request)
            return
        }
        writer.Header().Set("Content-Encoding", "gzip")
        gz := gzip.NewWriter(writer)
        defer func() { _ = gz.Close() }()
        gzr := GzipResponseWriter{Writer: gz, ResponseWriter: writer}
        handlerFunc(gzr, request)
    }
}
