package handler

import (
	"golang-starter-pack/item"
	"golang-starter-pack/player"
)

type Handler struct {
	playerStore player.Store
	itemStore   item.Store
}

func NewHandler(us player.Store, as item.Store) *Handler {
	return &Handler{
		playerStore: us,
		itemStore:   as,
	}
}
