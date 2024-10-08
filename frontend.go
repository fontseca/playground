package playground

import (
  "html/template"
  "log/slog"
  "net/http"
  "net/url"
)

// Scanner scans an incoming HTTP request, parses it, sends it to the backend,
// and writes the formatted response to the HTTP response writer.
func Scanner(w http.ResponseWriter, r *http.Request) {
  req := parse(r)
  response := backend(r.Context(), req)

  w.Header().Set("Content-Type", "text/plain; charset=utf-8")
  w.WriteHeader(http.StatusOK)
  template.HTMLEscape(w, response.Bytes())
}

// Renderer renders the website template and writes it to the HTTP response writer.
func Renderer(w http.ResponseWriter, r *http.Request) {
  website().Render(r.Context(), w)
}

// request represents an HTTP request with a method, target URL, and headers.
type request struct {
  // method is the HTTP method used for the request.
  method string

  // target is the URL to which the request is sent.
  target *url.URL

  // header contains the HTTP headers to be included in the request.
  header http.Header
}

// parse extracts the HTTP method and target URL from an incoming HTTP request
// and returns a new request struct containing these values.
func parse(r *http.Request) *request {
  var err error
  req := new(request)

  if err = r.ParseForm(); nil != err {
    slog.Error(err.Error())
  }

  target := r.PostFormValue("request_target")
  method := r.PostFormValue("request_method")
  headerKeys := r.PostForm["header-key"]
  headerValues := r.PostForm["header-value"]

  req.header = http.Header{}
  for n := range len(headerKeys) {
    key, value := headerKeys[n], headerValues[n] // must be trimmed already
    if "" == key && "" == value {
      continue
    }

    req.header.Add(http.CanonicalHeaderKey(key), value)
  }

  req.method = method

  req.target, err = url.Parse(target)
  if nil != err {
    slog.Error(err.Error())
    req.target = &url.URL{}
  }

  return req
}
