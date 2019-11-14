package model

import (
	"errors"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type Player struct {
	gorm.Model
	Username   string `gorm:"unique_index;not null"`
	Email      string `gorm:"unique_index;not null"`
	Password   string `gorm:"not null"`
	Bio        *string
	Image      *string
	Followers  []Follow `gorm:"foreignkey:FollowingID"`
	Followings []Follow `gorm:"foreignkey:FollowerID"`
	Favorites  []Item   `gorm:"many2many:favorites;"`
}

type Follow struct {
	Follower    Player
	FollowerID  uint `gorm:"primary_key" sql:"type:int not null"`
	Following   Player
	FollowingID uint `gorm:"primary_key" sql:"type:int not null"`
}

func (p *Player) HashPassword(plain string) (string, error) {
	if len(plain) == 0 {
		return "", errors.New("password should not be empty")
	}
	h, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	return string(h), err
}

func (u *Player) CheckPassword(plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plain))
	return err == nil
}

// FollowedBy Followings should be pre loaded
func (u *Player) FollowedBy(id uint) bool {
	if u.Followers == nil {
		return false
	}
	for _, f := range u.Followers {
		if f.FollowerID == id {
			return true
		}
	}
	return false
}
