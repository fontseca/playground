package playground

import (
  "bytes"
  "fmt"
  "maps"
  "net/http"
  "slices"
  "sort"
  "strings"
  "time"
)

// responseBuilder is used to construct an HTTP response message with custom start lines, headers, and body content.
type responseBuilder struct {
  startLine []byte
  header    http.Header
  body      bytes.Buffer
}

func newResponseBuilder() *responseBuilder {
  return &responseBuilder{
    header: http.Header{},
  }
}

// SetStartLine sets the start line of the HTTP response with the provided protocol version and status.
func (r *responseBuilder) SetStartLine(proto, status string) {
  r.startLine = []byte(fmt.Sprintf("%s %s", strings.TrimSpace(proto), strings.TrimSpace(status)))
}

// SetHeaders copies the provided headers into the HTTP response headers.
func (r *responseBuilder) SetHeaders(h http.Header) {
  maps.Copy(r.header, h)
}

// DefaultHeaders sets some default headers for the HTTP response, discarding any previous headers.
func (r *responseBuilder) DefaultHeaders() {
  r.header = http.Header{}
  r.header.Set("Content-Length", fmt.Sprintf("%d", r.body.Len()))
  r.header.Set("Content-Type", "text/plain; charset=utf-8")
  r.header.Set("Date", time.Now().Format(time.RFC1123))
  r.header.Set("Server", "fontseca.dev/playground (v1.0)")
}

// Write appends the provided byte slice to the body of the HTTP response.
func (r *responseBuilder) Write(p []byte) (n int, err error) {
  return r.body.Write(p)
}

// WriteError writes an error message to the HTTP response, discarding any previous written bytes to the body.
func (r *responseBuilder) WriteError(err error) {
  r.body.Reset()
  r.body.WriteString(err.Error())
}

func (r *responseBuilder) build() *bytes.Buffer {
  buffer := &bytes.Buffer{}
  buffer.Grow(len(r.startLine) + 1 + len(r.header) + len(r.body.Bytes())) // approximate growth

  if len(r.startLine) == 0 {
    r.SetStartLine("HTTP/1.0", "200 OK")
  }

  buffer.Write(r.startLine)
  buffer.WriteRune('\n')

  keys := make([]string, 0, len(r.header))

  for key := range maps.Keys(r.header) {
    keys = append(keys, key)
  }

  sort.Strings(keys)

  for key := range slices.Values(keys) {
    buffer.WriteString(fmt.Sprintf("%s: %s\n", key, r.header.Get(key)))
  }

  buffer.WriteRune('\n')

  buffer.Write(r.body.Bytes())
  return buffer
}

// Bytes returns the HTTP response as a byte slice.
func (r *responseBuilder) Bytes() []byte {
  return r.build().Bytes()
}

// String returns the HTTP response as a string.
func (r *responseBuilder) String() string {
  return r.build().String()
}
