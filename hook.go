package main

import (
	"github.com/CharLemAznable/gokits"
	"io/ioutil"
	"net/http"
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
	if 0 == len(projectKey) {
		gokits.ResponseText(writer, "Request Illegal")
		return
	}
	if !projectKeyRegexp.MatchString(projectKey) {
		gokits.ResponseText(writer, "Ignored")
		return
	}
	bytes, _ := ioutil.ReadAll(request.Body)
	payload, ok := gokits.UnJson(string(bytes),
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

	for _, condition := range payload.QualityGate.Conditions {
		if "new_coverage" == condition.Metric {
			msg += "> 新代码覆盖率: "
			if "NO_VALUE" == condition.Status {
				msg += "<font color=\"comment\">-</font>"
			} else if "OK" == condition.Status {
				msg += "<font color=\"info\">" + condition.Value + "%</font>"
			} else {
				msg += "<font color=\"warning\">" + condition.Value + "%</font>"
			}
			msg += "\n"

		} else if "new_duplicated_lines_density" == condition.Metric {
			msg += "> 新代码重复率: "
			if "NO_VALUE" == condition.Status {
				msg += "<font color=\"comment\">-</font>"
			} else if "OK" == condition.Status {
				msg += "<font color=\"info\">" + condition.Value + "%</font>"
			} else {
				msg += "<font color=\"warning\">" + condition.Value + "%</font>"
			}
			msg += "\n"

		} else if "new_maintainability_rating" == condition.Metric {
			msg += "> 新代码可维护率: "
			if "NO_VALUE" == condition.Status {
				msg += "<font color=\"comment\">-</font>"
			} else if "OK" == condition.Status {
				msg += "<font color=\"info\">" + parseRatingValueToGrade(condition.Value) + "</font>"
			} else {
				msg += "<font color=\"warning\">" + parseRatingValueToGrade(condition.Value) + "</font>"
			}
			msg += "\n"

		} else if "new_reliability_rating" == condition.Metric {
			msg += "> 新代码可靠率: "
			if "NO_VALUE" == condition.Status {
				msg += "<font color=\"comment\">-</font>"
			} else if "OK" == condition.Status {
				msg += "<font color=\"info\">" + parseRatingValueToGrade(condition.Value) + "</font>"
			} else {
				msg += "<font color=\"warning\">" + parseRatingValueToGrade(condition.Value) + "</font>"
			}
			msg += "\n"

		} else if "new_security_rating" == condition.Metric {
			msg += "> 新代码安全率: "
			if "NO_VALUE" == condition.Status {
				msg += "<font color=\"comment\">-</font>"
			} else if "OK" == condition.Status {
				msg += "<font color=\"info\">" + parseRatingValueToGrade(condition.Value) + "</font>"
			} else {
				msg += "<font color=\"warning\">" + parseRatingValueToGrade(condition.Value) + "</font>"
			}
			msg += "\n"
		}
	}

	msg += "\n查看详情, 请点击[链接](" + payload.Branch.Url + ")"
	return msg
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
	} else {
		return "?"
	}
}
