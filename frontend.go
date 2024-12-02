package playground

import (
  "context"
  "errors"
  "fmt"
  "html/template"
  "io"
  "log/slog"
  "mime/multipart"
  "net/http"
  "net/url"
  "os"
  "path"
  "runtime"
  "strings"
)

// Scanner scans an incoming HTTP request, parses it, sends it to the backend,
// and writes the formatted response to the HTTP response writer.
func Scanner(ctx context.Context, w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/plain; charset=utf-8")
  var (
    response = &responseBuilder{}
    req, err = parse(r)
  )

  if nil != err {
    slog.Error("Could not parse request", slog.Group("error", slog.String("message", err.Error())))
    response.WriteError(err)
    response.DefaultHeaders()
  } else {
    response = backend(ctx, req)
  }

  w.WriteHeader(http.StatusOK)
  template.HTMLEscape(w, response.Bytes())
}

// Renderer renders the website template and writes it to the HTTP response writer.
func Renderer(w http.ResponseWriter, r *http.Request) {
  abortWithAlert := func(alert string) {
    if "" != alert {
      alert = fmt.Sprint("<script>alert('", alert, "');</script>")
    }

    website("", "", alert).Render(r.Context(), w)
  }

  collname := r.URL.Query().Get("collname")
  if r.Method != http.MethodPost && "" == collname {
    abortWithAlert("")
    return
  }

  var (
    collfile io.ReadCloser
    err      error
  )

  if "" != collname { /* Lookup file in collections folder.  */
    var (
      _, current, _, _ = runtime.Caller(0)
      collfullname     = fmt.Sprint(path.Join(path.Dir(current), "public", "collections", collname), ".json")
      file             *os.File
    )

    file, err = os.Open(collfullname)
    if nil != err {
      slog.Error("could not open collection file", slog.Group("error", slog.String("message", err.Error())))
      abortWithAlert("Playground could find this collection :o")
      return
    }

    stats, err := file.Stat()
    if nil != err {
      slog.Error("file.Stat() failed", slog.Group("error", slog.String("message", err.Error())))
      abortWithAlert("Playground could find this collection :o")
      return
    }

    if 1024*1024 < stats.Size() {
      abortWithAlert("Playground only accepts file sizes less than 1 MB.")
      return
    }

    collfile = file /* All right.  */
  } else if r.Method == http.MethodPost { /* Collection file comes from form data.  */
    var collfileheader *multipart.FileHeader
    collfile, collfileheader, err = r.FormFile("coll")
    if nil != err {
      slog.Error("could not open collection file", slog.Group("error", slog.String("message", err.Error())))
      abortWithAlert("Playground could not open file.")
      return
    }

    if !strings.Contains(collfileheader.Header.Get("Content-Type"), "application/json") {
      abortWithAlert("Playground only accepts `application/json` files.")
      return
    }

    if 1024*1024 < collfileheader.Size {
      abortWithAlert("Playground only accepts file sizes less than 1 MB.")
      return
    }
  }

  defer collfile.Close()

  collsrc, colltree, err := collGen(collfile)
  if nil != err {
    slog.Error("could not generate from collection file", slog.Group("error", slog.String("message", err.Error())))
    abortWithAlert("Playground suffered an internal failure while processing your JSON file.")
    return
  }

  website(colltree, collsrc, "").Render(r.Context(), w)
}

// request represents an HTTP request with a method, target URL, and headers.
type request struct {
  // method is the HTTP method used for the request.
  method string

  // target is the URL to which the request is sent.
  target *url.URL

  // header contains the HTTP headers to be included in the request.
  header http.Header

  // The HTTP body of the request.
  body string
}

// parse extracts the HTTP method and target URL from an incoming HTTP request
// and returns a new request struct containing these values.
func parse(r *http.Request) (*request, error) {
  var err error
  req := new(request)

  if err = r.ParseForm(); nil != err {
    slog.Error(err.Error())
  }

  target := r.PostFormValue("request_target")
  method := r.PostFormValue("request_method")
  headerKeys := r.PostForm["header-key"]
  headerValues := r.PostForm["header-value"]
  if len(r.PostForm["http-request-body"]) > 0 {
    if 5<<20 <= len(r.PostForm["http-request-body"][0]) {
      return nil, errors.New("request body too long")
    }
    req.body = r.PostForm["http-request-body"][0]
  }

  req.header = http.Header{}
  for n := range len(headerKeys) {
    key, value := strings.TrimSpace(headerKeys[n]), strings.TrimSpace(headerValues[n])
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

  return req, nil
}
