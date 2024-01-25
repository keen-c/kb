package controllers

import "kb/kb/models"

type Game struct {
	Answer bool            `json:"answer"`
	Game   models.Question `json:"game,omitempty"`
}
