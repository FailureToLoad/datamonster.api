package web

import (
	"context"
	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"log"
	"net/http"
	"slices"
	"strings"
)

const (
	permissionDeniedErrorMessage = "Permission denied"
)

type CustomClaims struct {
	Permissions []string `json:"permissions"`
}

func (c CustomClaims) Validate(_ context.Context) error {
	return nil
}

func (c CustomClaims) HasPermissions(expectedClaims []string) bool {
	if len(expectedClaims) == 0 {
		return false
	}
	for _, scope := range expectedClaims {
		if !slices.Contains(c.Permissions, scope) {
			return false
		}
	}
	return true
}

type ctxUserIdKey string

const UserIdKey ctxUserIdKey = "userId"

func ValidatePermissions(expectedClaims []string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
		claims := token.CustomClaims.(*CustomClaims)
		if !claims.HasPermissions(expectedClaims) {
			errorMessage := ErrorMessage{Message: permissionDeniedErrorMessage}
			if err := WriteJSON(w, http.StatusForbidden, errorMessage); err != nil {
				log.Printf("Failed to write error message: %v", err)
			}
			return
		}
		subjectParts := strings.Split(token.RegisteredClaims.Subject, "|")
		ctx := context.WithValue(r.Context(), UserIdKey, subjectParts[1])
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}
}
