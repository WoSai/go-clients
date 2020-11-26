package dingtalk

import (
	"encoding/json"
	"net/http"

	"github.com/jacexh/requests"
)

func UnmarshalAndParseError(v Response) requests.Interceptor {
	return func(request *http.Request, response *http.Response, bytes []byte) error {
		if err := json.Unmarshal(bytes, v); err != nil {
			return err
		}
		return v.GotErr()
	}
}
