package fileinput_test

import (
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/doloopwhile/go-fileinput"
)

func TestLines(t *testing.T) {
	assert := assert.New(t)
	data := map[string]string{
		"a": "a01\na02\na03",
		"b": "b01\nb02\nb03\n",
		"c": "c01\nc02\nc03",
	}
	sc := &fileinput.Scanner{
		Args: []string{"a", "b", "c"},
		Open: func(name string) (io.ReadCloser, error) {
			return ioutil.NopCloser(strings.NewReader(data[name])), nil
		},
	}
	lines := []string{}
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	assert.NoError(sc.Err())
	expected := []string{
		"a01", "a02", "a03",
		"b01", "b02", "b03",
		"c01", "c02", "c03",
	}
	assert.Equal(expected, lines)
}
