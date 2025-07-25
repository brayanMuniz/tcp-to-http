package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeadersParse(t *testing.T) {

	// Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Valid single header, extra white space
	headers = NewHeaders()
	data = []byte("  Host: localhost:42069  \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 27, n)
	assert.False(t, done)

	// Valid 2 headers, extra white space
	headers = NewHeaders()
	data = []byte("  Host: localhost:42069\r\n Anime: forSure  \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, 25, n)
	assert.Equal(t, "localhost:42069", headers["host"])

	n, done, err = headers.Parse(data[n:]) // continue to where it left off
	assert.Equal(t, 19, n)
	assert.Equal(t, "forSure", headers["anime"])
	assert.False(t, done)

	// Valid same key multiple values
	headers = NewHeaders()
	data = []byte("Anime: onimai\r\n Anime: friren\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, 15, n)
	assert.Equal(t, "onimai", headers["anime"])

	n, done, err = headers.Parse(data[n:]) // continue to where it left off
	assert.Equal(t, 16, n)
	assert.Equal(t, "onimai, friren", headers["anime"])
	assert.False(t, done)

	// Invalid key
	headers = NewHeaders()
	data = []byte("H@st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}
