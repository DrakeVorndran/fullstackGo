package handler

import (
	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
	"golang-starter-pack/model"
)

type playerUpdateRequest struct {
	Player struct {
		Username string `json:"username"`
		Email    string `json:"email" validate:"email"`
		Password string `json:"password"`
		Bio      string `json:"bio"`
		Image    string `json:"image"`
	} `json:"player"`
}

func newPlayerUpdateRequest() *playerUpdateRequest {
	return new(playerUpdateRequest)
}

func (r *playerUpdateRequest) populate(u *model.Player) {
	r.Player.Username = u.Username
	r.Player.Email = u.Email
	r.Player.Password = u.Password
	if u.Bio != nil {
		r.Player.Bio = *u.Bio
	}
	if u.Image != nil {
		r.Player.Image = *u.Image
	}
}

func (r *playerUpdateRequest) bind(c echo.Context, u *model.Player) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	u.Username = r.Player.Username
	u.Email = r.Player.Email
	if r.Player.Password != u.Password {
		h, err := u.HashPassword(r.Player.Password)
		if err != nil {
			return err
		}
		u.Password = h
	}
	u.Bio = &r.Player.Bio
	u.Image = &r.Player.Image
	return nil
}

type playerRegisterRequest struct {
	Player struct {
		Username string `json:"username" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	} `json:"player"`
}

func (r *playerRegisterRequest) bind(c echo.Context, u *model.Player) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	u.Username = r.Player.Username
	u.Email = r.Player.Email
	h, err := u.HashPassword(r.Player.Password)
	if err != nil {
		return err
	}
	u.Password = h
	return nil
}

type playerLoginRequest struct {
	Player struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	} `json:"player"`
}

func (r *playerLoginRequest) bind(c echo.Context) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	return nil
}

type itemCreateRequest struct {
	Items struct {
		Title       string   `json:"title" validate:"required"`
		Description string   `json:"description" validate:"required"`
		Body        string   `json:"body" validate:"required"`
		Tags        []string `json:"tagList, omitempty"`
	} `json:"item"`
}

func (r *itemCreateRequest) bind(c echo.Context, a *model.Item) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	a.Title = r.Items.Title
	a.Slug = slug.Make(r.Items.Title)
	a.Description = r.Items.Description
	a.Body = r.Items.Body
	if r.Items.Tags != nil {
		for _, t := range r.Items.Tags {
			a.Tags = append(a.Tags, model.Tag{Tag: t})
		}
	}
	return nil
}

type itemUpdateRequest struct {
	Items struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Body        string   `json:"body"`
		Tags        []string `json:"tagList"`
	} `json:"item"`
}

func (r *itemUpdateRequest) populate(a *model.Item) {
	r.Items.Title = a.Title
	r.Items.Description = a.Description
	r.Items.Body = a.Body
}

func (r *itemUpdateRequest) bind(c echo.Context, a *model.Item) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	a.Title = r.Items.Title
	a.Slug = slug.Make(a.Title)
	a.Description = r.Items.Description
	a.Body = r.Items.Body
	return nil
}

type createCommentRequest struct {
	Comment struct {
		Body string `json:"body" validate:"required"`
	} `json:"comment"`
}

func (r *createCommentRequest) bind(c echo.Context, cm *model.Comment) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	cm.Body = r.Comment.Body
	cm.PlayerID = playerIDFromToken(c)
	return nil
}
