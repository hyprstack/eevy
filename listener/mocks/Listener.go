package mocks

import "github.com/stretchr/testify/mock"

import "github.com/hevnly/eevy/event"
import "github.com/hevnly/eevy/listener/config"

type Listener struct {
	mock.Mock
}

func (m *Listener) Exec(evt event.Event) {
	m.Called(evt)
}
func (m *Listener) GetType() string {
	ret := m.Called()

	r0 := ret.Get(0).(string)

	return r0
}
func (m *Listener) GetConfig() config.Listener {
	ret := m.Called()

	r0 := ret.Get(0).(config.Listener)

	return r0
}
