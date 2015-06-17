package mocks

import (
	"github.com/stretchr/testify/mock"

	"github.com/hevnly/eevy/event"
	"github.com/hevnly/eevy/handler/config"
)

type Handler struct {
	mock.Mock

	Name string
}

func (m *Handler) Exec(evt event.Event) {
	m.Called(evt)
}
func (m *Handler) GetType() string {
	ret := m.Called()

	r0 := ret.Get(0).(string)

	return r0
}
func (m *Handler) GetName() string {
	return m.Name
}
func (m *Handler) SetName(s string) {
	m.Called(s)
}
func (m *Handler) GetMessage() string {
	ret := m.Called()

	r0 := ret.Get(0).(string)

	return r0
}
func (m *Handler) GetConfig() config.Handler {
	ret := m.Called()

	r0 := ret.Get(0).(config.Handler)

	return r0
}
func (m *Handler) String() string {
	ret := m.Called()

	r0 := ret.Get(0).(string)

	return r0
}
func (m *Handler) Init(s string) {
	m.Called(s)
}
