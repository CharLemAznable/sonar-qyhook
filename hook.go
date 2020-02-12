package main

import (
    "fmt"
    "github.com/CharLemAznable/gokits"
    "net/http"
    "strconv"
)

type SonarPayload struct {
    ServerUrl   string                  `json:"serverUrl"`
    TaskId      string                  `json:"taskId"`
    Status      string                  `json:"status"`
    AnalysedAt  string                  `json:"analysedAt"`
    ChangedAt   string                  `json:"changedAt"`
    Project     SonarPayloadProject     `json:"project"`
    Branch      SonarPayloadBranch      `json:"branch"`
    QualityGate SonarPayloadQualityGate `json:"qualityGate"`
}

type SonarPayloadProject struct {
    Key  string `json:"key"`
    Name string `json:"name"`
    Url  string `json:"url"`
}

type SonarPayloadBranch struct {
    Name   string `json:"name"`
    Type   string `json:"type"`
    IsMain bool   `json:"isMain"`
    Url    string `json:"url"`
}

type SonarPayloadQualityGate struct {
    Name       string                             `json:"name"`
    Status     string                             `json:"status"`
    Conditions []SonarPayloadQualityGateCondition `json:"conditions"`
}

type SonarPayloadQualityGateCondition struct {
    Metric           string `json:"metric"`
    Operator         string `json:"operator"`
    Value            string `json:"value"`
    Status           string `json:"status"`
    OnLeakPeriod     bool   `json:"onLeakPeriod"`
    ErrorThreshold   string `json:"errorThreshold"`
    WarningThreshold string `json:"warningThreshold"`
}

func qyhook(writer http.ResponseWriter, request *http.Request) {
    projectKey := request.Header.Get(SonarProjectKeyHeaderName)
    if "" == projectKey {
        gokits.ResponseText(writer, "Request Illegal")
        return
    }
    if !projectKeyRegexp.MatchString(projectKey) {
        gokits.ResponseText(writer, "Ignored")
        return
    }
    body, _ := gokits.RequestBody(request)
    payload, ok := gokits.UnJson(body,
        new(SonarPayload)).(*SonarPayload)
    if !ok || nil == payload {
        gokits.ResponseText(writer, "Request Illegal")
        return
    }

    go func() {
        msg := buildMsg(projectKey, payload)
        sendQyMessage(msg)
    }()

    gokits.ResponseText(writer, "OK")
}

func buildMsg(projectKey string, payload *SonarPayload) string {
    msg := "项目`" + projectKey + "` `" +
        payload.Branch.Name + "`分析: "
    if "SUCCESS" == payload.Status {
        msg += "完成"
    } else {
        msg += "未完成"
    }
    msg += "\n"

    msg += "> 质量阈: "
    if "OK" == payload.QualityGate.Status {
        msg += "<font color=\"info\">正常</font>"
    } else {
        msg += "<font color=\"warning\">错误</font>"
    }
    msg += "\n>\n"

    qualityMap := map[string]string{}
    for _, condition := range payload.QualityGate.Conditions {
        appender, ok := SonarMetricAppenderMap[condition.Metric]
        if !ok {
            continue
        }
        title, ok := SonarMetricTitleMap[condition.Metric]
        if !ok {
            continue
        }
        qualityMap[condition.Metric] = appender(title, condition)
    }
    for _, metric := range SonarMetricNameArray {
        quality, ok := qualityMap[metric]
        if !ok {
            continue
        }
        msg += quality
    }

    msg += "\n查看详情, 请点击[链接](" + payload.Branch.Url + ")"
    return msg
}

func ratingValueAppender(title string, condition SonarPayloadQualityGateCondition) string {
    str := "> " + title
    if "NO_VALUE" == condition.Status {
        str += "<font color=\"comment\">-</font>"
    } else if "OK" == condition.Status {
        str += "<font color=\"info\">" + roundRatingValue(condition.Value) + "%</font>"
    } else {
        str += "<font color=\"warning\">" + roundRatingValue(condition.Value) + "%</font>"
    }
    str += "\n"
    return str
}

func roundRatingValue(value string) string {
    floatValue, _ := strconv.ParseFloat(value, 64)
    return fmt.Sprintf("%.1f", floatValue)
}

func gradeAppender(title string, condition SonarPayloadQualityGateCondition) string {
    str := "> " + title
    if "NO_VALUE" == condition.Status {
        str += "<font color=\"comment\">-</font>"
    } else if "OK" == condition.Status {
        str += "<font color=\"info\">" + parseRatingValueToGrade(condition.Value) + "</font>"
    } else {
        str += "<font color=\"warning\">" + parseRatingValueToGrade(condition.Value) + "</font>"
    }
    str += "\n"
    return str
}

func parseRatingValueToGrade(value string) string {
    if "1" == value {
        return "A"
    } else if "2" == value {
        return "B"
    } else if "3" == value {
        return "C"
    } else if "4" == value {
        return "D"
    } else if "5" == value {
        return "E"
    } else {
        return "?"
    }
}
