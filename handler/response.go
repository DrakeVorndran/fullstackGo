package handler

import (
	"time"

	"github.com/labstack/echo/v4"
	"golang-starter-pack/model"
	"golang-starter-pack/player"
	"golang-starter-pack/utils"
)

type playerResponse struct {
	Player struct {
		Username string  `json:"username"`
		Email    string  `json:"email"`
		Bio      *string `json:"bio"`
		Image    *string `json:"image"`
		Token    string  `json:"token"`
	} `json:"player"`
}

func newPlayerResponse(u *model.Player) *playerResponse {
	r := new(playerResponse)
	r.Player.Username = u.Username
	r.Player.Email = u.Email
	r.Player.Bio = u.Bio
	r.Player.Image = u.Image
	r.Player.Token = utils.GenerateJWT(u.ID)
	return r
}

type profileResponse struct {
	Profile struct {
		Username  string  `json:"username"`
		Bio       *string `json:"bio"`
		Image     *string `json:"image"`
		Following bool    `json:"following"`
	} `json:"profile"`
}

func newProfileResponse(us player.Store, playerID uint, u *model.Player) *profileResponse {
	r := new(profileResponse)
	r.Profile.Username = u.Username
	r.Profile.Bio = u.Bio
	r.Profile.Image = u.Image
	r.Profile.Following, _ = us.IsFollower(u.ID, playerID)
	return r
}

type itemResponse struct {
	Slug           string    `json:"slug"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Body           string    `json:"body"`
	TagList        []string  `json:"tagList"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	Favorited      bool      `json:"favorited"`
	FavoritesCount int       `json:"favoritesCount"`
	Author         struct {
		Username  string  `json:"username"`
		Bio       *string `json:"bio"`
		Image     *string `json:"image"`
		Following bool    `json:"following"`
	} `json:"author"`
}

type singleItemResponse struct {
	Item *itemResponse `json:"item"`
}

type itemListResponse struct {
	Items      []*itemResponse `json:"items"`
	ItemsCount int             `json:"itemsCount"`
}

func newItemResponse(c echo.Context, a *model.Item) *singleItemResponse {
	ar := new(itemResponse)
	ar.TagList = make([]string, 0)
	ar.Slug = a.Slug
	ar.Title = a.Title
	ar.Description = a.Description
	ar.Body = a.Body
	ar.CreatedAt = a.CreatedAt
	ar.UpdatedAt = a.UpdatedAt
	for _, t := range a.Tags {
		ar.TagList = append(ar.TagList, t.Tag)
	}
	for _, u := range a.Favorites {
		if u.ID == playerIDFromToken(c) {
			ar.Favorited = true
		}
	}
	ar.FavoritesCount = len(a.Favorites)
	ar.Author.Username = a.Author.Username
	ar.Author.Image = a.Author.Image
	ar.Author.Bio = a.Author.Bio
	ar.Author.Following = a.Author.FollowedBy(playerIDFromToken(c))
	return &singleItemResponse{ar}
}

func newItemListResponse(us player.Store, playerID uint, items []model.Item, count int) *itemListResponse {
	r := new(itemListResponse)
	r.Items = make([]*itemResponse, 0)
	for _, a := range items {
		ar := new(itemResponse)
		ar.TagList = make([]string, 0)
		ar.Slug = a.Slug
		ar.Title = a.Title
		ar.Description = a.Description
		ar.Body = a.Body
		ar.CreatedAt = a.CreatedAt
		ar.UpdatedAt = a.UpdatedAt
		for _, t := range a.Tags {
			ar.TagList = append(ar.TagList, t.Tag)
		}
		for _, u := range a.Favorites {
			if u.ID == playerID {
				ar.Favorited = true
			}
		}
		ar.FavoritesCount = len(a.Favorites)
		ar.Author.Username = a.Author.Username
		ar.Author.Image = a.Author.Image
		ar.Author.Bio = a.Author.Bio
		ar.Author.Following, _ = us.IsFollower(a.AuthorID, playerID)
		r.Items = append(r.Items, ar)
	}
	r.ItemsCount = count
	return r
}

type commentResponse struct {
	ID        uint      `json:"id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Author    struct {
		Username  string  `json:"username"`
		Bio       *string `json:"bio"`
		Image     *string `json:"image"`
		Following bool    `json:"following"`
	} `json:"author"`
}

type singleCommentResponse struct {
	Comment *commentResponse `json:"comment"`
}

type commentListResponse struct {
	Comments []commentResponse `json:"comments"`
}

func newCommentResponse(c echo.Context, cm *model.Comment) *singleCommentResponse {
	comment := new(commentResponse)
	comment.ID = cm.ID
	comment.Body = cm.Body
	comment.CreatedAt = cm.CreatedAt
	comment.UpdatedAt = cm.UpdatedAt
	comment.Author.Username = cm.Player.Username
	comment.Author.Image = cm.Player.Image
	comment.Author.Bio = cm.Player.Bio
	comment.Author.Following = cm.Player.FollowedBy(playerIDFromToken(c))
	return &singleCommentResponse{comment}
}

func newCommentListResponse(c echo.Context, comments []model.Comment) *commentListResponse {
	r := new(commentListResponse)
	cr := commentResponse{}
	r.Comments = make([]commentResponse, 0)
	for _, i := range comments {
		cr.ID = i.ID
		cr.Body = i.Body
		cr.CreatedAt = i.CreatedAt
		cr.UpdatedAt = i.UpdatedAt
		cr.Author.Username = i.Player.Username
		cr.Author.Image = i.Player.Image
		cr.Author.Bio = i.Player.Bio
		cr.Author.Following = i.Player.FollowedBy(playerIDFromToken(c))

		r.Comments = append(r.Comments, cr)
	}
	return r
}

type tagListResponse struct {
	Tags []string `json:"tags"`
}

func newTagListResponse(tags []model.Tag) *tagListResponse {
	r := new(tagListResponse)
	for _, t := range tags {
		r.Tags = append(r.Tags, t.Tag)
	}
	return r
}
