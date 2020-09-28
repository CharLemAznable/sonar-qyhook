package main

import (
    "errors"
    "github.com/CharLemAznable/gokits"
    "github.com/tidwall/gjson"
    "go.etcd.io/bbolt"
    "net/http"
    "net/url"
)

type BadgeInfo struct {
    ProjectKey  string
    BranchName  string
    Metric      string
    MetricLabel string
    MetricValue string
    MetricColor string
}

func badge(writer http.ResponseWriter, request *http.Request) {
    badgeInfo, err := readRequestBadgeInfo(request)
    if nil != err {
        if gokits.IsAjaxRequest(request) {
            gokits.ResponseJson(writer,
                gokits.Json(map[string]string{"msg": err.Error()}))
        } else {
            http.Error(writer, "Request Parameters Error", http.StatusBadRequest)
        }
        return
    }

    label := url.QueryEscape(badgeInfo.MetricLabel)
    message := url.QueryEscape(badgeInfo.MetricValue)
    color := url.QueryEscape(badgeInfo.MetricColor)

    req := request
    req.URL.RawQuery = "label=" + label + "&message=" + message + "&color=" + color + "&logo=sonarqube"
    req.URL.Path = ""
    shieldsProxy.ServeHTTP(writer, req)
}

func readRequestBadgeInfo(request *http.Request) (*BadgeInfo, error) {
    projectKey := request.FormValue("projectKey")
    if "" == projectKey {
        return nil, errors.New("缺少参数projectKey")
    }
    branchName := request.FormValue("branchName")
    if "" == branchName {
        branchName = "master"
    }
    metric := request.FormValue("metric")
    if "" == metric {
        metric = "quality_gate"
    }
    badgeInfo := new(BadgeInfo)
    badgeInfo.ProjectKey = projectKey
    badgeInfo.BranchName = branchName
    badgeInfo.Metric = metric

    label, ok := SonarMetricLabelMap[badgeInfo.Metric]
    if !ok {
        badgeInfo.MetricLabel = "unknown"
        badgeInfo.MetricValue = "unknown"
        badgeInfo.MetricColor = "lightgray"
        return badgeInfo, nil
    }
    badgeInfo.MetricLabel = label

    err := db.View(func(tx *bbolt.Tx) error {
        bucket := tx.Bucket([]byte(PayloadBucket))
        payloadValue := string(bucket.Get([]byte(projectKey + ":" + branchName)))
        if "" == payloadValue {
            badgeInfo.MetricValue = "unknown"
            badgeInfo.MetricColor = "lightgray"
            return nil
        }
        payloadInfo, ok := gokits.UnJson(payloadValue,
            new(SonarPayload)).(*SonarPayload)
        if !ok || nil == payloadInfo {
            badgeInfo.MetricValue = "unexpected"
            badgeInfo.MetricColor = "red"
            return nil
        }
        statusPath, ok := SonarMetricStatusJsonPathMap[badgeInfo.Metric]
        if !ok {
            badgeInfo.MetricValue = "unknown"
            badgeInfo.MetricColor = "lightgray"
            return nil
        }
        status := gjson.Get(payloadValue, statusPath)
        if status.Exists() {
            badgeInfo.MetricValue = status.String()
            badgeInfo.MetricColor = SonarMetricStatusMapper(status.String())
        } else {
            badgeInfo.MetricValue = "unknown"
            badgeInfo.MetricColor = "lightgray"
            return nil
        }
        path, ok := SonarMetricValueJsonPathMap[badgeInfo.Metric]
        if !ok {
            return nil // no further value should be parsed
        }
        parser := SonarMetricValueParserMap[badgeInfo.Metric]
        mapper, ok := SonarMetricValueColorMapperMap[badgeInfo.Metric]
        value := gjson.Get(payloadValue, path)
        if value.Exists() {
            badgeInfo.MetricValue = parser(status.String(), value.String())
            if ok {
                badgeInfo.MetricColor = mapper(badgeInfo.MetricValue)
            }
        } else {
            badgeInfo.MetricValue = "unknown"
            badgeInfo.MetricColor = "lightgray"
        }
        return nil
    })
    if err != nil {
        badgeInfo.MetricValue = "unexpected"
        badgeInfo.MetricColor = "red"
        return badgeInfo, nil
    }

    return badgeInfo, nil
}
