package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"golang-starter-pack/router/middleware"
	"golang-starter-pack/utils"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestSignUpCaseSuccess(t *testing.T) {
	tearDown()
	setup()
	var (
		reqJSON = `{"player":{"username":"alice","email":"alice@realworld.io","password":"secret"}}`
	)
	req := httptest.NewRequest(echo.POST, "/api/players", strings.NewReader(reqJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	assert.NoError(t, h.SignUp(c))
	if assert.Equal(t, http.StatusCreated, rec.Code) {
		m := responseMap(rec.Body.Bytes(), "player")
		assert.Equal(t, "alice", m["username"])
		assert.Equal(t, "alice@realworld.io", m["email"])
		assert.Nil(t, m["bio"])
		assert.Nil(t, m["image"])
		assert.NotEmpty(t, m["token"])
	}
}

func TestLoginCaseSuccess(t *testing.T) {
	tearDown()
	setup()
	var (
		reqJSON = `{"player":{"email":"player1@realworld.io","password":"secret"}}`
	)
	req := httptest.NewRequest(echo.POST, "/api/players/login", strings.NewReader(reqJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	assert.NoError(t, h.Login(c))
	if assert.Equal(t, http.StatusOK, rec.Code) {
		m := responseMap(rec.Body.Bytes(), "player")
		assert.Equal(t, "player1", m["username"])
		assert.Equal(t, "player1@realworld.io", m["email"])
		assert.NotEmpty(t, m["token"])
	}
}

func TestLoginCaseFailed(t *testing.T) {
	tearDown()
	setup()
	var (
		reqJSON = `{"player":{"email":"playerx@realworld.io","password":"secret"}}`
	)
	req := httptest.NewRequest(echo.POST, "/api/players/login", strings.NewReader(reqJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	assert.NoError(t, h.Login(c))
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestCurrentPlayerCaseSuccess(t *testing.T) {
	tearDown()
	setup()
	jwtMiddleware := middleware.JWT(utils.JWTSecret)
	req := httptest.NewRequest(echo.GET, "/api/players/login", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, authHeader(utils.GenerateJWT(1)))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := jwtMiddleware(func(context echo.Context) error {
		return h.CurrentPlayer(c)
	})(c)
	assert.NoError(t, err)
	if assert.Equal(t, http.StatusOK, rec.Code) {
		m := responseMap(rec.Body.Bytes(), "player")
		assert.Equal(t, "player1", m["username"])
		assert.Equal(t, "player1@realworld.io", m["email"])
		assert.NotEmpty(t, m["token"])
	}
}

func TestCurrentPlayerCaseInvalid(t *testing.T) {
	tearDown()
	setup()
	jwtMiddleware := middleware.JWT(utils.JWTSecret)
	req := httptest.NewRequest(echo.GET, "/api/players/login", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, authHeader(utils.GenerateJWT(100)))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := jwtMiddleware(func(context echo.Context) error {
		return h.CurrentPlayer(c)
	})(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestUpdatePlayerEmail(t *testing.T) {
	tearDown()
	setup()
	var (
		player1UpdateReq = `{"player":{"email":"player1@player1.me"}}`
	)
	jwtMiddleware := middleware.JWT(utils.JWTSecret)
	req := httptest.NewRequest(echo.PUT, "/api/player", strings.NewReader(player1UpdateReq))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, authHeader(utils.GenerateJWT(1)))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := jwtMiddleware(func(context echo.Context) error {
		return h.UpdatePlayer(c)
	})(c)
	assert.NoError(t, err)
	if assert.Equal(t, http.StatusOK, rec.Code) {
		m := responseMap(rec.Body.Bytes(), "player")
		assert.Equal(t, "player1", m["username"])
		assert.Equal(t, "player1@player1.me", m["email"])
		assert.NotEmpty(t, m["token"])
	}
}

func TestUpdatePlayerMultipleFields(t *testing.T) {
	tearDown()
	setup()
	var (
		player1UpdateReq = `{"player":{"username":"player11","email":"player11@player11.me","bio":"player11 bio"}}`
	)
	jwtMiddleware := middleware.JWT(utils.JWTSecret)
	req := httptest.NewRequest(echo.PUT, "/api/player", strings.NewReader(player1UpdateReq))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, authHeader(utils.GenerateJWT(1)))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := jwtMiddleware(func(context echo.Context) error {
		return h.UpdatePlayer(c)
	})(c)
	assert.NoError(t, err)
	if assert.Equal(t, http.StatusOK, rec.Code) {
		m := responseMap(rec.Body.Bytes(), "player")
		assert.Equal(t, "player11", m["username"])
		assert.Equal(t, "player11@player11.me", m["email"])
		assert.Equal(t, "player11 bio", m["bio"])
		assert.NotEmpty(t, m["token"])
	}
}

func TestGetProfileCaseSuccess(t *testing.T) {
	tearDown()
	setup()
	jwtMiddleware := middleware.JWT(utils.JWTSecret)
	req := httptest.NewRequest(echo.GET, "/api/profiles/:username", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, authHeader(utils.GenerateJWT(1)))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/profiles/:username")
	c.SetParamNames("username")
	c.SetParamValues("player1")
	err := jwtMiddleware(func(context echo.Context) error {
		return h.GetProfile(c)
	})(c)
	assert.NoError(t, err)
	if assert.Equal(t, http.StatusOK, rec.Code) {
		m := responseMap(rec.Body.Bytes(), "profile")
		assert.Equal(t, "player1", m["username"])
		assert.Equal(t, "player1 bio", m["bio"])
		assert.Equal(t, "http://realworld.io/player1.jpg", m["image"])
		assert.Equal(t, false, m["following"])
	}
}

func TestGetProfileCaseNotFound(t *testing.T) {
	tearDown()
	setup()
	jwtMiddleware := middleware.JWT(utils.JWTSecret)
	req := httptest.NewRequest(echo.GET, "/api/profiles/:username", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, authHeader(utils.GenerateJWT(1)))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/profiles/:username")
	c.SetParamNames("username")
	c.SetParamValues("playerx")
	err := jwtMiddleware(func(context echo.Context) error {
		return h.GetProfile(c)
	})(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestFollowCaseSuccess(t *testing.T) {
	tearDown()
	setup()
	jwtMiddleware := middleware.JWT(utils.JWTSecret)
	req := httptest.NewRequest(echo.POST, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, authHeader(utils.GenerateJWT(1)))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/profiles/:username/follow")
	c.SetParamNames("username")
	c.SetParamValues("player2")
	err := jwtMiddleware(func(context echo.Context) error {
		return h.Follow(c)
	})(c)
	assert.NoError(t, err)
	if assert.Equal(t, http.StatusOK, rec.Code) {
		m := responseMap(rec.Body.Bytes(), "profile")
		assert.Equal(t, "player2", m["username"])
		assert.Equal(t, "player2 bio", m["bio"])
		assert.Equal(t, "http://realworld.io/player2.jpg", m["image"])
		assert.Equal(t, true, m["following"])
	}
}

func TestFollowCaseInvalidPlayer(t *testing.T) {
	tearDown()
	setup()
	jwtMiddleware := middleware.JWT(utils.JWTSecret)
	req := httptest.NewRequest(echo.POST, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, authHeader(utils.GenerateJWT(1)))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/profiles/:username/follow")
	c.SetParamNames("username")
	c.SetParamValues("playerx")
	err := jwtMiddleware(func(context echo.Context) error {
		return h.Follow(c)
	})(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestUnfollow(t *testing.T) {
	tearDown()
	setup()
	jwtMiddleware := middleware.JWT(utils.JWTSecret)
	req := httptest.NewRequest(echo.DELETE, "/api/profiles/:username/follow", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, authHeader(utils.GenerateJWT(1)))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/profiles/:username/follow")
	c.SetParamNames("username")
	c.SetParamValues("player2")
	err := jwtMiddleware(func(context echo.Context) error {
		return h.Unfollow(c)
	})(c)
	assert.NoError(t, err)
	if assert.Equal(t, http.StatusOK, rec.Code) {
		m := responseMap(rec.Body.Bytes(), "profile")
		assert.Equal(t, "player2", m["username"])
		assert.Equal(t, "player2 bio", m["bio"])
		assert.Equal(t, "http://realworld.io/player2.jpg", m["image"])
		assert.Equal(t, false, m["following"])
	}
}
