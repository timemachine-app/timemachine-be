package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/timemachine-app/timemachine-be/internal/config"
	"github.com/timemachine-app/timemachine-be/superbase"
)

type AppleTokenResponse struct {
	AccessToken string `json:"access_token"`
	IDToken     string `json:"id_token"`
}

type AccountHandler struct {
	signInWithAppleConfig config.SignInWithAppleConfig
	supabaseClient        *superbase.SupabaseClient
	jwtSecret             string
}

func NewAccountHandler(signInWithAppleConfig config.SignInWithAppleConfig, supabaseClient *superbase.SupabaseClient, jwtSecret string) *AccountHandler {
	return &AccountHandler{
		signInWithAppleConfig: signInWithAppleConfig,
		supabaseClient:        supabaseClient,
		jwtSecret:             jwtSecret,
	}
}

func (h *AccountHandler) SignInWithApple(c *gin.Context) {
	var req struct {
		Code       string `json:"code"`
		ExternalId string `json:"externalId"`
		Email      string `json:"email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	clientSecret, err := generateClientSecret(h.signInWithAppleConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate client secret"})
		return
	}

	resp, err := http.PostForm("https://appleid.apple.com/auth/token", url.Values{
		"client_id":     {h.signInWithAppleConfig.AppleClientId},
		"client_secret": {clientSecret},
		"code":          {req.Code},
		"grant_type":    {"authorization_code"},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get token from Apple"})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
		return
	}

	var tokenResp AppleTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal token response"})
		return
	}

	user, err := h.supabaseClient.GetUser(req.ExternalId)
	if err != nil {
		if err.Error() != "user not found" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
			return
		}

		user = superbase.User{
			Email:          req.Email,
			ExternalUserId: req.ExternalId,
		}

		userId, err := h.supabaseClient.AddUser(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add user"})
			return
		}
		user.UserId = userId
	}

	// Generate JWT token for the user
	jwtToken, err := GenerateJWTToken(user, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"jwt_token": jwtToken,
	})
}

func (h *AccountHandler) DeleteAccount(c *gin.Context) {
	var req struct {
		UserId string `json:"userId"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := h.supabaseClient.DeleteUser(req.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "true"})
}

func generateClientSecret(singInWithAppleConfig config.SignInWithAppleConfig) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"iss": singInWithAppleConfig.TeamId,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour).Unix(),
		"aud": "https://appleid.apple.com",
		"sub": singInWithAppleConfig.AppleClientId,
	})
	token.Header["kid"] = singInWithAppleConfig.KeyId

	privateKey, err := jwt.ParseECPrivateKeyFromPEM([]byte(singInWithAppleConfig.PrivateKey))
	if err != nil {
		return "", err
	}

	return token.SignedString(privateKey)
}

func GenerateJWTToken(user superbase.User, jwtSecret string) (string, error) {
	claims := jwt.MapClaims{
		"sub": user.UserId,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}
