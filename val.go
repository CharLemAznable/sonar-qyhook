package main

import (
    "github.com/CharLemAznable/gokits"
)

const SonarProjectKeyHeaderName = "X-SonarQube-Project"

type MetricAppender func(string, SonarPayloadQualityGateCondition) string
type MetricValueParser func(string) string
type MetricColorMapper func(string) string

var SonarMetricNameArray []string
var SonarMetricTitleMap map[string]string
var SonarMetricAppenderMap map[string]MetricAppender
var SonarMetricLabelMap map[string]string
var SonarMetricJsonPathMap map[string]string
var SonarMetricValueParserMap map[string]MetricValueParser
var SonarMetricColorMapperMap map[string]MetricColorMapper

func init() {
    SonarMetricNameArray = []string{
        "reliability_rating",
        "new_reliability_rating",
        "security_rating",
        "new_security_rating",
        "sqale_rating",
        "new_maintainability_rating",
        "coverage",
        "new_coverage",
        "duplicated_lines_density",
        "new_duplicated_lines_density",
    }
    SonarMetricTitleMap = map[string]string{
        "reliability_rating":           "可靠性比率：　　",
        "new_reliability_rating":       "新代码可靠率：　",
        "security_rating":              "安全比率：　　　",
        "new_security_rating":          "新代码安全率：　",
        "sqale_rating":                 "可维护性：　　　",
        "new_maintainability_rating":   "新代码可维护性：",
        "coverage":                     "覆盖率：　　　　",
        "new_coverage":                 "新代码覆盖率：　",
        "duplicated_lines_density":     "重复行：　　　　",
        "new_duplicated_lines_density": "新代码重复行：　",
    }
    SonarMetricAppenderMap = map[string]MetricAppender{
        "reliability_rating":           gradeAppender,
        "new_reliability_rating":       gradeAppender,
        "security_rating":              gradeAppender,
        "new_security_rating":          gradeAppender,
        "sqale_rating":                 gradeAppender,
        "new_maintainability_rating":   gradeAppender,
        "coverage":                     ratingValueAppender,
        "new_coverage":                 ratingValueAppender,
        "duplicated_lines_density":     ratingValueAppender,
        "new_duplicated_lines_density": ratingValueAppender,
    }
    SonarMetricLabelMap = map[string]string{
        "analysis_status":              "代码分析",
        "quality_gate":                 "质量阈",
        "reliability_rating":           "可靠性比率",
        "new_reliability_rating":       "新代码可靠率",
        "security_rating":              "安全比率",
        "new_security_rating":          "新代码安全率",
        "sqale_rating":                 "可维护性",
        "new_maintainability_rating":   "新代码可维护性",
        "coverage":                     "覆盖率",
        "new_coverage":                 "新代码覆盖率",
        "duplicated_lines_density":     "重复行",
        "new_duplicated_lines_density": "新代码重复行",
    }
    SonarMetricJsonPathMap = map[string]string{
        "analysis_status":              "status",
        "quality_gate":                 "qualityGate.status",
        "reliability_rating":           "qualityGate.conditions.#(metric==\"reliability_rating\").value",
        "new_reliability_rating":       "qualityGate.conditions.#(metric==\"new_reliability_rating\").value",
        "security_rating":              "qualityGate.conditions.#(metric==\"security_rating\").value",
        "new_security_rating":          "qualityGate.conditions.#(metric==\"new_security_rating\").value",
        "sqale_rating":                 "qualityGate.conditions.#(metric==\"sqale_rating\").value",
        "new_maintainability_rating":   "qualityGate.conditions.#(metric==\"new_maintainability_rating\").value",
        "coverage":                     "qualityGate.conditions.#(metric==\"coverage\").value",
        "new_coverage":                 "qualityGate.conditions.#(metric==\"new_coverage\").value",
        "duplicated_lines_density":     "qualityGate.conditions.#(metric==\"duplicated_lines_density\").value",
        "new_duplicated_lines_density": "qualityGate.conditions.#(metric==\"new_duplicated_lines_density\").value",
    }
    SonarMetricValueParserMap = map[string]MetricValueParser{
        "analysis_status":              func(s string) string { return s },
        "quality_gate":                 func(s string) string { return s },
        "reliability_rating":           parseRatingValueToGrade,
        "new_reliability_rating":       parseRatingValueToGrade,
        "security_rating":              parseRatingValueToGrade,
        "new_security_rating":          parseRatingValueToGrade,
        "sqale_rating":                 parseRatingValueToGrade,
        "new_maintainability_rating":   parseRatingValueToGrade,
        "coverage":                     roundRatingValue,
        "new_coverage":                 roundRatingValue,
        "duplicated_lines_density":     roundRatingValue,
        "new_duplicated_lines_density": roundRatingValue,
    }
    statusMapper := func(s string) string {
        return gokits.Condition("SUCCESS" == s || "OK" == s, "brightgreen", "lightgray").(string)
    }
    gradeMapper := func(s string) string {
        if "A" == s {
            return "brightgreen"
        } else if "B" == s {
            return "green"
        } else if "C" == s {
            return "yellow"
        } else if "D" == s {
            return "orange"
        } else if "E" == s {
            return "red"
        } else {
            return "lightgray"
        }
    }
    ratingMapper := func(s string) string { return "lightgray" }
    SonarMetricColorMapperMap = map[string]MetricColorMapper{
        "analysis_status":              statusMapper,
        "quality_gate":                 statusMapper,
        "reliability_rating":           gradeMapper,
        "new_reliability_rating":       gradeMapper,
        "security_rating":              gradeMapper,
        "new_security_rating":          gradeMapper,
        "sqale_rating":                 gradeMapper,
        "new_maintainability_rating":   gradeMapper,
        "coverage":                     ratingMapper,
        "new_coverage":                 ratingMapper,
        "duplicated_lines_density":     ratingMapper,
        "new_duplicated_lines_density": ratingMapper,
    }
}
