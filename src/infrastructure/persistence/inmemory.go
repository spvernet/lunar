package persistence

import (
	"encoding/json"
	"lunar/src/domain"
	"sort"
	"sync"
	"time"
)

type MemoryStore struct {
	mu      sync.RWMutex
	rockets map[string]*domain.Rocket
	seen    map[string]map[int]struct{}
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		rockets: make(map[string]*domain.Rocket),
		seen:    make(map[string]map[int]struct{}),
	}
}

// Apply implementa idempotencia + last-write-wins como comentamos.
func (s *MemoryStore) Apply(env domain.MessageEnvelope) error {
	ch, num, kind, raw := env.Metadata.Channel, env.Metadata.MessageNum, env.Metadata.MessageType, env.Message

	s.mu.Lock()
	defer s.mu.Unlock()

	// idempotencia
	if s.seen[ch] == nil {
		s.seen[ch] = make(map[int]struct{})
	}
	if _, dup := s.seen[ch][num]; dup {
		return nil
	}
	s.seen[ch][num] = struct{}{}

	r := s.ensureRocket(ch)

	// Para eventos no conmutativos: ignora si es más antiguo que el último aplicado
	isNonCommutative := kind == domain.TypeLaunched || kind == domain.TypeMissionChanged || kind == domain.TypeExploded
	if isNonCommutative && num < r.LastMsgNum {
		return nil // no “deshacemos” estado
	}

	switch kind {
	case domain.TypeLaunched:
		var p domain.RocketLaunchedPayload
		if err := json.Unmarshal(raw, &p); err != nil {
			return err
		}
		r.Type = p.Type
		r.Mission = p.Mission
		if r.Speed < p.LaunchSpeed {
			r.Speed = p.LaunchSpeed
		} // no reducimos velocidad

	case domain.TypeSpeedIncreased:
		var p domain.RocketSpeedDeltaPayload
		if err := json.Unmarshal(raw, &p); err != nil {
			return err
		}
		r.Speed += p.By

	case domain.TypeSpeedDecreased:
		var p domain.RocketSpeedDeltaPayload
		if err := json.Unmarshal(raw, &p); err != nil {
			return err
		}
		r.Speed -= p.By

	case domain.TypeMissionChanged:
		var p domain.RocketMissionChangedPayload
		if err := json.Unmarshal(raw, &p); err != nil {
			return err
		}
		r.Mission = p.NewMission

	case domain.TypeExploded:
		r.Status = domain.StatusExploded
	}

	if num > r.LastMsgNum {
		r.LastMsgNum = num
	}
	r.UpdatedAt = time.Now()
	return nil
}

func (s *MemoryStore) Get(channel string) (domain.Rocket, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.rockets[channel]
	if !ok {
		return domain.Rocket{}, false, nil
	}
	return *r, true, nil
}

func (s *MemoryStore) List(sortBy, order string) ([]domain.Rocket, error) {
	s.mu.RLock()
	items := make([]domain.Rocket, 0, len(s.rockets))
	for _, r := range s.rockets {
		items = append(items, *r)
	}
	s.mu.RUnlock()

	less := func(i, j int) bool {
		switch sortBy {
		case "speed":
			if order == "desc" {
				return items[i].Speed > items[j].Speed
			}
			return items[i].Speed < items[j].Speed
		case "updated_at":
			if order == "desc" {
				return items[i].UpdatedAt.After(items[j].UpdatedAt)
			}
			return items[i].UpdatedAt.Before(items[j].UpdatedAt)
		default: // channel
			if order == "desc" {
				return items[i].Channel > items[j].Channel
			}
			return items[i].Channel < items[j].Channel
		}
	}
	sort.Slice(items, less)
	return items, nil
}

func (s *MemoryStore) ensureRocket(ch string) *domain.Rocket {
	r, ok := s.rockets[ch]
	if !ok {
		r = &domain.Rocket{Channel: ch, Status: domain.StatusActive}
		s.rockets[ch] = r
	}
	return r
}
