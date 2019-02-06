package methods

type ErrMsq struct {
	Error string
}

type Query struct {
	Method  string             `json:"method"`
	Answers map[string]Answers `json:"answers"`
}

type Answers struct {
	HttpStatus      int                 `json:"http_status"`
	MimeType        string              `json:"mime_type"`
	WaitTime        string              `json:"wait_time,omitempty"`
	Queries         map[string]string   `json:"query_arguments,omitempty"`
	RequestHeaders  map[string]string   `json:"request_headers,omitempty"`
	ResponseHeaders map[string]string   `json:"response_headers,omitempty"`
	PostForm        map[string]string   `json:"post_arguments,omitempty"`
	Data            string              `json:"data"`
	DataFromFile    string              `json:"data_from_file"`
	Cookies         map[string][]Cookie `json:"cookies,omitempty"`
}

type Cookie struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	MaxAge   int    `json:"max_age,omitempty"`
	Path     string `json:"path"`
	Domain   string `json:"domain,omitempty"`
	Secure   bool   `json:"secure,omitempty"`
	HttpOnly bool   `json:"http_only,omitempty"`
}

type ResponseArgs struct {
	Code     int
	MimeType string
	Data     []byte
}
