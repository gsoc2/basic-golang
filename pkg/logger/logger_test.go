package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	defer func() {
		output = os.Stdout
	}()
	r, w, err := os.Pipe()
	require.NoError(t, err)
	output = w

	read := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, err := io.Copy(&buf, r)
		require.NoError(t, err)
		read <- buf.String()
	}()

	log := New().
		ID("id123").
		Data(Data{"1": "1", "2": "2"}).
		Root(Data{"4": "4"})

	log.Err(fmt.Errorf("normal error")).Error("foo", Data{"1": "11", "3": "3"})
	log.Err(errors.New("pkg error")).Error("foo", Data{"1": "11", "3": "3"})
	log.Warn("foo", Data{"1": "11", "3": "3"})
	log.Info("foo", Data{"1": "11", "3": "3"})
	log.Debug("foo", Data{"1": "11", "3": "3"})

	err = w.Close()
	require.NoError(t, err)

	line := <-read

	assert.Contains(t, line, `"timestamp":`)
	assert.Contains(t, line, `"hostname":`)
	assert.Contains(t, line, `"id":"id123"`)
	assert.Contains(t, line, `"1":"11"`)
	assert.Contains(t, line, `"2":"2"`)
	assert.Contains(t, line, `"3":"3"`)
	assert.Contains(t, line, `"4":"4"`)
	assert.Contains(t, line, `"error":"normal error",`)
	assert.Contains(t, line, `"error":"pkg error",`)
	assert.Contains(t, line, `"level":"info"`)
	assert.Contains(t, line, `"message":"foo"`)
}
