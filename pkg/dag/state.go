package dag

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

type State struct {
	mu    sync.RWMutex
	state map[string]any
}

func NewState(initialState map[string]any) *State {
	return &State{
		state: initialState,
	}
}

func (s *State) Keys() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	keys := make([]string, 0, len(s.state))
	for k := range s.state {
		keys = append(keys, k)
	}
	return keys
}

func (s *State) Set(key string, value interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.state[key] = value
}

func (s *State) Get(key string) interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.state[key]
	if !ok {
		return nil
	}
	return val
}

func (s *State) MarshalJSON() ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return json.Marshal(s.state)
}

func (s *State) UnmarshalJSON(data []byte) error {
	var tempMap map[string]interface{}
	if err := json.Unmarshal(data, &tempMap); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.state = tempMap
	return nil
}

func (s *State) Filtered() *State {
	out := NewState(map[string]any{})

	keys := s.Keys()
	for _, key := range keys {
		if strings.HasPrefix(key, "_") {
			continue
		}

		value := s.Get(key)
		if value != nil {
			out.Set(key, value)
		}
	}

	return out
}

func GetStateKey[T any](ctx context.Context, key string) (T, error) {
	value := GetState(ctx).Get(key)
	if value == nil {
		return *new(T), nil
	}
	typed, ok := value.(T)
	if !ok {
		return *new(T), fmt.Errorf("expected type %T, got %T", new(T), value)
	}
	return typed, nil
}

func SetStateKey[T any](ctx context.Context, key string, value T) {
	GetState(ctx).Set(key, value)
}
