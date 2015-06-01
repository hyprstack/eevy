package listener

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/hevnly/eevy/event"
	"github.com/hevnly/eevy/handler/mocks"
)

func TestStub(t *testing.T) {
	assert.True(t, true, "This is good. Canary test passing")
}

func TestAdd(t *testing.T) {

	rootList := Listener{}
	rootList.Name = ""

	rootList.Add("_", &mocks.Handler{})
	assert.Equal(t, len(rootList.Subset["*"].Handlers), 1, "Adding '_'")

	rootList.Add("*", &mocks.Handler{})
	assert.Equal(t, len(rootList.Subset["*"].Handlers), 2, "Adding '*'")

	rootList.Add("test", &mocks.Handler{})
	assert.Equal(t, len(rootList.Subset["*"].Handlers), 2, "Adding 'test', check root wildcard")
	assert.Equal(t, len(rootList.Subset["test"].Handlers), 1, "Adding 'test', check root subset")

	rootList.Add("test.*", &mocks.Handler{})
	assert.Equal(t, len(rootList.Subset["*"].Handlers), 2, "Adding 'test.*', check root wildcard")
	assert.Equal(t, len(rootList.Subset["test"].Handlers), 1, "Adding 'test.*', check root subset")
	assert.Equal(t, len(rootList.Subset["test"].Subset["*"].Handlers), 1, "Adding 'test.*', check subset wildcard")

	rootList.Add("test.sub", &mocks.Handler{})
	assert.Equal(t, len(rootList.Subset["*"].Handlers), 2, "Adding 'test.sub', check root wildcard")
	assert.Equal(t, len(rootList.Subset["test"].Handlers), 1, "Adding 'test.sub', check root subset")
	assert.Equal(t, len(rootList.Subset["test"].Subset["*"].Handlers), 1, "Adding 'test.sub', check subset wildcard")
	assert.Equal(t, len(rootList.Subset["test"].Subset["sub"].Handlers), 1, "Adding 'test.sub', check subset sub")

	rootList.Add("test1.*", &mocks.Handler{})
	assert.Equal(t, len(rootList.Subset["*"].Handlers), 2, "Adding 'test1.*', check root wildcard")
	assert.Equal(t, len(rootList.Subset["test"].Handlers), 1, "Adding 'test1.*', check root subset")
	assert.Equal(t, len(rootList.Subset["test"].Subset["*"].Handlers), 1, "Adding 'test1.*', check subset wildcard")
	assert.Equal(t, len(rootList.Subset["test"].Subset["sub"].Handlers), 1, "Adding 'test1.*', check subset sub")
	assert.Equal(t, len(rootList.Subset["test"].Subset["sub"].Handlers), 1, "Adding 'test1.*', check test1 wildcard")
}

func TestExec(t *testing.T) {

	rootList, l := createRootList()

	for _, list := range *l {
		list.On("IsRepeater").Return(false)
		list.On("Exec", mock.AnythingOfType("event.Event")).Return()
	}

	events := [9]event.Event{
		event.Event{Event: "blah.blah"},
		event.Event{Event: "blah"},
		event.Event{Event: "test"},
		event.Event{Event: "test.blah"},
		event.Event{Event: "test.sub"},
		event.Event{Event: "test.sub."},
		event.Event{Event: "test1"},
		event.Event{Event: "test1.yeah"},
		event.Event{Event: "test1.yeah.crap"},
	}

	for _, event := range events {
		rootList.Exec(event)
	}

	for evt, list := range *l {
		switch evt {
		case "*":
			list.AssertNumberOfCalls(t, "Exec", 9)
		case "_":
			list.AssertNumberOfCalls(t, "Exec", 9)
		case "test":
			list.AssertNumberOfCalls(t, "Exec", 1)
		case "test.*":
			list.AssertNumberOfCalls(t, "Exec", 4)
		case "test.sub":
			list.AssertNumberOfCalls(t, "Exec", 1)
		case "test1":
			list.AssertNumberOfCalls(t, "Exec", 1)
		case "test1.*":
			list.AssertNumberOfCalls(t, "Exec", 3)
		}
	}
}

func createRootList() (*Listener, *map[string]*mocks.Handler) {

	rootList := Listener{}
	rootList.Name = ""

	l := make(map[string]*mocks.Handler)
	l["_"] = &mocks.Handler{}
	l["*"] = &mocks.Handler{}
	l["test"] = &mocks.Handler{}
	l["test.*"] = &mocks.Handler{}
	l["test.sub"] = &mocks.Handler{}
	l["test.sub."] = &mocks.Handler{}
	l["test1.*"] = &mocks.Handler{}
	l["test1"] = &mocks.Handler{}

	for name, list := range l {
		rootList.Add(name, list)
	}
	return &rootList, &l
}
