package si

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

type ContextKey string

func (c ContextKey) String() string {
	return string(c)
}

type Context struct {
	Request  *http.Request
	Response http.ResponseWriter
}

// SetAttribute sets a key-value pair in the context
func (ctx *Context) SetAttribute(key ContextKey, value interface{}) {
	ctx.Request = ctx.Request.WithContext(
		context.WithValue(ctx.Request.Context(), key, value),
	)
}

// GetAttribute gets a value from the context
func (ctx *Context) GetAttribute(key string) interface{} {
	return ctx.Request.Context().Value(key)
}

// IP gets the IP address of the client
func (ctx *Context) IP() string {
	if ctx.Request.Header.Get("X-Forwarded-For") != "" {
		return ctx.Request.Header.Get("X-Forwarded-For")
	}

	return ctx.Request.RemoteAddr
}

// -----
// Header methods
// -----

// Accepts checks if the client accepts a certain type of content
func (ctx *Context) Accepts(offers ...string) bool {
	acceptHeader := ctx.HeaderString("Accept")

	acceptArray := strings.Split(acceptHeader, ",")

	for _, offer := range offers {
		for _, accept := range acceptArray {
			accept = strings.TrimSpace(accept)

			if strings.Contains(accept, offer) {
				return true
			}

			if strings.Contains(accept, "*/*") {
				return true
			}

			if strings.Contains(accept, "*") {
				return true
			}
		}
	}

	return false
}

// HeaderString gets a header value as a string
func (ctx *Context) HeaderString(key string) string {
	value := ctx.Request.Header.Get(key)
	if value == "" {
		return ""
	}

	return value
}

// CookieString gets a cookie value as a string
func (ctx *Context) CookieString(key string) string {
	cookie, err := ctx.Request.Cookie(key)
	if err != nil {
		return ""
	}

	return cookie.Value
}

// -----
// URL methods
// -----

// Method gets the HTTP method
func (ctx *Context) Method() string {
	return ctx.Request.Method
}

// Host gets the host
func (ctx *Context) Host() string {
	return ctx.Request.Host
}

// FullUrl gets the full URL
func (ctx *Context) FullUrl() string {
	return ctx.Request.URL.String()
}

// Path gets the path
func (ctx *Context) Path() string {
	return ctx.Request.URL.Path
}

// PathAndQuery gets the path and query
func (ctx *Context) PathAndQuery() string {
	if ctx.Request.URL.RawQuery == "" {
		return ctx.Request.URL.Path
	}

	return ctx.Request.URL.Path + "?" + ctx.Request.URL.RawQuery
}

// Query gets the query parameters
func (ctx *Context) Query() map[string]any {
	query := map[string]any{}

	for key, value := range ctx.Request.URL.Query() {
		query[key] = value
	}

	return query
}

// QueryString gets a query parameter as a string
func (ctx *Context) QueryString(key string) string {
	value := ctx.Request.URL.Query().Get(key)
	if value == "" {
		return ""
	}

	return value
}

// QueryInt gets a query parameter as an int
func (ctx *Context) QueryInt(key string) int {
	value, err := strconv.Atoi(ctx.Request.URL.Query().Get(key))
	if err != nil {
		return 0
	}

	return value
}

// QueryBool gets a query parameter as a bool
func (ctx *Context) QueryBool(key string) bool {
	value, err := strconv.ParseBool(ctx.Request.URL.Query().Get(key))
	if err != nil {
		return false
	}

	return value
}

// ParamString gets a URL parameter as a string
func (ctx *Context) ParamString(key string) string {
	value := chi.URLParam(ctx.Request, key)
	if value == "" {
		return ""
	}

	return value
}

// ParamInt gets a URL parameter as an int
func (ctx *Context) ParamInt(key string) int {
	value, err := strconv.Atoi(chi.URLParam(ctx.Request, key))
	if err != nil {
		return 0
	}

	return value
}

// ParamBool gets a URL parameter as a bool
func (ctx *Context) ParamBool(key string) bool {
	value, err := strconv.ParseBool(chi.URLParam(ctx.Request, key))
	if err != nil {
		return false
	}

	return value
}

// -----
// Body methods
// -----

// GetFormData gets the form data
func (ctx *Context) GetFormData() (map[string][]string, error) {
	req := ctx.Request

	err := req.ParseForm()
	if err != nil {
		return nil, err
	}

	err = req.ParseMultipartForm(32 << 20)
	if err != nil {
		return nil, err
	}

	if len(req.Form) == 0 {
		return req.PostForm, nil
	}

	return req.Form, nil
}

// GetRawContent gets the raw content
func (ctx *Context) GetRawContent() ([]byte, error) {
	req := ctx.Request

	buf := &bytes.Buffer{}
	_, _ = io.Copy(buf, req.Body)
	err := req.Body.Close()
	if err != nil {
		return nil, err
	}

	b := buf.Bytes()
	// reset the body so it can be read again
	req.Body = io.NopCloser(bytes.NewBuffer(b))

	return b, nil
}

// UnmarshalJSONBody unmarshals the JSON body
func (ctx *Context) UnmarshalJSONBody(v interface{}) error {
	body, err := ctx.GetRawContent()
	if err != nil {
		return err
	}

	return json.Unmarshal(body, v)
}

// -----
// Response methods
// -----

// WriteHeader writes a header
func (ctx *Context) WriteHeader(key string, value string) {
	ctx.Response.Header().Add(key, value)
}

// WriteHeaders writes multiple headers
func (ctx *Context) WriteHeaders(headers map[string]string) {
	for key, value := range headers {
		ctx.Response.Header().Add(key, value)
	}
}

// WriteStatus writes a status code
// should be called after all headers are written
func (ctx *Context) WriteStatus(statusCode int) {
	ctx.Response.WriteHeader(statusCode)
}

// SetCookie writes a cookie
// Can be used to clear a cookie by setting MaxAge to -1
func (ctx *Context) SetCookie(cookie *http.Cookie) {
	ctx.Response.Header().Add("Set-Cookie", cookie.String())
}

// SendBytes sends a byte array
func (ctx *Context) SendBytes(data []byte, statusCode int) {
	if statusCode == 0 {
		statusCode = http.StatusOK
	}

	ctx.Response.WriteHeader(statusCode)
	_, _ = ctx.Response.Write(data)
}

// SendString sends a string
func (ctx *Context) SendString(data string, statusCode int) {
	if statusCode == 0 {
		statusCode = http.StatusOK
	}

	ctx.Response.WriteHeader(statusCode)
	_, _ = fmt.Fprint(ctx.Response, data)
}

// SendJSON sends a JSON response
func (ctx *Context) SendJSON(data any, statusCode int) {
	if statusCode == 0 {
		statusCode = http.StatusOK
	}

	ctx.WriteHeader("Content-Type", "application/json")
	ctx.WriteStatus(statusCode)
	if data != nil {
		err := json.NewEncoder(ctx.Response).Encode(data)
		if err != nil {
			_, _ = fmt.Fprintf(ctx.Response, `
			{
				"error": {
					"code": 500,
					"message": "Internal Server Error. Could not encode JSON."
				}
			}`)
		}
	}
}

// SendErrorJSON sends an error JSON response
func (ctx *Context) SendErrorJSON(message string, statusCode int) {
	if statusCode == 0 {
		statusCode = http.StatusInternalServerError
	}

	ctx.SendJSON(map[string]interface{}{
		"error": map[string]interface{}{
			"code":    statusCode,
			"message": message,
		},
	}, statusCode)
}

// SendStream sends a stream
func (ctx *Context) SendStream(stream io.ReadCloser, statusCode int) {
	if statusCode == 0 {
		statusCode = http.StatusOK
	}

	ctx.WriteStatus(statusCode)
	_, _ = io.Copy(ctx.Response, stream)
}

// -----
// Response headers methods
// -----

// Redirect redirects to a URL
func (ctx *Context) Redirect(url string, statusCode int) {
	if statusCode == 0 {
		statusCode = http.StatusFound
	}

	http.Redirect(ctx.Response, ctx.Request, url, statusCode)
}

// WriteEarlyHintScript writes a preload hint for a script
func (ctx *Context) WriteEarlyHintScript(path string) {
	ctx.WriteHeader("Link", fmt.Sprintf("</%s>; rel=preload; as=script", path))
	ctx.WriteStatus(103)
}

// WriteEarlyHintStyle writes a preload hint for a style
func (ctx *Context) WriteEarlyHintStyle(path string) {
	ctx.WriteHeader("Link", fmt.Sprintf("</%s>; rel=preload; as=style", path))
	ctx.WriteStatus(103)
}

// NotFound sends a 404 response
func (ctx *Context) NotFound() {
	ctx.Response.WriteHeader(404)
	ctx.Response.Header().Add("X-Error-Code", "404")
}

// -----
// Shortcuts
// -----

// SB is a shortcut for SendBytes
func (ctx *Context) SB(data []byte) {
	ctx.SendBytes(data, 0)
}

// SS is a shortcut for SendString
func (ctx *Context) SS(data string) {
	ctx.SendString(data, 0)
}

// SJ is a shortcut for SendJSON
func (ctx *Context) SJ(data any) {
	ctx.SendJSON(data, 0)
}
