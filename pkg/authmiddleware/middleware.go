package authmiddleware

import (
	"context"
	"net/http"
	"strings"

	authpb "github.com/Vovarama1992/go-ai-messenger/proto/authpb"
)

type contextKey string

const userIDKey = contextKey("userID")

type AuthMiddleware struct {
	authClient authpb.AuthServiceClient
}

func NewAuthMiddleware(client authpb.AuthServiceClient) *AuthMiddleware {
	return &AuthMiddleware{authClient: client}
}

func (a *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		resp, err := a.authClient.ValidateToken(r.Context(), &authpb.ValidateTokenRequest{Token: token})
		if err != nil || resp.UserId == 0 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, resp.UserId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Извлечение userID из контекста в handler
func GetUserID(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(userIDKey).(int64)
	return userID, ok
}
