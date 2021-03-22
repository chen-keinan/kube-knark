package khttp

type HTTPRequestData struct {
	Method        string
	Host          string
	RequestURI    string
	SourceIP      string
	DestinationIP string
	Headers       map[string]string
	ReqBody       string
	StartTime     string
	Proto         string
}

type HTTPResponseData struct {
	SourceIP      string
	DestinationIP string
	Headers       []string
	ResBody       string
	StatusCode    int
	StartTime     string
}

type HTTPNetData struct {
	HttpRequestData  *HTTPRequestData
	HttpResponseData *HTTPResponseData
}

func NewHTTPNetData(HttpRequestData *HTTPRequestData, HttpResponseData *HTTPResponseData) *HTTPNetData {
	return &HTTPNetData{HttpRequestData: HttpRequestData, HttpResponseData: HttpResponseData}
}
