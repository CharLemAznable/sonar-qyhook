package main

const SonarProjectKeyHeaderName = "X-SonarQube-Project"

type MetricAppender func(string, SonarPayloadQualityGateCondition) string

var SonarMetricNameArray []string
var SonarMetricTitleMap map[string]string
var SonarMetricAppenderMap map[string]MetricAppender

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
}
