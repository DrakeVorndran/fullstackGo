package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"golang-starter-pack/router"
	"golang-starter-pack/router/middleware"
	"golang-starter-pack/utils"
)

func TestListItemsCaseSuccess(t *testing.T) {
	tearDown()
	setup()
	e := router.New()
	req := httptest.NewRequest(echo.GET, "/api/items", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	assert.NoError(t, h.Items(c))
	if assert.Equal(t, http.StatusOK, rec.Code) {
		var aa itemListResponse
		err := json.Unmarshal(rec.Body.Bytes(), &aa)
		assert.NoError(t, err)
		assert.Equal(t, 2, aa.ItemsCount)
	}
}

func TestGetItemsCaseSuccess(t *testing.T) {
	tearDown()
	setup()
	req := httptest.NewRequest(echo.GET, "/api/items/:slug", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/items/:slug")
	c.SetParamNames("slug")
	c.SetParamValues("item1-slug")
	assert.NoError(t, h.GetItem(c))
	if assert.Equal(t, http.StatusOK, rec.Code) {
		var a singleItemResponse
		err := json.Unmarshal(rec.Body.Bytes(), &a)
		assert.NoError(t, err)
		assert.Equal(t, "item1-slug", a.Item.Slug)
		assert.Equal(t, 2, len(a.Item.TagList))
	}
}

func TestCreateItemsCaseSuccess(t *testing.T) {
	tearDown()
	setup()
	var (
		reqJSON = `{"item":{"title":"item2", "description":"item2", "body":"item2", "tagList":["tag1","tag2"]}}`
	)
	jwtMiddleware := middleware.JWT(utils.JWTSecret)
	req := httptest.NewRequest(echo.POST, "/api/items", strings.NewReader(reqJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, authHeader(utils.GenerateJWT(1)))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := jwtMiddleware(func(context echo.Context) error {
		return h.CreateItem(c)
	})(c)
	assert.NoError(t, err)
	if assert.Equal(t, http.StatusCreated, rec.Code) {
		var a singleItemResponse
		err := json.Unmarshal(rec.Body.Bytes(), &a)
		assert.NoError(t, err)
		assert.Equal(t, "item2", a.Item.Slug)
		assert.Equal(t, "item2", a.Item.Description)
		assert.Equal(t, "item2", a.Item.Title)
		assert.Equal(t, "player1", a.Item.Author.Username)
		assert.Equal(t, 2, len(a.Item.TagList))
	}
}

func TestUpdateItemsCaseSuccess(t *testing.T) {
	tearDown()
	setup()
	var (
		reqJSON = `{"item":{"title":"item1 part 2", "tagList":["tag3"]}}`
	)
	jwtMiddleware := middleware.JWT(utils.JWTSecret)
	req := httptest.NewRequest(echo.PUT, "/api/items/:slug", strings.NewReader(reqJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, authHeader(utils.GenerateJWT(1)))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/items/:slug")
	c.SetParamNames("slug")
	c.SetParamValues("item1-slug")
	err := jwtMiddleware(func(context echo.Context) error {
		return h.UpdateItem(c)
	})(c)
	assert.NoError(t, err)
	if assert.Equal(t, http.StatusOK, rec.Code) {
		var a singleItemResponse
		err := json.Unmarshal(rec.Body.Bytes(), &a)
		assert.NoError(t, err)
		assert.Equal(t, "item1 part 2", a.Item.Title)
		assert.Equal(t, "item1-part-2", a.Item.Slug)
		assert.Equal(t, 1, len(a.Item.TagList))
		assert.Equal(t, "tag3", a.Item.TagList[0])
	}
}

func TestFeedCaseSuccess(t *testing.T) {
	tearDown()
	setup()
	jwtMiddleware := middleware.JWT(utils.JWTSecret)
	req := httptest.NewRequest(echo.GET, "/api/items/feed", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, authHeader(utils.GenerateJWT(1)))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := jwtMiddleware(func(context echo.Context) error {
		return h.Feed(c)
	})(c)
	assert.NoError(t, err)
	if assert.Equal(t, http.StatusOK, rec.Code) {
		var a itemListResponse
		err := json.Unmarshal(rec.Body.Bytes(), &a)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(a.Items))
		assert.Equal(t, a.ItemsCount, len(a.Items))
		assert.Equal(t, "item2 title", a.Items[0].Title)
		assert.Equal(t, "item2 title", a.Items[0].Title)
	}
}

func TestDeleteItemCaseSuccess(t *testing.T) {
	tearDown()
	setup()
	jwtMiddleware := middleware.JWT(utils.JWTSecret)
	req := httptest.NewRequest(echo.DELETE, "/api/items/:slug", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, authHeader(utils.GenerateJWT(1)))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/items/:slug")
	c.SetParamNames("slug")
	c.SetParamValues("item1-slug")
	err := jwtMiddleware(func(context echo.Context) error {
		return h.DeleteItem(c)
	})(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetCommentsCaseSuccess(t *testing.T) {
	tearDown()
	setup()
	jwtMiddleware := middleware.JWT(utils.JWTSecret)
	req := httptest.NewRequest(echo.GET, "/api/items/:slug/comments", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, authHeader(utils.GenerateJWT(2)))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/items/:slug/comments")
	c.SetParamNames("slug")
	c.SetParamValues("item1-slug")
	err := jwtMiddleware(func(context echo.Context) error {
		return h.GetComments(c)
	})(c)
	assert.NoError(t, err)
	if assert.Equal(t, http.StatusOK, rec.Code) {
		var cc commentListResponse
		err := json.Unmarshal(rec.Body.Bytes(), &cc)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(cc.Comments))
	}
}

func TestAddCommentCaseSuccess(t *testing.T) {
	tearDown()
	setup()
	var (
		reqJSON = `{"comment":{"body":"item1 comment2 by player2"}}`
	)
	jwtMiddleware := middleware.JWT(utils.JWTSecret)
	req := httptest.NewRequest(echo.POST, "/api/items/:slug/comments", strings.NewReader(reqJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, authHeader(utils.GenerateJWT(2)))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/items/:slug/comments")
	c.SetParamNames("slug")
	c.SetParamValues("item1-slug")
	err := jwtMiddleware(func(context echo.Context) error {
		return h.AddComment(c)
	})(c)
	assert.NoError(t, err)
	if assert.Equal(t, http.StatusCreated, rec.Code) {
		var c singleCommentResponse
		err := json.Unmarshal(rec.Body.Bytes(), &c)
		assert.NoError(t, err)
		assert.Equal(t, "item1 comment2 by player2", c.Comment.Body)
		assert.Equal(t, "player2", c.Comment.Author.Username)
	}
}

func TestDeleteCommentCaseSuccess(t *testing.T) {
	tearDown()
	setup()
	jwtMiddleware := middleware.JWT(utils.JWTSecret)
	req := httptest.NewRequest(echo.DELETE, "/api/items/:slug/comments/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, authHeader(utils.GenerateJWT(1)))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/items/:slug/comments/:id")
	c.SetParamNames("slug")
	c.SetParamValues("item1-slug")
	c.SetParamNames("id")
	c.SetParamValues("1")
	err := jwtMiddleware(func(context echo.Context) error {
		return h.DeleteComment(c)
	})(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestFavoriteCaseSuccess(t *testing.T) {
	tearDown()
	setup()
	jwtMiddleware := middleware.JWT(utils.JWTSecret)
	req := httptest.NewRequest(echo.POST, "/api/items/:slug/favorite", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, authHeader(utils.GenerateJWT(2)))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/items/:slug/comments")
	c.SetParamNames("slug")
	c.SetParamValues("item1-slug")
	err := jwtMiddleware(func(context echo.Context) error {
		return h.Favorite(c)
	})(c)
	assert.NoError(t, err)
	if assert.Equal(t, http.StatusOK, rec.Code) {
		var a singleItemResponse
		err := json.Unmarshal(rec.Body.Bytes(), &a)
		assert.NoError(t, err)
		assert.Equal(t, "item1 title", a.Item.Title)
		assert.True(t, a.Item.Favorited)
		assert.Equal(t, 1, a.Item.FavoritesCount)
	}
}

func TestUnfavoriteCaseSuccess(t *testing.T) {
	tearDown()
	setup()
	jwtMiddleware := middleware.JWT(utils.JWTSecret)
	req := httptest.NewRequest(echo.DELETE, "/api/items/:slug/favorite", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, authHeader(utils.GenerateJWT(1)))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/items/:slug/favorite")
	c.SetParamNames("slug")
	c.SetParamValues("item2-slug")
	err := jwtMiddleware(func(context echo.Context) error {
		return h.Unfavorite(c)
	})(c)
	assert.NoError(t, err)
	if assert.Equal(t, http.StatusOK, rec.Code) {
		var a singleItemResponse
		err := json.Unmarshal(rec.Body.Bytes(), &a)
		assert.NoError(t, err)
		assert.Equal(t, "item2 title", a.Item.Title)
		assert.False(t, a.Item.Favorited)
		assert.Equal(t, 0, a.Item.FavoritesCount)
	}
}

func TestGetTagsCaseSuccess(t *testing.T) {
	tearDown()
	setup()
	req := httptest.NewRequest(echo.GET, "/api/tags", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	assert.NoError(t, h.Tags(c))
	if assert.Equal(t, http.StatusOK, rec.Code) {
		var tt tagListResponse
		err := json.Unmarshal(rec.Body.Bytes(), &tt)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(tt.Tags))
		assert.Contains(t, tt.Tags, "tag1")
		assert.Contains(t, tt.Tags, "tag2")
	}
}
