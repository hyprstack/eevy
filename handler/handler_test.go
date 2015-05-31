package handler

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hevnly/eevy/event"
)

func TestStub(t *testing.T) {
	assert.True(t, true, "This is good. Canary test passing")
}

func TestGetMagicString(t *testing.T) {

	ms := findMagicStrings("${test}(){{{${test2.test3.test-4}..vwvvw")
	assert.Equal(t, 2, len(ms), "Ensure we get the correct amount of strings")

	assert.Equal(t, "${test}", ms[0], "Check that the first magic string is correct")
	assert.Equal(t, "${test2.test3.test-4}", ms[1], "Check that the second magic string is correct")
}

func TestMagicString(t *testing.T) {

	msg := "{\"testKey\":\"testValue\"}"
	evt := event.Event{
		Event: "test.event",
	}
	json.Unmarshal([]byte(msg), &evt.Message)

	{
		r := magicString("${message.testKey}", evt)
		assert.Equal(t, "testValue", r, "Ensure message is changed")
	}
	{
		r := magicString("${message.testKey}-test.test1", evt)
		assert.Equal(t, "testValue-test.test1", r, "Ensure only variable is changed")
	}
	{
		r := magicString("${message}", evt)
		assert.Equal(t, msg, r, "Ensure entire message is copied")
	}
}
