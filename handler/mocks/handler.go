package mocks

import "github.com/stretchr/testify/mock"

import "github.com/hevnly/eevy/event"
import "github.com/hevnly/eevy/handler/config"

type Handler struct {
	mock.Mock
}

func (m *Handler) Exec(evt event.Event) {
	m.Called(evt)
}
func (m *Handler) GetType() string {
	ret := m.Called()

	r0 := ret.Get(0).(string)

	return r0
}
func (m *Handler) GetConfig() config.Handler {
	ret := m.Called()

	r0 := ret.Get(0).(config.Handler)

	return r0
}
