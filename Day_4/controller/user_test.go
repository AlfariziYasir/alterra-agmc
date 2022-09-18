package controller

import (
	"api-mvc/config"
	"api-mvc/db/postgres"
	"api-mvc/db/redis"
	"api-mvc/model"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func InitEcho() *echo.Echo {
	config.Cfg()
	e := echo.New()

	return e
}

func TestGetUsers(t *testing.T) {
	var testCases = []struct {
		name                 string
		path                 string
		expectStatus         int
		expextBodyStartsWith string
	}{
		{
			name:                 "success",
			path:                 "/users/",
			expectStatus:         http.StatusOK,
			expextBodyStartsWith: "{\"status\":200,\"message\":\"get users is success\",",
		},
	}

	e := InitEcho()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	db, _ := postgres.NewClient()
	rds, _ := redis.NewClient()

	control := NewController(rds.Conn(), db.Conn())

	for _, testCase := range testCases {
		c.SetPath(testCase.path)

		if assert.NoError(t, control.GetUsers(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			body := rec.Body.String()
			assert.True(t, strings.HasPrefix(body, testCase.expextBodyStartsWith))
		}
	}
}

func TestCreateUser(t *testing.T) {
	var testCases = []struct {
		name                 string
		path                 string
		expectStatus         int
		expextBodyStartsWith string
		payload              model.UserRequest
	}{
		{
			name:                 "success",
			path:                 "/register/",
			expectStatus:         http.StatusCreated,
			expextBodyStartsWith: "{\"status\":201,\"message\":\"create user is success\",",
			payload: model.UserRequest{
				Name:     "testing",
				Email:    "example@email.com",
				Password: "123456",
			},
		},
		{
			name:                 "conflict",
			path:                 "/register/",
			expectStatus:         http.StatusConflict,
			expextBodyStartsWith: "{\"status\":409,",
			payload: model.UserRequest{
				Name:     "tes",
				Email:    "example@email.com",
				Password: "123456",
			},
		},
	}

	db, _ := postgres.NewClient()
	rds, _ := redis.NewClient()

	control := NewController(rds.Conn(), db.Conn())

	for _, testCase := range testCases {
		payload, _ := json.Marshal(testCase.payload)

		e := InitEcho()
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		c.SetPath(testCase.path)

		if assert.NoError(t, control.CreateUser(c)) {
			assert.Equal(t, testCase.expectStatus, rec.Code)
			body := rec.Body.String()
			assert.True(t, strings.HasPrefix(body, testCase.expextBodyStartsWith))
		}
	}
}
