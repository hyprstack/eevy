// Used to mock a listener in unit testing.
package mocks

import "github.com/stretchr/testify/mock"

import "github.com/hevnly/eevy/event"

type Listener struct {
	mock.Mock
}

func (m *Listener) Exec(evt event.Event) {
	m.Called(evt)
}
func (m *Listener) GetMessage(evt event.Event) string {
	ret := m.Called(evt)

	r0 := ret.Get(0).(string)

	return r0
}
