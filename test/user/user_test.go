package user

import (
	"backend-auth/controllers"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateUser(t *testing.T) {
	t.Run("should return 200 status ok", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/books", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		controller := controllers.Controller{}
		controller.CreateUser(c)

		response := map[string]string{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NotNil(t, response["access_token"])
		assert.NotNil(t, response["refresh_token"])

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}
