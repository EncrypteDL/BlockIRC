package pkg

import (
	"bytes"
	"sync"
)

type SaslState struct {
	sync.RWMutex

	started bool

	buffer *bytes.Buffer
	mech   string

	authcid string
}

func NewSaslState() *SaslState {
	return &SaslState{buffer: &bytes.Buffer{}}
}

func (s *SaslState) Reset() {
	s.Lock()
	defer s.Unlock()

	s.started = false
	s.buffer.Reset()
	s.mech = ""
	s.authcid = ""
}

func (s *SaslState) Started() bool {
	s.RLock()
	defer s.RUnlock()

	return s.started
}

func (s *SaslState) Start() {
	s.Lock()
	defer s.Unlock()

	s.started = true
}

func (s *SaslState) WriteString(data string) {
	s.Lock()
	defer s.Unlock()

	s.buffer.WriteString(data)
}

func (s SaslState) Len() int {
	s.RLock()
	defer s.RUnlock()

	return s.buffer.Len()
}

func (s *SaslState) String() string {
	s.RLock()
	defer s.RUnlock()

	return s.buffer.String()
}

func (s *SaslState) Login(authcid string) {
	s.Lock()
	defer s.Unlock()

	s.started = false
	s.buffer.Reset()
	s.mech = ""

	s.authcid = authcid
}

func (s *SaslState) Id() string {
	s.RLock()
	defer s.RUnlock()

	return s.authcid
}
