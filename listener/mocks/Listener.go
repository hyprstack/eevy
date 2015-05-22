package mocks

import "github.com/stretchr/testify/mock"

import "hevnly/eevy/event"
import "hevnly/eevy/config"

type Listener struct {
	mock.Mock
}

func (m *Listener) Exec(evt event.Event) {
	m.Called(evt)
}
func (m *Listener) Init(conf config.Listener) {
	m.Called(conf)
}
func (m *Listener) GetMessage(evt event.Event) string {
	ret := m.Called(evt)

	r0 := ret.Get(0).(string)

	return r0
}
