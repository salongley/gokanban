package adapters

import (
	"gokanban/internal/models"
	"sync"
	"time"
)

// --- B. In-Memory Repository (Database Adapter) ---
type InMemoryService struct {
	sync.Mutex
	cards  []models.Card
	nextID int
}

func NewInMemoryService() *InMemoryService {
	return &InMemoryService{
		cards:  []models.Card{},
		nextID: 1,
	}
}

func (s *InMemoryService) List() []models.Card {
	return s.cards
}

func (s *InMemoryService) Create(title, column string) models.Card {
	s.Lock()
	defer s.Unlock()
	card := models.Card{
		ID:       s.nextID,
		Title:    title,
		Column:   column,
		EditedAt: time.Now().Format(time.Kitchen),
	}
	s.cards = append(s.cards, card)
	s.nextID++
	return card
}

func (s *InMemoryService) Update(id int, card models.Card) models.Card {
	s.Lock()
	defer s.Unlock()
	return models.Card{ID: id, Title: card.Title, Column: card.Column}
}
func (s *InMemoryService) Find(id int) models.Card {
	for _, card := range s.cards {
		if card.ID == id {
			return card
		}
	}

	return models.Card{}
}
