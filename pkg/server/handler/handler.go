package handler

import (
	"20dojo-online/pkg/server/model"
)

type Handler struct {
	model model.ModelInterface
}

func NewHandler(m model.ModelInterface) *Handler {
	return &Handler{
		model: m,
	}
}
