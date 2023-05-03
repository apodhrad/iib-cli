package format

import "encoding/json"

const INDENT string = "  "

func Json(v any, indent bool) (string, error) {
	var out []byte
	var err error

	if indent {
		out, err = json.MarshalIndent(v, "", INDENT)
	} else {
		out, err = json.Marshal(v)
	}
	return string(out), err
}
