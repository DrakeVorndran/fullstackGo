package item

import (
	"golang-starter-pack/model"
)

type Store interface {
	GetBySlug(string) (*model.Item, error)
	GetPlayerItemBySlug(playerID uint, slug string) (*model.Item, error)
	CreateItem(*model.Item) error
	UpdateItem(*model.Item, []string) error
	DeleteItem(*model.Item) error
	List(offset, limit int) ([]model.Item, int, error)
	ListByTag(tag string, offset, limit int) ([]model.Item, int, error)
	ListByAuthor(username string, offset, limit int) ([]model.Item, int, error)
	ListByWhoFavorited(username string, offset, limit int) ([]model.Item, int, error)
	ListFeed(playerID uint, offset, limit int) ([]model.Item, int, error)

	AddComment(*model.Item, *model.Comment) error
	GetCommentsBySlug(string) ([]model.Comment, error)
	GetCommentByID(uint) (*model.Comment, error)
	DeleteComment(*model.Comment) error

	AddFavorite(*model.Item, uint) error
	RemoveFavorite(*model.Item, uint) error
	ListTags() ([]model.Tag, error)
}
