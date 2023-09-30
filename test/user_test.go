package test

import (
	"backend-auth/controllers"
	"backend-auth/database"
	"backend-auth/models"
	"backend-auth/utils"
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
	"strings"
	"testing"
	"time"
)

var controller *controllers.Controller
var validator *utils.CustomValidator
var dbConnection *gorm.DB
var redisClient *redis.Client
var bgCtx context.Context

func TestMain(m *testing.M) {
	if err := godotenv.Load("../.env"); err != nil {
		panic(fmt.Sprintf("Cannot initialize env vars for tests: %s", err.Error()))
	}
	bgCtx = context.Background()
	cleanup() // used to delete any data saved in any data source
	controller = utils.InitializeController()
	validator = utils.InitializeValidator()
	exitVal := m.Run()
	os.Exit(exitVal)
}

func cleanup() {
	dbConnection, _ = database.InitializeConnection()
	dbConnection.Unscoped().Where("1 = 1").Delete(&models.UserTokens{})
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
	req := httptest.NewRequest(http.MethodGet, "/users/count", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	if assert.NoError(t, controller.GetUsersCount(ctx)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "0", rec.Body.String())
	}
}

func TestCreateUserSuccessfully(t *testing.T) {
	userJson := readRequestFile("requests/user/successful.json")
	ctx, _, rec := sendRequest(http.MethodPost, "/users", strings.NewReader(userJson), validator)

	if assert.NoError(t, controller.CreateUser(ctx)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		output := map[string]string{}
		err := json.Unmarshal(rec.Body.Bytes(), &output)
		assert.Nil(t, err)
		expectedAccessTokenExpiration := time.Now().Add(time.Hour * 24).Unix()
		expectedRefreshTokenExpiration := time.Now().Add(time.Hour * 24 * 90).Unix()
		deviceID := parseJWT(t, output["access_token"], expectedAccessTokenExpiration)
		parseJWT(t, output["refresh_token"], expectedRefreshTokenExpiration)
		var user models.User
		var userTokens models.UserTokens
		dbConnection.First(&user)
		assert.NotNil(t, user)
		dbConnection.Where("user_id = ?", user.ID).First(&userTokens)
		assert.Equal(t, userTokens.AccessToken, output["access_token"])
		assert.Equal(t, userTokens.RefreshToken, output["refresh_token"])
		cachedAccessToken, _ := redisClient.Get(bgCtx, fmt.Sprintf("access-token::%d::%s", user.ID, deviceID)).Result()
		cachedRefreshToken, _ := redisClient.Get(bgCtx, fmt.Sprintf("refresh-token::%d::%s", user.ID, deviceID)).Result()
		assert.Equal(t, cachedAccessToken, userTokens.AccessToken)
		assert.Equal(t, cachedRefreshToken, userTokens.RefreshToken)
	}
}

func parseJWT(t *testing.T, stringToken string, expectedExpiration int64) string {
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
