package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/failuretoload/datamonster/helpers"
	"github.com/failuretoload/datamonster/web"
	"github.com/golang-jwt/jwt/v5"
	"github.com/workos/workos-go/v4/pkg/usermanagement"
)

const (
	invalidJWTErrorMessage = "bad credentials"
)

func GetKeyFunc(ctx context.Context) (keyfunc.Keyfunc, error) {
	usermanagement.SetAPIKey(helpers.SafeGetEnv("WORKOS_API_KEY"))
	jwksUrl, workosErr := usermanagement.GetJWKSURL(helpers.SafeGetEnv("WORKOS_CLIENT_ID"))
	if workosErr != nil {
		return nil, fmt.Errorf("unable to retrieve JWKS: %w", workosErr)
	}
	return keyfunc.NewDefaultCtx(ctx, []string{jwksUrl.String()})
}

func ValidateJWTNew(keyfunc keyfunc.Keyfunc, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeaderParts := strings.Fields(r.Header.Get("Authorization"))
		if len(authHeaderParts) > 0 && strings.ToLower(authHeaderParts[0]) != "bearer" {
			errorMessage := web.ErrorMessage{Message: invalidJWTErrorMessage}
			if err := web.WriteJSON(w, http.StatusUnauthorized, errorMessage); err != nil {
				log.Printf("Failed to write error message: %v", err)
			}
			return
		}
		log.Default().Println(authHeaderParts[1])
		parsed, err := jwt.Parse(authHeaderParts[1], keyfunc.Keyfunc)
		if err != nil {
			web.Unauthorized(w, invalidJWTErrorMessage)
			return
		}
		if !parsed.Valid {
			web.Unauthorized(w, invalidJWTErrorMessage)
			return
		}
		userId, subjectErr := parsed.Claims.GetSubject()
		if subjectErr != nil {
			web.Unauthorized(w, invalidJWTErrorMessage)
			return
		}

		ctx := context.WithValue(r.Context(), web.UserIdKey, userId)
		log.Default().Println("authorization successful")
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
