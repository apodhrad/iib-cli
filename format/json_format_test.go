package format

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Employee struct {
	Name     string `json:"name"`
	Location string `json:"location"`
}

type Company struct {
	Name      string     `json:"name"`
	Employees []Employee `json:"employees"`
}

var JOE = Employee{
	Name:     "Joe",
	Location: "Boston",
}

var KATE = Employee{
	Name:     "Kate",
	Location: "New York",
}

var COOL_COMPANY = Company{
	Name:      "Cool Company",
	Employees: []Employee{JOE, KATE},
}

var COOL_COMPANY_JSON = `{"name":"Cool Company","employees":[{"name":"Joe","location":"Boston"},{"name":"Kate","location":"New York"}]}`

var COOL_COMPANY_JSON_INDENT = `{
  "name": "Cool Company",
  "employees": [
    {
      "name": "Joe",
      "location": "Boston"
    },
    {
      "name": "Kate",
      "location": "New York"
    }
  ]
}`

func TestJson(t *testing.T) {
	out, err := Json(COOL_COMPANY, false)
	assert.Nil(t, err)
	assert.Equal(t, COOL_COMPANY_JSON, out)
}

func TestJsonIndent(t *testing.T) {
	out, err := Json(COOL_COMPANY, true)
	assert.Nil(t, err)
	assert.Equal(t, COOL_COMPANY_JSON_INDENT, out)
}
