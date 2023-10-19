package test

import (
	"backend-auth/internal/utils"
	"backend-auth/internal/utils/connection"
	"backend-auth/pkg/database"
	"backend-auth/pkg/models"
	"backend-auth/pkg/servers"
	"backend-auth/pkg/workers"
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

var server *servers.Server
var worker *workers.Worker
var validator *utils.CustomValidator
var dbConnection *gorm.DB
var redisClient *redis.Client
var bgCtx context.Context

func TestMain(m *testing.M) {
	if err := godotenv.Load("../configs/dev/.env"); err != nil {
		panic(fmt.Sprintf("Cannot initialize env vars for tests: %s", err.Error()))
	}
	bgCtx = context.Background()
	cleanup() // used to delete any data saved in any data source
	server = connection.InitializeServer()
	worker = connection.InitializeWorker()
	validator = connection.InitializeValidator()
	exitVal := m.Run()
	os.Exit(exitVal)
}

func cleanup() {
	dbConnection, _ = database.InitializeConnection()
	dbConnection.Unscoped().Where("1 = 1").Delete(&models.GeneratedRefreshToken{})
	dbConnection.Unscoped().Where("1 = 1").Delete(&models.UserTokens{})
	dbConnection.Unscoped().Where("1 = 1").Delete(&models.UsedRefreshToken{})
	dbConnection.Unscoped().Where("1 = 1").Delete(&models.User{})
	opts, err := redis.ParseURL(os.ExpandEnv("redis://${REDIS_USERNAME}:${REDIS_PASSWORD}@${REDIS_HOST}:${REDIS_PORT}/#{REDIS_DB}"))
	if err != nil {
		panic(err)
	}
	redisClient = redis.NewClient(opts)
	redisClient.FlushAll(bgCtx)
}

func TestUsersCount(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/users/count", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	if assert.NoError(t, server.GetUsersCount(ctx)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.NotNil(t, "0", rec.Body.String())
	}
}

func TestCreateUserSuccessfully(t *testing.T) {
	userJson := readRequestFile("requests/user/successful.json")
	ctx, _, rec := sendRequest(echo.POST, "/users", strings.NewReader(userJson), validator, nil)

	if assert.NoError(t, server.CreateUser(ctx)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		output := map[string]string{}
		err := json.Unmarshal(rec.Body.Bytes(), &output)
		assert.Nil(t, err)
		expectedAccessTokenExpiration := time.Now().Add(time.Hour * 24).Unix()
		expectedRefreshTokenExpiration := time.Now().Add(time.Hour * 24 * 90).Unix()
		deviceID := checkJWT(t, output["access_token"], expectedAccessTokenExpiration)
		checkJWT(t, output["refresh_token"], expectedRefreshTokenExpiration)
		var user models.User
		var userTokens models.UserTokens
		dbConnection.Last(&user)
		assert.NotNil(t, user)
		dbConnection.Where("user_id = ?", user.ID).Last(&userTokens)
		assert.Equal(t, userTokens.AccessToken, output["access_token"])
		assert.Equal(t, userTokens.RefreshToken, output["refresh_token"])
		cachedAccessToken, _ := redisClient.Get(bgCtx, fmt.Sprintf("access-token::%d::%s", user.ID, deviceID)).Result()
		cachedRefreshToken, _ := redisClient.Get(bgCtx, fmt.Sprintf("refresh-token::%d::%s", user.ID, deviceID)).Result()
		assert.Equal(t, cachedAccessToken, userTokens.AccessToken)
		assert.Equal(t, cachedRefreshToken, userTokens.RefreshToken)
	}
}

func TestRefreshTokens(t *testing.T) {
	userJson := readRequestFile("requests/user/refresh.json")
	ctx, _, rec := sendRequest(echo.POST, "/users", strings.NewReader(userJson), validator, nil)
	if assert.NoError(t, server.CreateUser(ctx)) {
		output := map[string]string{}
		_ = json.Unmarshal(rec.Body.Bytes(), &output)
		refreshToken := output["refresh_token"]
		userID, deviceID, refreshExpiration := parseJWT(refreshToken)
		headers := map[string]string{
			"x-user-id":      strconv.Itoa(int(userID)),
			"x-device-id":    deviceID,
			"x-token-expiry": strconv.FormatInt(refreshExpiration, 10),
			"Authorization":  fmt.Sprintf("Bearer %s", refreshToken),
		}
		ctx, _, rec = sendRequest(echo.GET, "/refresh-tokens", nil, validator, headers)
		if assert.NoError(t, server.RefreshTokens(ctx)) {
			assert.Equal(t, rec.Code, http.StatusOK)
			output := map[string]string{}
			_ = json.Unmarshal(rec.Body.Bytes(), &output)
			refreshedAccessToken := output["access_token"]
			refreshedRefreshToken := output["refresh_token"]
			var userTokens []models.UserTokens
			dbConnection.Where("user_id = ?", userID).Find(&userTokens)
			assert.Equal(t, len(userTokens), 1)
			assert.Equal(t, userTokens[0].AccessToken, refreshedAccessToken)
			assert.Equal(t, userTokens[0].RefreshToken, refreshedRefreshToken)
			cachedAccessToken, _ := redisClient.Get(bgCtx, fmt.Sprintf("access-token::%d::%s", userID, deviceID)).Result()
			cachedRefreshToken, _ := redisClient.Get(bgCtx, fmt.Sprintf("refresh-token::%d::%s", userID, deviceID)).Result()
			assert.Equal(t, userTokens[0].AccessToken, cachedAccessToken)
			assert.Equal(t, userTokens[0].RefreshToken, cachedRefreshToken)
			var usedRefreshToken models.UsedRefreshToken
			dbConnection.Last(&usedRefreshToken)
			assert.Equal(t, usedRefreshToken.UserID, userID)
			assert.Equal(t, usedRefreshToken.RefreshToken, refreshToken)
			assert.Equal(t, usedRefreshToken.RefreshTokenExpiry.Unix(), refreshExpiration)
			found, _ := redisClient.SIsMember(bgCtx, "used-refresh-tokens", refreshToken).Result()
			assert.True(t, found)
			found, _ = redisClient.SIsMember(bgCtx, "used-refresh-tokens", refreshedRefreshToken).Result()
			assert.False(t, found)
		}
		// sending the request again with the same refresh token, should return 400
		ctx, _, rec = sendRequest(echo.GET, "/refresh-tokens", nil, validator, headers)
		server.RefreshTokens(ctx)
		assert.Equal(t, rec.Code, http.StatusBadRequest)
	}
}

func parseJWT(stringToken string) (uint, string, int64) {
	token, _, _ := new(jwt.Parser).ParseUnverified(stringToken, jwt.MapClaims{})
	claims, _ := token.Claims.(jwt.MapClaims)
	expirationTime := time.Unix(int64(claims["exp"].(float64)), 0).Unix()
	deviceID := claims["device_id"].(string)
	userID := uint(claims["id"].(float64))
	return userID, deviceID, expirationTime
}

func checkJWT(t *testing.T, stringToken string, expectedExpiration int64) string {
	token, _, err := new(jwt.Parser).ParseUnverified(stringToken, jwt.MapClaims{})
	assert.Nil(t, err)
	claims, ok := token.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	expirationTime := time.Unix(int64(claims["exp"].(float64)), 0).Unix()
	// instead of freezing the time, we subtract the two dates and make sure that the difference is less than 5 seconds
	// they should be the exact same, but we're adding this buffer just in case
	timeDifferenceInSeconds := expectedExpiration - expirationTime
	assert.Less(t, timeDifferenceInSeconds, int64(5))
	assert.NotNil(t, claims["id"])
	assert.NotNil(t, claims["device_id"])
	return claims["device_id"].(string)
}
