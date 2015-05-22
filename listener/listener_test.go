package listener

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"hevnly/eevy/event"
	"hevnly/eevy/listener/mocks"
)

func TestStub(t *testing.T) {
	assert.True(t, true, "This is good. Canary test passing")
}

func TestAdd(t *testing.T) {

	rootList := EventListener{}
	rootList.Name = ""

	rootList.Add("_", new(mocks.Listener))
	assert.Equal(t, len(rootList.Subset["*"].Listeners), 1, "Adding '_'")

	rootList.Add("*", new(mocks.Listener))
	assert.Equal(t, len(rootList.Subset["*"].Listeners), 2, "Adding '*'")

	rootList.Add("test", new(mocks.Listener))
	assert.Equal(t, len(rootList.Subset["*"].Listeners), 2, "Adding 'test', check root wildcard")
	assert.Equal(t, len(rootList.Subset["test"].Listeners), 1, "Adding 'test', check root subset")

	rootList.Add("test.*", new(mocks.Listener))
	assert.Equal(t, len(rootList.Subset["*"].Listeners), 2, "Adding 'test.*', check root wildcard")
	assert.Equal(t, len(rootList.Subset["test"].Listeners), 1, "Adding 'test.*', check root subset")
	assert.Equal(t, len(rootList.Subset["test"].Subset["*"].Listeners), 1, "Adding 'test.*', check subset wildcard")

	rootList.Add("test.sub", new(mocks.Listener))
	assert.Equal(t, len(rootList.Subset["*"].Listeners), 2, "Adding 'test.sub', check root wildcard")
	assert.Equal(t, len(rootList.Subset["test"].Listeners), 1, "Adding 'test.sub', check root subset")
	assert.Equal(t, len(rootList.Subset["test"].Subset["*"].Listeners), 1, "Adding 'test.sub', check subset wildcard")
	assert.Equal(t, len(rootList.Subset["test"].Subset["sub"].Listeners), 1, "Adding 'test.sub', check subset sub")

	rootList.Add("test1.*", new(mocks.Listener))
	assert.Equal(t, len(rootList.Subset["*"].Listeners), 2, "Adding 'test1.*', check root wildcard")
	assert.Equal(t, len(rootList.Subset["test"].Listeners), 1, "Adding 'test1.*', check root subset")
	assert.Equal(t, len(rootList.Subset["test"].Subset["*"].Listeners), 1, "Adding 'test1.*', check subset wildcard")
	assert.Equal(t, len(rootList.Subset["test"].Subset["sub"].Listeners), 1, "Adding 'test1.*', check subset sub")
	assert.Equal(t, len(rootList.Subset["test"].Subset["sub"].Listeners), 1, "Adding 'test1.*', check test1 wildcard")
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

func TestGetMagicString(t *testing.T) {

	lb := ListenerBase{}
	ms := lb.findMagicStrings("${test}(){{{${test2.test3.test-4}..vwvvw")
	assert.Equal(t, 2, len(ms), "Ensure we get the correct amount of strings")

	assert.Equal(t, "${test}", ms[0], "Check that the first magic string is correct")
	assert.Equal(t, "${test2.test3.test-4}", ms[1], "Check that the second magic string is correct")
}

func TestMagicString(t *testing.T) {

	lb := ListenerBase{}
	evt := event.Event{
		Event:   "test.event",
		Message: "{\"testKey\":\"testValue\"}",
	}

	r1 := lb.magicString("${message.testKey}", evt)
	assert.Equal(t, "testValue", r1, "Ensure message is changed")

	r2 := lb.magicString("${message.testKey}-test.test1", evt)
	assert.Equal(t, "testValue-test.test1", r2, "Ensure only variable is changed")

	r3 := lb.magicString("${message}", evt)
	assert.Equal(t, "{\"testKey\":\"testValue\"}", r3, "Ensure entire message is copied")
}

func TestGetMessage(t *testing.T) {

	lb1 := ListenerBase{}
	evt1 := event.Event{
		Event:   "test.event",
		Message: "testMessage1",
	}
	r1 := lb1.GetMessage(evt1)
	assert.Equal(t, "testMessage1", r1, "Empty message body")

	lb2 := ListenerBase{
		Message: "${message}",
	}
	evt2 := event.Event{
		Event:   "test.event",
		Message: "testMessage2",
	}
	r2 := lb2.GetMessage(evt2)
	assert.Equal(t, "testMessage2", r2, "Asses '${message}'")

	lb3 := ListenerBase{
		Message: "{${message}}.blah-blah!-{}",
	}
	evt3 := event.Event{
		Event:   "test.event",
		Message: "testMessage3",
	}
	r3 := lb3.GetMessage(evt3)
	assert.Equal(t, "{testMessage3}.blah-blah!-{}", r3, "Asses json message body")
}

func createRootList() (*EventListener, *map[string]*mocks.Listener) {

	rootList := EventListener{}
	rootList.Name = ""

	l := make(map[string]*mocks.Listener)
	l["_"] = new(mocks.Listener)
	l["*"] = new(mocks.Listener)
	l["test"] = new(mocks.Listener)
	l["test.*"] = new(mocks.Listener)
	l["test.sub"] = new(mocks.Listener)
	l["test.sub."] = new(mocks.Listener)
	l["test1.*"] = new(mocks.Listener)
	l["test1"] = new(mocks.Listener)

	for name, list := range l {
		rootList.Add(name, list)
	}
	return &rootList, &l
}
