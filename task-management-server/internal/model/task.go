package model

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type Task struct {
	Id   int       `json:"id"`
	Text string    `json:"text" validate:"required,min=1"`
	Tags []string  `json:"tags" validate:"dive,required"`
	Due  time.Time `json:"due" validate:"required"`
}

func (t *Task) Validate() error {
	validate := validator.New()
	return validate.Struct(t)
}
