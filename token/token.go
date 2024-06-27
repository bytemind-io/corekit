/*
Copyright 2024 The corego Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package token

import (
	"sync"

	"github.com/pkoukk/tiktoken-go"
)

// TikToken represents the tiktoken interface.
type TikToken interface {
	// Encode encodes the tiktoken.
	Encode() (*tiktoken.Tiktoken, error)
}

// Service represents the tiktoken service.
type Service struct {
	encoder map[string]TikToken
	mu      sync.RWMutex
}

// Open opens the tiktoken service.
func (s *Service) Open() error {
	s.Register("gpt-3.5-turbo", &gpt3turbo{})
	s.Register("gpt-4", &gpt4{})
	s.Register("gpt-4o", &gpt4o{})
	return nil
}

// Register registers the model with the tiktoken service.
func (s *Service) Register(model string, encoder TikToken) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.encoder[model] = encoder
}

// Deregister the model from the tiktoken service.
func (s *Service) Deregister(model string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.encoder, model)
}

// List returns the list of models.
func (s *Service) List() (map[string]*tiktoken.Tiktoken, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	diags := make(map[string]*tiktoken.Tiktoken, len(s.encoder))
	for k, v := range s.encoder {
		sv, err := v.Encode()
		if err != nil {
			continue
		}
		diags[k] = sv
	}
	return diags, nil
}
