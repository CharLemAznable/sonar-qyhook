package main

import (
    "flag"
    "github.com/BurntSushi/toml"
    "github.com/CharLemAznable/gokits"
    "github.com/CharLemAznable/varys-go-driver"
    "github.com/kataras/golog"
    "net/http/httputil"
    "net/url"
    "regexp"
    "strings"
)

type Config struct {
    gokits.HttpServerConfig

    LogLevel string

    VarysBaseUrl string
    QyWxAgentId  string

    ProjectKeyPattern string

    ShieldsBadgeUrl string
}

var globalConfig = &Config{}
var projectKeyRegexp *regexp.Regexp
var shieldsProxy *httputil.ReverseProxy

const DefaultShieldsProxyURL = "https://img.shields.io/static/v1"

func init() {
    configFile := ""
    flag.StringVar(&configFile, "configFile",
        "config.toml", "config file path")
    flag.Parse()
    if _, err := toml.DecodeFile(configFile, globalConfig); err != nil {
        golog.Errorf("config file decode error: %s", err.Error())
    }

    fixedConfig(globalConfig)
    projectKeyRegexp = regexp.MustCompile(globalConfig.ProjectKeyPattern)
    shieldsProxy = buildShieldsProxy(globalConfig)
}

func fixedConfig(config *Config) {
    gokits.If(0 == config.Port, func() {
        config.Port = 17258
    })
    gokits.If("" != config.ContextPath, func() {
        gokits.Unless(strings.HasPrefix(config.ContextPath, "/"),
            func() { config.ContextPath = "/" + config.ContextPath })
        gokits.If(strings.HasSuffix(config.ContextPath, "/"),
            func() { config.ContextPath = config.ContextPath[:len(config.ContextPath)-1] })
    })
    gokits.If("" == config.LogLevel, func() {
        config.LogLevel = "info"
    })

    gokits.Unless("" == config.VarysBaseUrl, func() {
        varys.ConfigInstance.Address = config.VarysBaseUrl
    })
    gokits.If("" == config.QyWxAgentId, func() {
        golog.Error("QyWxAgentId config REQUIRED")
        panic("QyWxAgentId config REQUIRED")
    })

    gokits.If("" == config.ProjectKeyPattern, func() {
        config.ProjectKeyPattern = "^.*$"
    })

    gokits.If("" == config.ShieldsBadgeUrl, func() {
        config.ShieldsBadgeUrl = DefaultShieldsProxyURL
    })

    gokits.GlobalHttpServerConfig = &config.HttpServerConfig

    golog.SetLevel(config.LogLevel)
    golog.Infof("config: %+v", *config)
}

func buildShieldsProxy(config *Config) *httputil.ReverseProxy {
    badgeUrl, err := url.Parse(config.ShieldsBadgeUrl)
    if err != nil {
        badgeUrl, _ = url.Parse(DefaultShieldsProxyURL)
    }
    return gokits.ReverseProxy(badgeUrl)
}
