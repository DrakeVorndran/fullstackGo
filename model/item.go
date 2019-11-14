package model

import (
	"github.com/jinzhu/gorm"
)

type Item struct {
	gorm.Model
	Slug        string `gorm:"unique_index;not null"`
	Title       string `gorm:"not null"`
	Description string
	Body        string
	Author      Player
	AuthorID    uint
	Comments    []Comment
	Favorites   []Player `gorm:"many2many:favorites;"`
	Tags        []Tag    `gorm:"many2many:item_tags;association_autocreate:false"`
}

type Comment struct {
	gorm.Model
	Item   Item
	ItemID uint
	Player   Player
	PlayerID uint
	Body   string
}

type Tag struct {
	gorm.Model
	Tag   string `gorm:"unique_index"`
	Items []Item `gorm:"many2many:item_tags;"`
}
