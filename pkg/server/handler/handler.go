package handler

import (
	"20dojo-online/pkg/server/model"
)

type Handler struct {
	model *model.Model
}

func NewHandler(m *model.Model) *Handler {
	return &Handler{
		model: m,
	}
}
