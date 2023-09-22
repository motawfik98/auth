package test

import (
	"backend-auth/controllers"
	"backend-auth/utils"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var controller *controllers.Controller

func TestMain(m *testing.M) {
	if err := godotenv.Load("../.env"); err != nil {
		panic(fmt.Sprintf("Cannot initialize env vars for tests: %s", err.Error()))
	}
	controller = utils.InitializeController()
	exitVal := m.Run()
	os.Exit(exitVal)
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
