package resp

import "time"

type ResourceVo struct {
	ID            int          `json:"id"`
	Name          string       `json:"name"`
	Url           string       `json:"url"`
	RequestMethod string       `json:"request_method"`
	IsAnonymous   int          `json:"is_anonymous"`
	CreatedAt     time.Time    `json:"created_at"`
	Children      []ResourceVo `json:"children"`
}
