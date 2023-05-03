package format

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/rodaine/table"
)

func stringsToInterfaces(strings []string) []interface{} {
	interfaceSlice := make([]interface{}, len(strings))
	for i, v := range strings {
		interfaceSlice[i] = v
	}
	return interfaceSlice
}

func newTable(headers ...string) table.Table {
	table.DefaultHeaderFormatter = func(format string, vals ...interface{}) string {
		return strings.ToUpper(fmt.Sprintf(format, vals...))
	}
	return table.New(stringsToInterfaces(headers)...)
}

func tableToString(tbl table.Table) string {
	var tblBuf bytes.Buffer
	tbl.WithWriter(&tblBuf)
	tbl.Print()
	return tblBuf.String()
}

func Table(data [][]string) (string, error) {
	if len(data) == 0 || len(data[0]) == 0 {
		return "", errors.New("Missing a line with headers!")
	}
	// headers
	table := newTable(data[0]...)

	// items
	for i := 1; i < len(data); i++ {
		table.AddRow(stringsToInterfaces(data[i])...)
	}

	out := tableToString(table)
	return out, nil
}
