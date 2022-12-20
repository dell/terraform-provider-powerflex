package helper

import (
	"bytes"
	"encoding/json"
)

// PrettyJSON indents and encode the data
func PrettyJSON(data interface{}) string {
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "\t")

	err := encoder.Encode(data)
	if err != nil {
		return ""
	}
	return buffer.String()
}
