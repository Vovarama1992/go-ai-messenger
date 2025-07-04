// @title Auth Service API
// @version 1.0
// @description Сервис для логина и регистрации пользователей
// @BasePath /

// @contact.name API Support
// @contact.email support@example.com

package httpadapter

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, handler *Handler) {
	// @Summary Логин
	// @Description Аутентификация пользователя и выдача JWT
	// @Accept json
	// @Produce json
	// @Param login body loginRequest true "Email и пароль"
	// @Success 200 {object} tokenResponse
	// @Failure 401 {string} string "unauthorized"
	// @Router /login [post]
	mux.HandleFunc("/login", handler.Login)

	// @Summary Регистрация
	// @Description Регистрация нового пользователя
	// @Accept json
	// @Produce json
	// @Param register body registerRequest true "Email и пароль"
	// @Success 200 {object} tokenResponse
	// @Failure 409 {string} string "conflict"
	// @Router /register [post]
	mux.HandleFunc("/register", handler.Register)
}
