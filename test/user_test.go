package test

import (
	"backend-auth/controllers"
	"backend-auth/database"
	"backend-auth/models"
	"backend-auth/utils"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
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
		//output := map[string]string{}
		//json.Unmarshal(rec.Body.Bytes(), &output)
		assert.Equal(t, "0", rec.Body.String())
	}
}

func TestCreateUserSuccessfully(t *testing.T) {
	userJson := readFileContent("requests/user/successful.json")
	ctx, _, rec := sendRequest(http.MethodPost, "/users", strings.NewReader(userJson), validator)

	if assert.NoError(t, controller.CreateUser(ctx)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}
