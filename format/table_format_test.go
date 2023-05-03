package format

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var DATA = [][]string{
	{"COMPANY", "NAME", "LOCATION"},
	{"Cool Company", "Joe", "Boston"},
	{"Cool Company", "Kate", "New York"},
}

var DATA_TABLE = `COMPANY       NAME  LOCATION  
Cool Company  Joe   Boston    
Cool Company  Kate  New York  
`

func TestFormatTable(t *testing.T) {
	out, err := Table(DATA)
	assert.Nil(t, err)
	assert.Equal(t, DATA_TABLE, out)
}

func TestFormatEmptyTable(t *testing.T) {
	out, err := Table([][]string{})
	assert.NotNil(t, err)
	assert.Equal(t, "", out)
}
