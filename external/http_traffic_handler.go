//nolint
package external

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"github.com/chen-keinan/kube-knark/pkg/model/netevent"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/chen-keinan/kube-knark/external/httpport"
	"github.com/hsiafan/glow/iox/filex"

	"bufio"

	"github.com/google/gopacket/tcpassembly/tcpreader"
)

// ConnectionKey contains src and dst endpoint identify a connection
type ConnectionKey struct {
	src Endpoint
	dst Endpoint
}

func (ck *ConnectionKey) reverse() ConnectionKey {
	return ConnectionKey{ck.dst, ck.src}
}

// return the src ip and port
func (ck *ConnectionKey) srcString() string {
	return ck.src.String()
}

// return the dst ip and port
func (ck *ConnectionKey) dstString() string {
	return ck.dst.String()
}

// HTTPConnectionHandler impl ConnectionHandler
type HTTPConnectionHandler struct {
	Option  *Option
	Printer *Printer
}

func (handler *HTTPConnectionHandler) handle(src Endpoint, dst Endpoint, connection *TCPConnection) {
	ck := ConnectionKey{src, dst}
	trafficHandler := &HTTPTrafficHandler{
		key:       ck,
		buffer:    new(bytes.Buffer),
		option:    handler.Option,
		printer:   handler.Printer,
		startTime: connection.lastTimestamp,
	}
	WaitGroup.Add(1)
	go trafficHandler.handle(connection)
}

func (handler *HTTPConnectionHandler) finish() {
	//handler.Printer.Finish()
}

// HTTPTrafficHandler parse a http connection traffic and send to Printer
type HTTPTrafficHandler struct {
	startTime time.Time
	endTime   time.Time
	key       ConnectionKey
	buffer    *bytes.Buffer
	option    *Option
	printer   *Printer
}

// read http request/response stream, and do output
func (h *HTTPTrafficHandler) handle(connection *TCPConnection) {
	defer WaitGroup.Done()
	defer connection.upStream.Close()
	defer connection.downStream.Close()
	// filter by args setting

	requestReader := bufio.NewReader(connection.upStream)
	defer discardAll(requestReader)
	responseReader := bufio.NewReader(connection.downStream)
	defer discardAll(responseReader)
	var hReqdata *netevent.HTTPRequestData
	var hResdata *netevent.HTTPResponseData
	for {
		h.buffer = new(bytes.Buffer)
		hReqdata = &netevent.HTTPRequestData{}
		hResdata = &netevent.HTTPResponseData{}
		filtered := false
		req, err := httpport.ReadRequest(requestReader)
		h.startTime = connection.lastTimestamp

		if err != nil {
			if err != io.EOF {
				fmt.Fprintln(os.Stderr, "Error parsing HTTP requests:", err)
			}
			break
		}

		if h.option.Host != "" && !wildcardMatch(req.Host, h.option.Host) {
			filtered = true
		}
		if h.option.Uri != "" && !wildcardMatch(req.RequestURI, h.option.Uri) {
			filtered = true
		}
		hReqdata.Host = req.Host
		expectContinue := req.Header.Get("Expect") == "100-continue"

		resp, err := httpport.ReadResponse(responseReader, nil)

		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				break
			} else {
				fmt.Fprintln(os.Stderr, "Error parsing HTTP response:", err, connection.clientID)
			}
			if !filtered {
				h.printRequest(req, hReqdata)
				h.writeLine("")
				h.printer.Send(netevent.NewHTTPNetData(hReqdata, hResdata))
			} else {
				discardAll(req.Body)
			}
			break
		}

		if h.option.StatusSet != nil && !h.option.StatusSet.Contains(resp.StatusCode) {
			filtered = true
		}

		if !filtered {
			h.printRequest(req, hReqdata)
			h.writeLine("")
			h.endTime = connection.lastTimestamp
			h.printResponse(resp, hResdata)
			h.printer.Send(netevent.NewHTTPNetData(hReqdata, hResdata))
		} else {
			discardAll(req.Body)
			discardAll(resp.Body)

		}

		if expectContinue {
			if resp.StatusCode == 100 {
				// read next response, the real response
				resp, err := httpport.ReadResponse(responseReader, nil)
				if err == io.EOF {
					fmt.Fprintln(os.Stderr, "Error parsing HTTP requests: unexpected end, ", err)
					break
				}
				if err == io.ErrUnexpectedEOF {
					fmt.Fprintln(os.Stderr, "Error parsing HTTP requests: unexpected end, ", err)
					// here return directly too, to avoid error when long polling connection is used
					break
				}
				if err != nil {
					fmt.Fprintln(os.Stderr, "Error parsing HTTP response:", err, connection.clientID)
					break
				}
				if !filtered {
					h.printResponse(resp, hResdata)
					h.printer.Send(netevent.NewHTTPNetData(hReqdata, hResdata))
				} else {
					discardAll(resp.Body)
				}
			} else if resp.StatusCode == 417 {

			}
		}
	}

	h.printer.Send(netevent.NewHTTPNetData(hReqdata, hResdata))
}

func (h *HTTPTrafficHandler) handleWebsocket(requestReader *bufio.Reader, responseReader *bufio.Reader) {
	//TODO: websocket

}

func (h *HTTPTrafficHandler) writeLineFormat(format string, a ...interface{}) {
	fmt.Fprintf(h.buffer, format, a...)
}

func (h *HTTPTrafficHandler) write(a ...interface{}) {
	fmt.Fprint(h.buffer, a...)
}

func (h *HTTPTrafficHandler) writeLine(a ...interface{}) {
	fmt.Fprintln(h.buffer, a...)
}

func (h *HTTPTrafficHandler) printRequestMark() {
	h.writeLine()
}

func (h *HTTPTrafficHandler) printHeader(header httpport.Header, hrd *netevent.HTTPRequestData) {
	hrd.Headers = make(map[string]string)
	for name, values := range header {
		for _, value := range values {
			hrd.Headers[name] = value
		}
	}
}

// print http request
func (h *HTTPTrafficHandler) printRequest(req *httpport.Request, hrd *netevent.HTTPRequestData) {
	defer discardAll(req.Body)
	h.printNormalRequest(req, hrd)
}

var blockHeaders = map[string]bool{
	"Content-Length":    true,
	"Transfer-Encoding": true,
	"Connection":        true,
	"Accept-Encoding:":  true,
}

// print http request curl command
func (h *HTTPTrafficHandler) printCurlRequest(req *httpport.Request) {
	//TODO: expect-100 continue handle

	h.writeLine()
	h.writeLine(strings.Repeat("*", 10), " REQUEST ", h.key.srcString(), " -----> ", h.key.dstString(), " // ", h.startTime.Format(time.RFC3339Nano))
	h.writeLineFormat("curl -X %v http://%v%v \\\n", req.Method, h.key.dstString(), req.RequestURI)
	var reader io.ReadCloser
	var deCompressed bool
	if h.option.DumpBody {
		reader = req.Body
		deCompressed = false
	} else {
		reader, deCompressed = h.tryDecompress(req.Header, req.Body)
	}

	if deCompressed {
		defer reader.Close()
	}
	seq := 0
	for name, values := range req.Header {
		seq++
		if blockHeaders[name] {
			continue
		}
		if deCompressed {
			if name == "Content-Encoding" {
				continue
			}
		}
		for idx, value := range values {
			if seq == len(req.Header) && idx == len(values)-1 {
				h.writeLineFormat("    -H '%v: %v'\n", name, value)
			} else {
				h.writeLineFormat("    -H '%v: %v' \\\n", name, value)
			}
		}
	}

	if req.ContentLength == 0 || req.Method == "GET" || req.Method == "HEAD" || req.Method == "TRACE" ||
		req.Method == "OPTIONS" {
		h.writeLine()
		return
	}

	if h.option.DumpBody {
		filename := "request-" + uriToFileName(req.RequestURI, h.startTime)
		h.writeLineFormat("    -d '@%v'", filename)

		err := filex.WriteAllFromReader(filename, reader)
		if err != nil {
			h.writeLine("dump to file failed:", err)
		}
	} else {
		br := bufio.NewReader(reader)
		// optimize for one line body
		firstLine, err := br.ReadString('\n')
		if err != nil && err != io.EOF {
			// read error
		} else if err == io.EOF && !strings.Contains(firstLine, "'") {
			h.writeLineFormat("    -d '%v'", strconv.Quote(firstLine))
		} else {
			h.writeLineFormat("    -d @- << HTTP_DUMP_BODY_EOF\n")
			h.write(firstLine)
			for {
				line, err := br.ReadString('\n')
				if err != nil && err != io.EOF {
					break
				}
				h.write(line)
				if err == io.EOF {
					h.writeLine("\nHTTP_DUMP_BODY_EOF")
					break
				}
			}
		}
	}

	h.writeLine()
}

// print http request
func (h *HTTPTrafficHandler) printNormalRequest(req *httpport.Request, hrd *netevent.HTTPRequestData) {
	if h.option.Level == "url" {
		hrd.Method = req.Method
		hrd.Host = req.Host
		hrd.RequestURI = req.RequestURI
		return
	}
	hrd.SourceIP = h.key.srcString()
	hrd.DestinationIP = h.key.dstString()
	hrd.StartTime = h.startTime.Format(time.RFC3339Nano)

	hrd.Method = req.Method
	hrd.Proto = req.Proto
	hrd.RequestURI = req.RequestURI
	h.printHeader(req.Header, hrd)

	if !(req.ContentLength == 0 || req.Method == "GET" || req.Method == "HEAD" || req.Method == "TRACE" ||
		req.Method == "OPTIONS") {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			hrd.ReqBody = string(body)
		}
	}
}

// print http response
func (h *HTTPTrafficHandler) printResponse(resp *httpport.Response, hrd *netevent.HTTPResponseData) {
	defer discardAll(resp.Body)
	if h.option.Level == "url" {
		return
	}
	hrd.SourceIP = h.key.srcString()
	hrd.DestinationIP = h.key.dstString()
	hrd.StartTime = h.startTime.Format(time.RFC3339Nano)
	hrd.StatusCode = resp.StatusCode
	hrd.Headers = resp.RawHeaders
	if !(resp.ContentLength == 0 || resp.StatusCode == 304 || resp.StatusCode == 204) {
		body, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			hrd.ResBody = string(body)
		}
	}
}

func (h *HTTPTrafficHandler) tryDecompress(header httpport.Header, reader io.ReadCloser) (io.ReadCloser, bool) {
	contentEncoding := header.Get("Content-Encoding")
	var nr io.ReadCloser
	var err error
	if contentEncoding == "" {
		// do nothing
		return reader, false
	} else if strings.Contains(contentEncoding, "gzip") {
		nr, err = gzip.NewReader(reader)
		if err != nil {
			return reader, false
		}
		return nr, true
	} else if strings.Contains(contentEncoding, "deflate") {
		nr, err = zlib.NewReader(reader)
		if err != nil {
			return reader, false
		}
		return nr, true
	} else {
		return reader, false
	}
}

// print http request/response body
func (h *HTTPTrafficHandler) printBody(header httpport.Header, reader io.ReadCloser) {

	// deal with content encoding such as gzip, deflate
	nr, decompressed := h.tryDecompress(header, reader)
	if decompressed {
		defer nr.Close()
	}

	// check mime type and charset
	contentType := header.Get("Content-Type")
	if contentType == "" {
		// TODO: detect content type using httpport.DetectContentType()
	}
	mimeTypeStr, charset := parseContentType(contentType)
	var mimeType = parseMimeType(mimeTypeStr)
	isText := mimeType.isTextContent()
	isBinary := mimeType.isBinaryContent()

	if !isText {
		err := h.printNonTextTypeBody(nr, contentType, isBinary)
		if err != nil {
			h.writeLine("{Read content error", err, "}")
		}
		return
	}

	var body string
	var err error
	if charset == "" {
		// response do not set charset, try to detect
		var data []byte
		data, err := ioutil.ReadAll(nr)
		if err == nil {
			// TODO: try to detect charset
			body = string(data)
		}
	} else {
		body, err = readToStringWithCharset(nr, charset)
	}
	if err != nil {
		h.writeLine("{Read body failed", err, "}")
		return
	}

	// prettify json
	if mimeType.subType == "json" || likeJSON(body) {
		var jsonValue interface{}
		_ = json.Unmarshal([]byte(body), &jsonValue)
		prettyJSON, err := json.MarshalIndent(jsonValue, "", "    ")
		if err == nil {
			body = string(prettyJSON)
		}
	}
	h.writeLine(body)
	h.writeLine()
}

func (h *HTTPTrafficHandler) printNonTextTypeBody(reader io.Reader, contentType string, isBinary bool) error {
	if h.option.Force && !isBinary {
		data, err := ioutil.ReadAll(reader)
		if err != nil {
			return err
		}
		// TODO: try to detect charset
		str := string(data)
		h.writeLine(str)
		h.writeLine()
	} else {
		h.writeLine("{Non-text body, content-type:", contentType, ", len:", discardAll(reader), "}")
	}
	return nil
}

func discardAll(r io.Reader) (dicarded int) {
	return tcpreader.DiscardBytesToEOF(r)
}

func uriToFileName(uri string, t time.Time) string {
	timeStr := t.Format("2006_01_02_15_04_05.000000")
	filename := strings.ReplaceAll(uri, "/", "_") + "-" + timeStr
	return filename[1:]
}
