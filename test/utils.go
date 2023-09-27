package test

import (
	"backend-auth/utils"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
)

func readFileContent(filename string) string {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func sendRequest(method string, target string, body io.Reader, validator *utils.CustomValidator) (echo.Context, *http.Request, *httptest.ResponseRecorder) {
	e := echo.New()
	e.Validator = validator
	req := httptest.NewRequest(method, target, body)
	if method == http.MethodPost || method == http.MethodPut {
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	}
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	return ctx, req, rec
}
