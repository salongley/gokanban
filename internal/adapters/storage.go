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

func (s *InMemoryService) Update(id int, title string) models.Card {
	return models.Card{ID: id, Title: title}
}
func (s *InMemoryService) Find(id int) models.Card { return models.Card{ID: id, Title: "Found"} }
