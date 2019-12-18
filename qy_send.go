package main

import (
    "github.com/CharLemAznable/gokits"
)

type QyWxAccessToken struct {
    CorpId string `json:"corpId"`
    Token  string `json:"token"`
}

func sendQyMessage(msg string) {
    rsp, err := gokits.NewHttpReq(appConfig.VarysBaseUrl + appConfig.QyWxAgentId).Get()
    if nil != err {
        _ = gokits.LOG.Warn("sendQyMessage: %s [ %s ]", err.Error(), rsp)
        return
    }

    token := gokits.UnJson(rsp, new(QyWxAccessToken)).(*QyWxAccessToken)
    qyWxMsg := map[string]interface{}{
        "touser":   "@all",
        "toparty":  "@all",
        "totag":    "@all",
        "msgtype":  "markdown",
        "agentid":  appConfig.QyWxAgentId,
        "safe":     0,
        "markdown": map[string]string{"content": msg},
    }
    _, _ = gokits.NewHttpReq(appConfig.QyWxMessageUrl + token.Token).
        Prop("Content-Type", "application/json;charset=utf-8").
        RequestBody(gokits.Json(qyWxMsg)).Post()
}
