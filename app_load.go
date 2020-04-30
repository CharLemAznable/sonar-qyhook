package main

import (
    "flag"
    "github.com/BurntSushi/toml"
    "github.com/CharLemAznable/gokits"
    "github.com/CharLemAznable/varys-go-driver"
    "net/http/httputil"
    "net/url"
    "regexp"
    "strings"
    "unsafe"
)

type AppConfig struct {
    gokits.HttpServerConfig
    VarysBaseUrl      string
    QyWxAgentId       string
    ProjectKeyPattern string
    ShieldsBadgeUrl   string
}

var appConfig AppConfig
var _configFile string
var projectKeyRegexp *regexp.Regexp
var shieldsProxy *httputil.ReverseProxy

func init() {
    gokits.LOG.LoadConfiguration("logback.xml")

    flag.StringVar(&_configFile, "configFile", "appConfig.toml", "config file path")
    flag.Parse()

    if _, err := toml.DecodeFile(_configFile, &appConfig); err != nil {
        gokits.LOG.Crashf("config file decode error: %s", err.Error())
    }

    gokits.If(0 == appConfig.Port, func() {
        appConfig.Port = 17258
    })
    gokits.If(0 != len(appConfig.ContextPath), func() {
        gokits.Unless(strings.HasPrefix(appConfig.ContextPath, "/"),
            func() { appConfig.ContextPath = "/" + appConfig.ContextPath })
        gokits.If(strings.HasSuffix(appConfig.ContextPath, "/"),
            func() { appConfig.ContextPath = appConfig.ContextPath[:len(appConfig.ContextPath)-1] })
    })
    gokits.Unless("" == appConfig.VarysBaseUrl, func() {
        varys.ConfigInstance.Address = appConfig.VarysBaseUrl
    })
    gokits.If("" == appConfig.QyWxAgentId, func() {
        gokits.LOG.Crashf("QyWxAgentId config REQUIRED")
    })
    gokits.If("" == appConfig.ProjectKeyPattern, func() {
        appConfig.ProjectKeyPattern = "^.*$"
    })
    projectKeyRegexp = regexp.MustCompile(appConfig.ProjectKeyPattern)
    gokits.If("" == appConfig.ShieldsBadgeUrl, func() {
        appConfig.ShieldsBadgeUrl = "https://img.shields.io/static/v1"
    })
    badgeUrl, err := url.Parse(appConfig.ShieldsBadgeUrl)
    if err != nil {
        badgeUrl, _ = url.Parse("https://img.shields.io/static/v1")
    }
    shieldsProxy = gokits.ReverseProxy(badgeUrl)

    gokits.GlobalHttpServerConfig = (*gokits.HttpServerConfig)(unsafe.Pointer(&appConfig))
    gokits.LOG.Debug("appConfig: %s", gokits.Json(appConfig))
}
