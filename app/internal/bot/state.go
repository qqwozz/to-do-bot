package bot

import "sync"

type convStep int

const (
	stepIdle convStep = iota
	stepTitle
	stepDesc
	stepDate
	stepTime
)

type conversation struct {
	step  convStep
	plan  PlanRequest
	store *sync.Map
}

// conversationStore хранит состояние диалогов по chatID.
type conversationStore struct {
	mu sync.Map
}

func newConversationStore() *conversationStore {
	return &conversationStore{}
}

func (s *conversationStore) get(chatID int64) *conversation {
	v, ok := s.mu.Load(chatID)
	if !ok {
		return &conversation{step: stepIdle, store: &sync.Map{}}
	}
	return v.(*conversation)
}

func (s *conversationStore) set(chatID int64, c *conversation) {
	s.mu.Store(chatID, c)
}

func (s *conversationStore) reset(chatID int64) {
	s.mu.Delete(chatID)
}
