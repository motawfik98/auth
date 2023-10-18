package test

import (
	"backend-auth/pkg/middleware"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestJWTAuth(t *testing.T) {
	e := echo.New()
	e.Use(middleware.JWTMiddleware())
	e.GET("/restricted", restrictedHandler)

	req, rec := sendMiddlewareRequest(e, map[string]string{})
	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	userJson := readRequestFile("requests/middleware/user-access.json")
	ctx, _, createRec := sendRequest(echo.POST, "/users", strings.NewReader(userJson), validator, nil)
	if assert.NoError(t, server.CreateUser(ctx)) {
		output := map[string]string{}
		_ = json.Unmarshal(createRec.Body.Bytes(), &output)
		accessToken := output["access_token"]
		refreshToken := output["refresh_token"]
		req, rec = sendMiddlewareRequest(e, map[string]string{"Authorization": fmt.Sprintf("Bearer %s", accessToken)})
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Empty(t, req.Header.Get("x-user-id"))
		assert.Empty(t, req.Header.Get("x-device-id"))
		assert.Empty(t, req.Header.Get("x-token-expiry"))
		_ = json.Unmarshal(rec.Body.Bytes(), &output)
		assert.NotEmpty(t, output["user_id"])
		assert.NotEmpty(t, output["device_id"])
		assert.NotEmpty(t, output["token_expiry"])
		req, rec = sendMiddlewareRequest(e, map[string]string{"Authorization": fmt.Sprintf("Bearer %s", refreshToken)})
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	}
}

func TestJWTRefresh(t *testing.T) {
	e := echo.New()
	e.Use(middleware.JWTRefreshMiddleware())
	e.GET("/restricted", restrictedHandler)

	_, rec := sendMiddlewareRequest(e, map[string]string{})
	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	userJson := readRequestFile("requests/middleware/user-refresh.json")
	ctx, _, createRec := sendRequest(echo.POST, "/users", strings.NewReader(userJson), validator, nil)
	if assert.NoError(t, server.CreateUser(ctx)) {
		output := map[string]string{}
		_ = json.Unmarshal(createRec.Body.Bytes(), &output)
		accessToken := output["access_token"]
		refreshToken := output["refresh_token"]
		_, rec = sendMiddlewareRequest(e, map[string]string{"Authorization": fmt.Sprintf("Bearer %s", accessToken)})
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		_, rec = sendMiddlewareRequest(e, map[string]string{"Authorization": fmt.Sprintf("Bearer %s", refreshToken)})
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func restrictedHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{
		"user_id":      c.Request().Header.Get("x-user-id"),
		"device_id":    c.Request().Header.Get("x-device-id"),
		"token_expiry": c.Request().Header.Get("x-token-expiry"),
	})
}

func sendMiddlewareRequest(e *echo.Echo, headers map[string]string) (*http.Request, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(echo.GET, "/restricted", nil)
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	rec := httptest.NewRecorder()
	// Using the ServerHTTP on echo will trigger the router and middleware
	e.ServeHTTP(rec, req)
	return req, rec
}
