package event

import (
	"testing"

	"encoding/json"

	"github.com/stretchr/testify/assert"
)

func TestStub(t *testing.T) {
	assert.True(t, true, "This is good. Canary test passing")
}

func TestUnmarshal(t *testing.T) {

	var evt Event
	msg := "{\"event\": \"test.cli\", \"message\": {\"test\": {\"test1.1\": 23}, \"test1\": \"2\"}}"
	err := json.Unmarshal([]byte(msg), &evt)
	assert.Equal(t, nil, err, "Check there were no errors")
}

func TestGet(t *testing.T) {
	var evt Event
	msg := "{\"event\": \"test.cli\", \"message\": {\"easy\": 123, \"test\": {\"test1-1\": \"result\"}, \"test1\": \"2\"}}"
	err := json.Unmarshal([]byte(msg), &evt)
	assert.Equal(t, nil, err, "Check there were no errors")

	val := evt.Get("easy")
	wf, ok := val.(float64)
	assert.True(t, ok, "Check that returned value is a number")
	assert.Equal(t, 123.0, wf, "Check there were no errors")

	val = evt.Get("test.test1-1")
	ws, ok := val.(string)
	assert.True(t, ok, "Check that returned value is a string")
	assert.Equal(t, "result", ws, "Check there were no errors")
}

func TestGetString(t *testing.T) {

	var evt Event
	msg := "{\"event\": \"test.cli\", \"message\": {\"easy\": 123, \"test\": {\"test1-1\": \"result\"}, \"test1\": \"2\"}}"
	err := json.Unmarshal([]byte(msg), &evt)
	assert.Equal(t, nil, err, "Check there were no errors")

	val := evt.GetString("test.test1-1")
	assert.Equal(t, "result", val, "Check we recieved the correct string")
}

func TestGetInt(t *testing.T) {

	var evt Event
	msg := "{\"event\": \"test.cli\", \"message\": {\"easy\": 123, \"test\": {\"test1-1\": \"result\"}, \"test1\": \"2\"}}"
	err := json.Unmarshal([]byte(msg), &evt)
	assert.Equal(t, nil, err, "Check there were no errors")

	val := evt.GetInt("easy")
	assert.Equal(t, 123, val, "Check we recieved the correct int")
}

func TestGetFloat(t *testing.T) {

	var evt Event
	msg := "{\"event\": \"test.cli\", \"message\": {\"easy\": 123, \"test\": {\"test1-1\": \"result\"}, \"test1\": \"2\"}}"
	err := json.Unmarshal([]byte(msg), &evt)
	assert.Equal(t, nil, err, "Check there were no errors")

	val := evt.GetFloat("easy")
	assert.Equal(t, 123.0, val, "Check we recieved the correct int")
}
