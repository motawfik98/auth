package test

import (
	"backend-auth/controllers"
	"backend-auth/database"
	"backend-auth/models"
	"backend-auth/utils"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

var controller *controllers.Controller
var validator *utils.CustomValidator

func TestMain(m *testing.M) {
	if err := godotenv.Load("../.env"); err != nil {
		panic(fmt.Sprintf("Cannot initialize env vars for tests: %s", err.Error()))
	}
	controller = utils.InitializeController()
	validator = utils.InitializeValidator()
	cleanup() // used to delete any data saved in any data source
	exitVal := m.Run()
	os.Exit(exitVal)
}

func cleanup() {
	dbConnection, _ := database.InitializeConnection()
	dbConnection.Unscoped().Where("1 = 1").Delete(&models.User{})
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
		parseJWT(t, output["access_token"], expectedAccessTokenExpiration)
		parseJWT(t, output["refresh_token"], expectedRefreshTokenExpiration)
	}
}

func parseJWT(t *testing.T, stringToken string, expectedExpiration int64) {
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
}
