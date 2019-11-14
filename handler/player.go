package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"golang-starter-pack/model"
	"golang-starter-pack/utils"
)

func (h *Handler) SignUp(c echo.Context) error {
	var u model.Player
	req := &playerRegisterRequest{}
	if err := req.bind(c, &u); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	if err := h.playerStore.Create(&u); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	return c.JSON(http.StatusCreated, newPlayerResponse(&u))
}

func (h *Handler) Login(c echo.Context) error {
	req := &playerLoginRequest{}
	if err := req.bind(c); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	u, err := h.playerStore.GetByEmail(req.Player.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if u == nil {
		return c.JSON(http.StatusForbidden, utils.AccessForbidden())
	}
	if !u.CheckPassword(req.Player.Password) {
		return c.JSON(http.StatusForbidden, utils.AccessForbidden())
	}
	return c.JSON(http.StatusOK, newPlayerResponse(u))
}

func (h *Handler) CurrentPlayer(c echo.Context) error {
	u, err := h.playerStore.GetByID(playerIDFromToken(c))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if u == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	return c.JSON(http.StatusOK, newPlayerResponse(u))
}

func (h *Handler) UpdatePlayer(c echo.Context) error {
	u, err := h.playerStore.GetByID(playerIDFromToken(c))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if u == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	req := newPlayerUpdateRequest()
	req.populate(u)
	if err := req.bind(c, u); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	if err := h.playerStore.Update(u); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	return c.JSON(http.StatusOK, newPlayerResponse(u))
}

func (h *Handler) GetProfile(c echo.Context) error {
	username := c.Param("username")
	u, err := h.playerStore.GetByUsername(username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if u == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	return c.JSON(http.StatusOK, newProfileResponse(h.playerStore, playerIDFromToken(c), u))
}

func (h *Handler) Follow(c echo.Context) error {
	followerID := playerIDFromToken(c)
	username := c.Param("username")
	u, err := h.playerStore.GetByUsername(username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if u == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	if err := h.playerStore.AddFollower(u, followerID); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	return c.JSON(http.StatusOK, newProfileResponse(h.playerStore, playerIDFromToken(c), u))
}
func (h *Handler) Unfollow(c echo.Context) error {
	followerID := playerIDFromToken(c)
	username := c.Param("username")
	u, err := h.playerStore.GetByUsername(username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}
	if u == nil {
		return c.JSON(http.StatusNotFound, utils.NotFound())
	}
	if err := h.playerStore.RemoveFollower(u, followerID); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}
	return c.JSON(http.StatusOK, newProfileResponse(h.playerStore, playerIDFromToken(c), u))
}
func playerIDFromToken(c echo.Context) uint {
	id, ok := c.Get("player").(uint)
	if !ok {
		return 0
	}
	return id
}
