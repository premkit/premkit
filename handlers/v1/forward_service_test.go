package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStripLeadingSlashIfPresent(t *testing.T) {
	before := "/test"
	after := stripLeadingSlashIfPresent(before)
	assert.Equal(t, "test", after)

	before = "test"
	after = stripLeadingSlashIfPresent(before)
	assert.Equal(t, "test", after)
}

func TestIsPathPrefix(t *testing.T) {
	servicePath := "service"
	requestPath := "/service/something"
	ok := isPathPrefix(servicePath, requestPath)
	assert.True(t, ok)

	servicePath = "/service"
	requestPath = "service/something"
	ok = isPathPrefix(servicePath, requestPath)
	assert.True(t, ok)

	servicePath = "/service"
	requestPath = "svc/something"
	ok = isPathPrefix(servicePath, requestPath)
	assert.False(t, ok)
}

func TestCreateForwardPath(t *testing.T) {
	servicePath := "service"
	requestPath := "/service/something"
	forwardPath := createForwardPath(servicePath, requestPath)
	assert.Equal(t, "/something", forwardPath)

	servicePath = "service/test/something"
	requestPath = "/service/test/something/one/two"
	forwardPath = createForwardPath(servicePath, requestPath)
	assert.Equal(t, "/one/two", forwardPath)
}

func TestCreateForwardPathWithQuery(t *testing.T) {
	servicePath := "service"
	requestPath := "/service/something?a=b"
	forwardPath := createForwardPath(servicePath, requestPath)
	assert.Equal(t, "/something?a=b", forwardPath)

	servicePath = "service/test/something"
	requestPath = "/service/test/something/one/two?a=b"
	forwardPath = createForwardPath(servicePath, requestPath)
	assert.Equal(t, "/one/two?a=b", forwardPath)
}
