package store

import (
	"github.com/jinzhu/gorm"
	"golang-starter-pack/model"
)

type ItemStore struct {
	db *gorm.DB
}

func NewItemStore(db *gorm.DB) *ItemStore {
	return &ItemStore{
		db: db,
	}
}

func (as *ItemStore) GetBySlug(s string) (*model.Item, error) {
	var m model.Item
	err := as.db.Where(&model.Item{Slug: s}).Preload("Favorites").Preload("Tags").Preload("Author").Find(&m).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &m, err
}

func (as *ItemStore) GetPlayerItemBySlug(playerID uint, slug string) (*model.Item, error) {
	var m model.Item
	err := as.db.Where(&model.Item{Slug: slug, AuthorID: playerID}).Find(&m).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &m, err
}

func (as *ItemStore) CreateItem(a *model.Item) error {
	tags := a.Tags
	tx := as.db.Begin()
	if err := tx.Create(&a).Error; err != nil {
		return err
	}
	for _, t := range a.Tags {
		err := tx.Where(&model.Tag{Tag: t.Tag}).First(&t).Error
		if err != nil && !gorm.IsRecordNotFoundError(err) {
			tx.Rollback()
			return err
		}
		if err := tx.Model(&a).Association("Tags").Append(t).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := tx.Where(a.ID).Preload("Favorites").Preload("Tags").Preload("Author").Find(&a).Error; err != nil {
		tx.Rollback()
		return err
	}
	a.Tags = tags
	return tx.Commit().Error
}

func (as *ItemStore) UpdateItem(a *model.Item, tagList []string) error {
	tx := as.db.Begin()
	if err := tx.Model(a).Update(a).Error; err != nil {
		return err
	}
	tags := make([]model.Tag, 0)
	for _, t := range tagList {
		tag := model.Tag{Tag: t}
		err := tx.Where(&tag).First(&tag).Error
		if err != nil && !gorm.IsRecordNotFoundError(err) {
			tx.Rollback()
			return err
		}
		tags = append(tags, tag)
	}
	if err := tx.Model(a).Association("Tags").Replace(tags).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where(a.ID).Preload("Favorites").Preload("Tags").Preload("Author").Find(a).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (as *ItemStore) DeleteItem(a *model.Item) error {
	return as.db.Delete(a).Error
}

func (as *ItemStore) List(offset, limit int) ([]model.Item, int, error) {
	var (
		items []model.Item
		count int
	)
	as.db.Model(&items).Count(&count)
	as.db.Preload("Favorites").Preload("Tags").Preload("Author").Offset(offset).Limit(limit).Order("created_at desc").Find(&items)
	return items, count, nil
}

func (as *ItemStore) ListByTag(tag string, offset, limit int) ([]model.Item, int, error) {
	var (
		t     model.Tag
		items []model.Item
		count int
	)
	err := as.db.Where(&model.Tag{Tag: tag}).First(&t).Error
	if err != nil {
		return nil, 0, err
	}
	as.db.Model(&t).Preload("Favorites").Preload("Tags").Preload("Author").Offset(offset).Limit(limit).Order("created_at desc").Association("Items").Find(&items)
	count = as.db.Model(&t).Association("Items").Count()
	return items, count, nil
}

func (as *ItemStore) ListByAuthor(username string, offset, limit int) ([]model.Item, int, error) {
	var (
		u     model.Player
		items []model.Item
		count int
	)
	err := as.db.Where(&model.Player{Username: username}).First(&u).Error
	if err != nil {
		return nil, 0, err
	}
	as.db.Where(&model.Item{AuthorID: u.ID}).Preload("Favorites").Preload("Tags").Preload("Author").Offset(offset).Limit(limit).Order("created_at desc").Find(&items)
	as.db.Where(&model.Item{AuthorID: u.ID}).Model(&model.Item{}).Count(&count)

	return items, count, nil
}

func (as *ItemStore) ListByWhoFavorited(username string, offset, limit int) ([]model.Item, int, error) {
	var (
		u     model.Player
		items []model.Item
		count int
	)
	err := as.db.Where(&model.Player{Username: username}).First(&u).Error
	if err != nil {
		return nil, 0, err
	}
	as.db.Model(&u).Preload("Favorites").Preload("Tags").Preload("Author").Offset(offset).Limit(limit).Order("created_at desc").Association("Favorites").Find(&items)
	count = as.db.Model(&u).Association("Favorites").Count()
	return items, count, nil
}

func (as *ItemStore) ListFeed(playerID uint, offset, limit int) ([]model.Item, int, error) {
	var (
		u     model.Player
		items []model.Item
		count int
	)
	err := as.db.First(&u, playerID).Error
	if err != nil {
		return nil, 0, err
	}
	var followings []model.Follow
	as.db.Model(&u).Preload("Following").Preload("Follower").Association("Followings").Find(&followings)
	var ids []uint
	for _, i := range followings {
		ids = append(ids, i.FollowingID)
	}
	as.db.Where("author_id in (?)", ids).Preload("Favorites").Preload("Tags").Preload("Author").Offset(offset).Limit(limit).Order("created_at desc").Find(&items)
	as.db.Where(&model.Item{AuthorID: u.ID}).Model(&model.Item{}).Count(&count)
	return items, count, nil
}

func (as *ItemStore) AddComment(a *model.Item, c *model.Comment) error {
	err := as.db.Model(a).Association("Comments").Append(c).Error
	if err != nil {
		return err
	}
	return as.db.Where(c.ID).Preload("Player").First(c).Error
}

func (as *ItemStore) GetCommentsBySlug(slug string) ([]model.Comment, error) {
	var m model.Item
	if err := as.db.Where(&model.Item{Slug: slug}).Preload("Comments").Preload("Comments.Player").First(&m).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return m.Comments, nil
}

func (as *ItemStore) GetCommentByID(id uint) (*model.Comment, error) {
	var m model.Comment
	if err := as.db.Where(id).First(&m).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (as *ItemStore) DeleteComment(c *model.Comment) error {
	return as.db.Delete(c).Error
}

func (as *ItemStore) AddFavorite(a *model.Item, playerID uint) error {
	usr := model.Player{}
	usr.ID = playerID
	return as.db.Model(a).Association("Favorites").Append(&usr).Error
}

func (as *ItemStore) RemoveFavorite(a *model.Item, playerID uint) error {
	usr := model.Player{}
	usr.ID = playerID
	return as.db.Model(a).Association("Favorites").Delete(&usr).Error
}

func (as *ItemStore) ListTags() ([]model.Tag, error) {
	var tags []model.Tag
	if err := as.db.Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}
