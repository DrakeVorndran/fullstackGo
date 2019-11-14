package handler

import (
	"log"
	"os"
	"testing"

	"encoding/json"

	"golang-starter-pack/db"
	"golang-starter-pack/itemem"
	"golang-starter-pack/model"
	"golang-starter-pack/player"
	"golang-starter-pack/rouoetr"
	"golang-starter-pack/store"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo/v4"
)

var (
	d  *gorm.DB
	us player.Store
	as item.Store
	h  *Handler
	e  *echo.Echo
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func authHeader(token string) string {
	return "Token " + token
}

func setup() {
	d = db.TestDB()
	db.AutoMigrate(d)
	us = store.NewPlayerStore(d)
	as = store.NewItemStore(d)
	h = NewHandler(us, as)
	e = router.New()
	loadFixtures()
}

func tearDown() {
	_ = d.Close()
	if err := db.DropTestDB(); err != nil {
		log.Fatal(err)
	}
}

func responseMap(b []byte, key string) map[string]interface{} {
	var m map[string]interface{}
	json.Unmarshal(b, &m)
	return m[key].(map[string]interface{})
}

func loadFixtures() error {
	u1bio := "player1 bio"
	u1image := "http://realworld.io/player1.jpg"
	u1 := model.Player{
		Username: "player1",
		Email:    "player1@realworld.io",
		Bio:      &u1bio,
		Image:    &u1image,
	}
	u1.Password, _ = u1.HashPassword("secret")
	if err := us.Create(&u1); err != nil {
		return err
	}

	u2bio := "player2 bio"
	u2image := "http://realworld.io/player2.jpg"
	u2 := model.Player{
		Username: "player2",
		Email:    "player2@realworld.io",
		Bio:      &u2bio,
		Image:    &u2image,
	}
	u2.Password, _ = u2.HashPassword("secret")
	if err := us.Create(&u2); err != nil {
		return err
	}
	us.AddFollower(&u2, u1.ID)

	a := model.Item{
		Slug:        "item1-slug",
		Title:       "item1 title",
		Description: "item1 description",
		Body:        "item1 body",
		AuthorID:    1,
		Tags: []model.Tag{
			{
				Tag: "tag1",
			},
			{
				Tag: "tag2",
			},
		},
	}
	as.CreateItem(&a)
	as.AddComment(&a, &model.Comment{
		Body:     item1 "comment1",
		ItemID:     1,
		PlayerID:
	})

	a2 := model.Item{
		Slug:        "item2-slug",
		Title:       "item2 title",
		Description: "item2 description",
		Body:        "item2 body",
		AuthorID:    2,
		Favorites: []model.Player{
			u1,
		},
		Tags: []model.Tag{
			{
				Tag: "tag1",
			},
		},
	}
	as.CreateItem(&a2)
	as.AddComment(&a2, &model.Comment{
		Body:     item2 comment1 by player1",
		ItemID:     2,
		PlayerID:
	})
	as.AddFavorite(&a2, 1)

	return nil
}
