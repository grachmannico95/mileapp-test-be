package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/grachmannico95/mileapp-test-be/internal/config"
	"github.com/grachmannico95/mileapp-test-be/internal/dto"
	"github.com/grachmannico95/mileapp-test-be/internal/service"
	"github.com/grachmannico95/mileapp-test-be/internal/util"
)

type AuthHandler struct {
	authService service.AuthService
	config      *config.Config
}

func NewAuthHandler(authService service.AuthService, config *config.Config) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		config:      config,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("validation failed", util.ParseValidationError(err)...))
		return
	}

	user, jwtToken, csrfToken, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse(err.Error()))
		return
	}

	response := dto.LoginResponse{}
	if h.config.Server.AuthCookie {
		response = dto.LoginResponse{
			User: dto.ToUserResponse(user),
		}
		h.setAuthCookies(c, jwtToken, csrfToken)
	} else {
		response = dto.LoginResponse{
			User:        dto.ToUserResponse(user),
			AccessToken: jwtToken,
			CSRFToken:   csrfToken,
		}
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("login successful", response))
}

func (h *AuthHandler) Logout(c *gin.Context) {
	h.clearAuthCookies(c)
	c.JSON(http.StatusOK, dto.SuccessResponse("logout successful", nil))
}

func (h *AuthHandler) setAuthCookies(c *gin.Context, jwtToken, csrfToken string) {
	maxAge := time.Now().Add(h.config.JWT.Expiry)
	sameSite := map[string]http.SameSite{
		"SameSite": http.SameSiteDefaultMode,
		"Lax":      http.SameSiteLaxMode,
		"Strict":   http.SameSiteStrictMode,
		"None":     http.SameSiteNoneMode,
	}

	accessTokenCookie := &http.Cookie{
		Name:     "access_token",
		Value:    jwtToken,
		Path:     "/",
		Domain:   h.config.Cookie.Domain,
		Expires:  maxAge,
		HttpOnly: h.config.Cookie.HTTPOnly,
		Secure:   h.config.Cookie.Secure,
		SameSite: sameSite[h.config.Cookie.SameSite],
	}

	csrfTokenCookie := &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Path:     "/",
		Domain:   h.config.Cookie.Domain,
		Expires:  maxAge,
		HttpOnly: false,
		Secure:   h.config.Cookie.Secure,
		SameSite: sameSite[h.config.Cookie.SameSite],
	}

	http.SetCookie(c.Writer, accessTokenCookie)
	http.SetCookie(c.Writer, csrfTokenCookie)
}

func (h *AuthHandler) clearAuthCookies(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", h.config.Cookie.Domain, h.config.Cookie.Secure, h.config.Cookie.HTTPOnly)
	c.SetCookie("csrf_token", "", -1, "/", h.config.Cookie.Domain, h.config.Cookie.Secure, false)
}
