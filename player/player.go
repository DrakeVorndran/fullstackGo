package player

import (
	"golang-starter-pack/model"
)

type Store interface {
	GetByID(uint) (*model.Player, error)
	GetByEmail(string) (*model.Player, error)
	GetByUsername(string) (*model.Player, error)
	Create(*model.Player) error
	Update(*model.Player) error
	AddFollower(player *model.Player, followerID uint) error
	RemoveFollower(player *model.Player, followerID uint) error
	IsFollower(playerID, followerID uint) (bool, error)
}
