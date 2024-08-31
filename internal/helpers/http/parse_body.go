package httphelpers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ParseJSON(req *http.Request, v any) error {
	if req.Body == nil {
		return fmt.Errorf("missing request body")
	}

	return json.NewDecoder(req.Body).Decode(v)
}
