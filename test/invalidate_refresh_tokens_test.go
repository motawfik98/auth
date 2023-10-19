package test

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestInvalidateRefreshToken(t *testing.T) {
	user := createUser(map[string]interface{}{}, false, dbConnection)
	userJson := map[string]string{
		"email":            user.Email,
		"password":         user.Password,
		"confirm_password": user.Password,
		"full_name":        user.FullName,
	}
	marshal, _ := json.Marshal(userJson)
	ctx, _, rec := sendRequest(echo.POST, "/users", strings.NewReader(string(marshal)), validator, nil)
	_ = server.CreateUser(ctx)
	assert.Equal(t, http.StatusCreated, rec.Code)
	output := map[string]string{}
	_ = json.Unmarshal(rec.Body.Bytes(), &output)
	refreshToken1 := output["refresh_token"]
	userID, deviceID, refreshExpiration := parseJWT(refreshToken1)
	headers := map[string]string{
		"x-user-id":      strconv.Itoa(int(userID)),
		"x-device-id":    deviceID,
		"x-token-expiry": strconv.FormatInt(refreshExpiration, 10),
		"Authorization":  fmt.Sprintf("Bearer %s", refreshToken1),
	}
	ctx, _, rec = sendRequest(echo.GET, "/refresh-tokens", nil, validator, headers)
	_ = server.RefreshTokens(ctx)
	assert.Equal(t, rec.Code, http.StatusOK)

	_ = json.Unmarshal(rec.Body.Bytes(), &output)
	refreshToken2 := output["refresh_token"]
	headers["Authorization"] = fmt.Sprintf("Bearer %s", refreshToken2)

	ctx, _, rec = sendRequest(echo.GET, "/refresh-tokens", nil, validator, headers)
	_ = server.RefreshTokens(ctx)
	assert.Equal(t, rec.Code, http.StatusOK)
	_ = json.Unmarshal(rec.Body.Bytes(), &output)
	refreshToken3 := output["refresh_token"]

	workerBody := map[string]string{
		"refresh_token": refreshToken1,
	}
	bytes, _ := json.Marshal(workerBody)
	err := worker.InvalidateCompromisedRefreshTokens(bytes)
	assert.NoError(t, err)

	found, _ := redisClient.SIsMember(bgCtx, "compromised-refresh-tokens", refreshToken1).Result()
	assert.True(t, found)
	found, _ = redisClient.SIsMember(bgCtx, "compromised-refresh-tokens", refreshToken2).Result()
	assert.True(t, found)
	found, _ = redisClient.SIsMember(bgCtx, "compromised-refresh-tokens", refreshToken3).Result()
	assert.True(t, found)
}
