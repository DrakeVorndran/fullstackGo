package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"golang-starter-pack/model"
	"golang-starter-pack/utils"
)

func (h *Handler) GetItem(c echo.Context) error {
	slug := c.Param("slug")
	a, err := h.itemStore.GetBySlug(slug)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if a == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	return c.JSON(http.StatusOK, newItemResponse(c, a))
}

func (h *Handler) Items(c echo.Context) error {
	tag := c.QueryParam("tag")
	author := c.QueryParam("author")
	favoritedBy := c.QueryParam("favorited")
	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil {
		offset = 0
	}
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 20
	}
	var items []model.Item
	var count int
	if tag != "" {
		items, count, err = h.itemStore.ListByTag(tag, offset, limit)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, nil)
		}
	} else if author != "" {
		items, count, err = h.itemStore.ListByAuthor(author, offset, limit)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, nil)
		}
	} else if favoritedBy != "" {
		items, count, err = h.itemStore.ListByWhoFavorited(favoritedBy, offset, limit)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, nil)
		}
	} else {
		items, count, err = h.itemStore.List(offset, limit)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, nil)
		}
	}
	return c.JSON(http.StatusOK, newItemListResponse(h.playerStore, playerIDFromToken(c), items, count))
}

func (h *Handler) Feed(c echo.Context) error {
	var items []model.Item
	var count int
	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil {
		offset = 0
	}
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 20
	}
	items, count, err = h.itemStore.ListFeed(playerIDFromToken(c), offset, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, newItemListResponse(h.playerStore, playerIDFromToken(c), items, count))
}

func (h *Handler) CreateItem(c echo.Context) error {
	var a model.Item
	req := &itemCreateRequest{}
	if err := req.bind(c, &a); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	a.AuthorID = playerIDFromToken(c)
	err := h.itemStore.CreateItem(&a)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	return c.JSON(http.StatusCreated, newItemResponse(c, &a))
}

func (h *Handler) UpdateItem(c echo.Context) error {
	slug := c.Param("slug")
	a, err := h.itemStore.GetPlayerItemBySlug(playerIDFromToken(c), slug)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if a == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	req := &itemUpdateRequest{}
	req.populate(a)
	if err := req.bind(c, a); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	if err = h.itemStore.UpdateItem(a, req.Items.Tags); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	return c.JSON(http.StatusOK, newItemResponse(c, a))
}

func (h *Handler) DeleteItem(c echo.Context) error {
	slug := c.Param("slug")
	a, err := h.itemStore.GetPlayerItemBySlug(playerIDFromToken(c), slug)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if a == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	err = h.itemStore.DeleteItem(a)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"result": "ok"})
}

func (h *Handler) AddComment(c echo.Context) error {
	slug := c.Param("slug")
	a, err := h.itemStore.GetBySlug(slug)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if a == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	var cm model.Comment
	req := &createCommentRequest{}
	if err := req.bind(c, &cm); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	if err = h.itemStore.AddComment(a, &cm); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	return c.JSON(http.StatusCreated, newCommentResponse(c, &cm))
}

func (h *Handler) GetComments(c echo.Context) error {
	slug := c.Param("slug")
	cm, err := h.itemStore.GetCommentsBySlug(slug)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	return c.JSON(http.StatusOK, newCommentListResponse(c, cm))
}

func (h *Handler) DeleteComment(c echo.Context) error {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 32)
	id := uint(id64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewError(err))
	}
	cm, err := h.itemStore.GetCommentByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if cm == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	if cm.PlayerID != playerIDFromToken(c) {
		return c.JSON(http.StatusUnauthorized, utils.NewError(errors.New("unauthorized action")))
	}
	if err := h.itemStore.DeleteComment(cm); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"result": "ok"})
}

func (h *Handler) Favorite(c echo.Context) error {
	slug := c.Param("slug")
	a, err := h.itemStore.GetBySlug(slug)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if a == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	if err := h.itemStore.AddFavorite(a, playerIDFromToken(c)); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	return c.JSON(http.StatusOK, newItemResponse(c, a))
}

func (h *Handler) Unfavorite(c echo.Context) error {
	slug := c.Param("slug")
	a, err := h.itemStore.GetBySlug(slug)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if a == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	if err := h.itemStore.RemoveFavorite(a, playerIDFromToken(c)); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	return c.JSON(http.StatusOK, newItemResponse(c, a))
}

func (h *Handler) Tags(c echo.Context) error {
	tags, err := h.itemStore.ListTags()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, newTagListResponse(tags))
}
