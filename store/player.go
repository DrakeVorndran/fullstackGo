package store

import (
	"github.com/jinzhu/gorm"
	"golang-starter-pack/model"
)

type PlayerStore struct {
	db *gorm.DB
}

func NewPlayerStore(db *gorm.DB) *PlayerStore {
	return &PlayerStore{
		db: db,
	}
}

func (us *PlayerStore) GetByID(id uint) (*model.Player, error) {
	var m model.Player
	if err := us.db.First(&m, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (us *PlayerStore) GetByEmail(e string) (*model.Player, error) {
	var m model.Player
	if err := us.db.Where(&model.Player{Email: e}).First(&m).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (us *PlayerStore) GetByUsername(username string) (*model.Player, error) {
	var m model.Player
	if err := us.db.Where(&model.Player{Username: username}).Preload("Followers").First(&m).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (us *PlayerStore) Create(u *model.Player) (err error) {
	return us.db.Create(u).Error
}

func (us *PlayerStore) Update(u *model.Player) error {
	return us.db.Model(u).Update(u).Error
}

func (us *PlayerStore) AddFollower(u *model.Player, followerID uint) error {
	return us.db.Model(u).Association("Followers").Append(&model.Follow{FollowerID: followerID, FollowingID: u.ID}).Error
}

func (us *PlayerStore) RemoveFollower(u *model.Player, followerID uint) error {
	f := model.Follow{
		FollowerID:  followerID,
		FollowingID: u.ID,
	}
	if err := us.db.Model(u).Association("Followers").Find(&f).Error; err != nil {
		return err
	}
	if err := us.db.Delete(f).Error; err != nil {
		return err
	}
	return nil
}

func (us *PlayerStore) IsFollower(playerID, followerID uint) (bool, error) {
	var f model.Follow
	if err := us.db.Where("following_id = ? AND follower_id = ?", playerID, followerID).Find(&f).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
