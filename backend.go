package playground

import (
  "bytes"
  "encoding/json"
  "io"
  "log/slog"
  "mime"
  "net/http"
  "time"
)

// response represents the response returned from the backend.
type response struct {
  // body stores the formatted JSON response from the backend.
  body *bytes.Buffer
}

type bodyFormatter interface {
  format(input []byte, output io.Writer, indent string)
}

var (
  textFormatter = &textFormatterImpl{}
  xmlFormatter  = &xmlFormatterImpl{}
  jsonFormatter = &jsonFormatterImpl{}
  htmlFormatter = &htmlFormatterImpl{}
)

var supportedMediaTypes = map[string]bodyFormatter{
  "application/xml":          xmlFormatter,
  "application/problem+xml":  xmlFormatter,
  "application/json":         jsonFormatter,
  "application/problem+json": jsonFormatter,
  "application/sql":          textFormatter,
  "application/yaml":         textFormatter, // TODO: implement a YAML formatter

  "text/xml":  xmlFormatter,
  "text/html": htmlFormatter,
}

// backend sends an HTTP request to the target specified in the input request
// and returns a response with a formatted JSON body. If an error occurs during
// the request, it logs the error and returns nil.
func backend(in *request) *response {
  client := http.Client{
    Timeout: 30 * time.Second,
  }

  req, err := http.NewRequest(in.method, in.target.String(), nil)
  if nil != err {
    slog.Error(err.Error())
    return nil
  }

  for k, v := range in.header {
    for _, vv := range v {
      req.Header.Add(k, vv)
    }
  }

  out := &response{
    body: &bytes.Buffer{},
  }

  res, err := client.Do(req)
  if nil != err {
    slog.Error(err.Error())
    return nil
  }

  var (
    contentType     = res.Header.Get("Content-Type")
    mediatype, _, _ = mime.ParseMediaType(contentType)
    formatter       bodyFormatter
  )

  for mtype, f := range supportedMediaTypes {
    if mediatype == mtype {
      formatter = f
    }
  }

  if nil == formatter {
    if typ, _ := splitMediaType(mediatype); "text" != typ {
      out.body.WriteString("Unsupported media type:" + mediatype)
      return out
    }

    formatter = textFormatter
  }

  err = json.Indent(out.body, result, "", "  ")
  if nil != err {
    slog.Error(err.Error())
    return nil
  }

  formatter.format(result, out.body, "  ")

  return out
}

// splitMediaType extracts the type and subtype from a MIME type.
func splitMediaType(v string) (typ string, subtype string) {
  if parts := strings.Split(v, "/"); len(parts) == 2 {
    return parts[0], strings.Split(parts[1], ";")[0]
  }
  return
}
