package main

import (
    "flag"
    "github.com/BurntSushi/toml"
    "github.com/CharLemAznable/gokits"
    "regexp"
    "strings"
)

type AppConfig struct {
    Port              int
    ContextPath       string
    VarysBaseUrl      string
    QyWxAgentId       string
    QyWxMessageUrl    string
    ProjectKeyPattern string
}

var appConfig AppConfig
var _configFile string
var projectKeyRegexp *regexp.Regexp

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
    gokits.If(0 == len(appConfig.VarysBaseUrl), func() {
        appConfig.VarysBaseUrl = "http://127.0.0.1:4236/varys/query-wechat-corp-token/"
    })
    gokits.If(0 == len(appConfig.QyWxAgentId), func() {
        gokits.LOG.Crashf("QyWxAgentId config REQUIRED")
    })
    gokits.If(0 == len(appConfig.QyWxMessageUrl), func() {
        appConfig.QyWxMessageUrl = "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token="
    })
    gokits.If(0 == len(appConfig.ProjectKeyPattern), func() {
        appConfig.ProjectKeyPattern = "^.*$"
    })
    projectKeyRegexp = regexp.MustCompile(appConfig.ProjectKeyPattern)

    gokits.LOG.Debug("appConfig: %s", gokits.Json(appConfig))
}
