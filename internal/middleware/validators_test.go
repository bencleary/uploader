package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

var validationCases = []struct {
	name          string
	key           string
	expectedError string
	expectedCode  int
}{
	{
		name:          "valid 32-byte key",
		key:           "12345678901234567890123456789012", // 32 bytes
		expectedError: "",
		expectedCode:  200,
	},
	{
		name:          "missing key",
		key:           "", // 32 bytes
		expectedError: "Encryption key is required",
		expectedCode:  401,
	},
	{
		name:          "missing key",
		key:           "",
		expectedError: "Encryption key is required",
		expectedCode:  401,
	},
	{
		name:          "invalid length",
		key:           "short",
		expectedError: "Encryption key is invalid",
		expectedCode:  401,
	},
	{
		name:          "invalid key format",
		key:           "00000000000000000000000000000000", // 32 bytes but invalid
		expectedError: "Encryption key is invalid",
		expectedCode:  401,
	},
}

func TestValidateEncryptionKey(t *testing.T) {
	e := echo.New()
	e.Use(ValidateEncryptionKey)
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	for _, test := range validationCases {
		t.Run(test.name, func(t *testing.T) {
			req.Header.Set("key", test.key)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			if rec.Code != test.expectedCode {
				t.Fatal("name", test.name, "expected code", test.expectedCode, "got", rec.Code)
			}

			if rec.Code != 200 {
				var errorResponse ErrorResponse
				if err := json.Unmarshal(rec.Body.Bytes(), &errorResponse); err != nil {
					t.Fatal("Failed to unmarshal error response", err)
				}

				if errorResponse.Message != test.expectedError {
					t.Fatal("name", test.name, "expected error", test.expectedError, "got", errorResponse.Message)
				}
			} else {
				if rec.Body.String() != "OK" {
					t.Fatal("name", test.name, "expected success response OK", "got", rec.Body.String())
				}
			}
		})
	}

}
