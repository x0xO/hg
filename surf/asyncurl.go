package surf

import "net/http"

// AsyncURL struct represents an asynchronous URL with additional information
// such as data, context, setHeaders, addHeaders, and url.
type AsyncURL struct {
	context    any
	addHeaders map[string]string
	setHeaders map[string]string
	url        string
	addCookies []http.Cookie
	data       []any
}

// asyncResponse struct represents an asynchronous response that embeds a
// *Response pointer and a context field of type any.
type asyncResponse struct {
	*Response
	context any
}

// asyncRequest struct represents an asynchronous request that embeds a
// *Request pointer and additional fields: context, setHeaders, and addHeaders,
// all of type map[string]string.
type asyncRequest struct {
	*Request
	context    any
	setHeaders map[string]string
	addHeaders map[string]string
	addCookies []http.Cookie
}

// NewAsyncURL creates a new AsyncURL object with the provided URL string and returns a pointer to
// it.
func NewAsyncURL(url string) *AsyncURL { return &AsyncURL{url: url} }

// Context sets the context of the AsyncURL object and returns a pointer to the updated object.
func (au *AsyncURL) Context(context any) *AsyncURL {
	au.context = context
	return au
}

// Data sets the data of the AsyncURL object and returns a pointer to the updated object.
func (au *AsyncURL) Data(data ...any) *AsyncURL {
	au.data = data
	return au
}

// SetHeaders sets the headers of the AsyncURL object and returns a pointer to the updated object.
func (au *AsyncURL) SetHeaders(headers map[string]string) *AsyncURL {
	au.setHeaders = headers
	return au
}

// AddHeaders adds headers to the AsyncURL object and returns a pointer to the updated object.
func (au *AsyncURL) AddHeaders(headers map[string]string) *AsyncURL {
	au.addHeaders = headers
	return au
}

// AddCookies adds cookies to the AsyncURL object and returns a pointer to the updated object.
func (au *AsyncURL) AddCookies(cookies ...http.Cookie) *AsyncURL {
	au.addCookies = cookies
	return au
}

// Context returns the context of the asyncResponse object.
func (ar asyncResponse) Context() any { return ar.context }
