package test

import (
	"backend-auth/internal/utils"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
)

func readFileContent(filename string, asString bool) ([]byte, string) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	if asString {
		return bytes, string(bytes)
	}
	return bytes, ""
}

func readRequestFile(filename string) string {
	_, content := readFileContent(filename, true)
	return content
}

func readResponseFile(filename string) map[string]string {
	bytes, _ := readFileContent(filename, false)
	output := map[string]string{}
	json.Unmarshal(bytes, &output)
	return output
}

func sendRequest(method, target string, body io.Reader, validator *utils.CustomValidator, headers map[string]string) (echo.Context, *http.Request, *httptest.ResponseRecorder) {
	e := echo.New()
	e.Validator = validator
	req := httptest.NewRequest(method, target, body)
	if method == http.MethodPost || method == http.MethodPut {
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	return ctx, req, rec
}
