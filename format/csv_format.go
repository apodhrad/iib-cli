package format

import (
	"bytes"
	"errors"
	"strings"
)

const CSV_DELIMETER string = ","

func Csv(data [][]string) (string, error) {
	if len(data) == 0 || len(data[0]) == 0 {
		return "", errors.New("Missing a line with headers!")
	}

	var out bytes.Buffer
	for i := 0; i < len(data); i++ {
		out.WriteString(strings.Join(data[i], CSV_DELIMETER))
		out.WriteString("\n")
	}

	return out.String(), nil
}
