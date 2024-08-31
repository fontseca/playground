package playground

import (
  "bytes"
  "encoding/json"
  "io"
  "log/slog"
  "net/http"
  "time"
)

// response represents the response returned from the backend.
type response struct {
  // body stores the formatted JSON response from the backend.
  body *bytes.Buffer
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

  res, err := client.Do(req)
  if nil != err {
    slog.Error(err.Error())
    return nil
  }

  result, err := io.ReadAll(res.Body)
  defer res.Body.Close()

  out := &response{
    body: &bytes.Buffer{},
  }

  err = json.Indent(out.body, result, "", "  ")
  if nil != err {
    slog.Error(err.Error())
    return nil
  }

  return out
}

// splitMediaType extracts the type and subtype from a MIME type.
func splitMediaType(v string) (typ string, subtype string) {
  if parts := strings.Split(v, "/"); len(parts) == 2 {
    return parts[0], strings.Split(parts[1], ";")[0]
  }
  return
}
