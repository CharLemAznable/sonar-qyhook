package main

import (
    "github.com/CharLemAznable/gokits"
    "github.com/CharLemAznable/varys-go-driver"
)

func sendQyMessage(msg string) {
    qyWxMsg := map[string]interface{}{
        "touser":   "@all",
        "toparty":  "@all",
        "totag":    "@all",
        "msgtype":  "markdown",
        "agentid":  globalConfig.QyWxAgentId,
        "safe":     0,
        "markdown": map[string]string{"content": msg},
    }
    _, _ = varys.WechatCorp(globalConfig.QyWxAgentId, "/cgi-bin/message/send").
        Prop("Content-Type", "application/json;charset=utf-8").
        RequestBody(gokits.Json(qyWxMsg)).Post()
}
