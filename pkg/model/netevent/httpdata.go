package netevent

//HTTPRequestData http request data for validation
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

//HTTPResponseData http response data for validation
type HTTPResponseData struct {
	SourceIP      string
	DestinationIP string
	Headers       []string
	ResBody       string
	StatusCode    int
	StartTime     string
}

//HTTPNetData encapsulate request and response
type HTTPNetData struct {
	HTTPRequestData  *HTTPRequestData
	HTTPResponseData *HTTPResponseData
}

//NewHTTPNetData return new NewHTTPNetData object
func NewHTTPNetData(HTTPRequestData *HTTPRequestData, HTTPResponseData *HTTPResponseData) *HTTPNetData {
	return &HTTPNetData{HTTPRequestData: HTTPRequestData, HTTPResponseData: HTTPResponseData}
}
