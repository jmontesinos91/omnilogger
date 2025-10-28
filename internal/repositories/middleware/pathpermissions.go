package middleware

import (
	"net/http"
	"strings"

	"github.com/jmontesinos91/ologs/logger"
	"github.com/jmontesinos91/osecurity/sts"
	"github.com/sirupsen/logrus"
)

type Paths string

const (
	full   Paths = "/v1/logs/{id},/v1/logs,/v1/log_messages"
	export Paths = "/v1/logs/export"
)

// ValidatePermission validates requested sources based on user permissions
func ValidatePermission(permission sts.Permission, path string, method string, logger *logger.ContextLogger) bool {
	logger.Log(logrus.DebugLevel, "ValidatePermission", "start validate permissions")

	action := strings.ReplaceAll(permission.Action, "_", "")
	action = strings.ToLower(action)

	logger.Log(logrus.DebugLevel, "ValidatePermission", "Action: "+action)

	switch action {
	case "full", "FULL":
		if (strings.Contains(string(full), path) || strings.Contains(string(export), path)) && (method == http.MethodGet || method == http.MethodOptions || method == http.MethodPost) {
			logger.Log(logrus.DebugLevel, "ValidatePermission", "Full: "+action)
			return true
		}
	default:
		logger.Log(logrus.DebugLevel, "ValidatePermission", "Default Action: "+action)
	}

	return false
}
