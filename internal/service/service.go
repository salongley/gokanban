package service

import "gokanban/internal/models"

// Service Interface (Port)
type CardService interface {
	List() []models.Card
	Create(title, column string) models.Card
	Update(id int, title string) models.Card
	Find(id int) models.Card
}
