package format

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const DATA_CSV string = `COMPANY,NAME,LOCATION
Cool Company,Joe,Boston
Cool Company,Kate,New York
`

func TestFormatCsv(t *testing.T) {
	out, err := Csv(DATA)
	assert.Nil(t, err)
	assert.Equal(t, DATA_CSV, out)
}
