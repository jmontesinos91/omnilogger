package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmontesinos91/ologs/logger"
	"github.com/jmontesinos91/omnilogger/internal/repositories/middleware"
	"github.com/jmontesinos91/osecurity/services/omnibackend/enum"
	"github.com/jmontesinos91/osecurity/sts"
	"github.com/jmontesinos91/terrors"
	"github.com/sirupsen/logrus"
)

// JwtVerifyMiddleware A custom middleware to validate and parse a JWT, it will propagate the claims through the context
func JwtVerifyMiddleware(logger *logger.ContextLogger, stsClient sts.ISTSClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Log(logrus.DebugLevel, "JwtVerifyMiddleware", "start jwt validation")
			claims, permissions, err := stsClient.ValidateTokenFromRequest(r, enum.LOGS)
			if err != nil {
				logger.Error(logrus.ErrorLevel, "JwtVerifyMiddleware", "JWT parsing failure: %v", err)
				terr := terrors.Unauthorized(terrors.ErrUnauthorized, "Invalid credentials", map[string]string{})
				RenderError(r.Context(), w, terr)
				return
			}

			// Validate paths with permissions
			errV := validateAccess(r, permissions, logger)
			if errV != nil {
				RenderError(r.Context(), w, errV)
				return
			}

			// All good, propagate token using context
			r = r.WithContext(stsClient.StoreClaimsV2InContext(r.Context(), claims))

			// Continue the chain
			next.ServeHTTP(w, r)
		})
	}
}

func validateAccess(r *http.Request, permissions *[]sts.Permission, logger *logger.ContextLogger) error {
	logger.Log(logrus.DebugLevel, "validateAccess", "start validate access")
	route := chi.RouteContext(r.Context()).RoutePattern()
	method := r.Method

	// Range all permissions
	for _, permission := range *permissions {
		if middleware.ValidatePermission(permission, route, method, logger) {
			return nil //if we found one permission we let the access
		}
	}

	//if we don't found a permission match we return an unauthorized
	return terrors.Unauthorized(terrors.ErrUnauthorized, "Invalid credentials", map[string]string{})
}
