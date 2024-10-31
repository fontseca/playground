package playground

import (
  "compress/flate"
  "compress/gzip"
  "context"
  "errors"
  "fmt"
  "io"
  "log/slog"
  "maps"
  "mime"
  "net"
  "net/http"
  "net/url"
  "strings"
  "time"
)

var allowedMethods = map[string]struct{}{
  http.MethodGet:    {},
  http.MethodPost:   {},
  http.MethodPut:    {},
  http.MethodPatch:  {},
  http.MethodDelete: {},
}

// maxBodyBytes is the accepted body size for a playground request.
const maxBodyBytes = 5 << 20 // 5 MB

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

var supportedEncodings = map[string]func(io.Reader) io.ReadCloser{
  "gzip":     func(r io.Reader) io.ReadCloser { out, _ := gzip.NewReader(r); return out },
  "flate":    flate.NewReader,
  "compress": flate.NewReader,
}

var errNoRequest = errors.New("internal service error")

// backend sends an HTTP request to the target specified in the input request
// and returns a playgroundResponse with a formatted JSON body. If an error occurs during
// the request, it logs the error and returns nil.
func backend(ctx context.Context, in *request) (response *responseBuilder) {
  client := &http.Client{
    Timeout: 30 * time.Second,
    Transport: &http.Transport{
      DialContext: (&net.Dialer{
        Timeout:   10 * time.Second,
        KeepAlive: 10 * time.Second,
      }).DialContext,
      TLSHandshakeTimeout: 10 * time.Second,
    },
    CheckRedirect: func(req *http.Request, via []*http.Request) error {
      if len(via) >= 5 {
        return http.ErrUseLastResponse
      }
      if !isValidTarget(req.URL) {
        return fmt.Errorf("redirect to unsafe URL blocked")
      }
      return nil
    },
  }

  response = newResponseBuilder()

  ctx, cancel := context.WithTimeout(ctx, client.Timeout)
  defer cancel()

  if _, ok := allowedMethods[in.method]; !ok {
    response.WriteError(fmt.Errorf("method %s is not allowed", in.method))
    response.DefaultHeaders()
    return
  }

  if !isValidTarget(in.target) {
    response.WriteError(fmt.Errorf("invalid target URL"))
    response.DefaultHeaders()
    return
  }

  var body io.Reader

  if "" != in.body {
    body = io.LimitReader(strings.NewReader(in.body), maxBodyBytes)
  }

  req, err := http.NewRequestWithContext(ctx, in.method, in.target.String(), body)
  if nil != err {
    slog.Error("http.NewRequestWithContext(...) failed",
      slog.Group("error", slog.String("message", err.Error())),
      slog.Group("request",
        slog.String("method", in.method),
        slog.String("target", in.target.String()),
        slog.Int("n_headers", len(in.header)),
        slog.String("headers", fmt.Sprintf("%v", in.header)),
      ),
    )

    response.WriteError(errNoRequest)
    response.DefaultHeaders()
    return
  }

  maps.Copy(req.Header, in.header)

  res, err := client.Do(req)
  if nil != err {
    switch {
    default:
      slog.Error("client.Do(...) failed", slog.Group("error", slog.String("message", err.Error())))
      response.WriteError(errNoRequest)
    case errors.Is(err, context.DeadlineExceeded):
      response.WriteError(errors.New("request timed out"))
    case strings.Contains(err.Error(), "unsupported protocol scheme"):
      response.WriteError(errors.New(strings.Split(err.Error(), ": ")[1]))
    case strings.Contains(err.Error(), "no such host"):
      response.WriteError(fmt.Errorf("could not connect to the server at %#q", in.target.String()))
    }

    response.DefaultHeaders()
    return
  }

  response.SetStartLine(res.Proto, res.Status)
  response.SetHeaders(res.Header)

  var (
    contentType     = res.Header.Get("Content-Type")
    mediatype, _, _ = mime.ParseMediaType(contentType)
    formatter       bodyFormatter
  )

  for mtype, f := range maps.All(supportedMediaTypes) {
    if mediatype == mtype {
      formatter = f
    }
  }

  if nil == formatter {
    if typ, _ := splitMediaType(mediatype); "text" != typ && "" != typ {
      response.WriteError(fmt.Errorf("unsupported media type %#q", mediatype))
      response.DefaultHeaders()
      return
    }

    formatter = textFormatter
  }

  res.Body = http.MaxBytesReader(nil, res.Body, maxBodyBytes)
  var bodyReader io.ReadCloser

  contentEncoding := res.Header.Get("Content-Encoding")
  for encoding, maker := range supportedEncodings {
    if contentEncoding == encoding {
      bodyReader = maker(res.Body)
    }
  }

  if nil == bodyReader {
    bodyReader = res.Body
  }

  defer bodyReader.Close()

  result, err := io.ReadAll(bodyReader)
  if nil != err {
    switch {
    default:
      response.WriteError(errNoRequest)
      slog.Error("io.ReadAll(...) failed", slog.Group("error", slog.String("message", err.Error())))
    case strings.Contains(err.Error(), "request body too large"):
      response.WriteError(fmt.Errorf("response body is too large"))
    }

    response.DefaultHeaders()
    return
  }

  formatter.format(result, response, "  ")

  return response
}

// splitMediaType extracts the type and subtype from a MIME type.
func splitMediaType(v string) (typ string, subtype string) {
  if parts := strings.Split(v, "/"); len(parts) == 2 {
    return parts[0], strings.Split(parts[1], ";")[0]
  }
  return
}

// isValidTarget checks if the URL is safe (not internal IPs or localhost)
func isValidTarget(target *url.URL) bool {
  ip := net.ParseIP(target.Hostname())
  if ip == nil {
    return true // Not an IP address, proceed as normal.
  }

  privateRanges := []*net.IPNet{
    {IP: net.IPv4(127, 0, 0, 0), Mask: net.CIDRMask(8, 32)},    // 127.0.0.0/8
    {IP: net.IPv4(10, 0, 0, 0), Mask: net.CIDRMask(8, 32)},     // 10.0.0.0/8
    {IP: net.IPv4(172, 16, 0, 0), Mask: net.CIDRMask(12, 32)},  // 172.16.0.0/12
    {IP: net.IPv4(192, 168, 0, 0), Mask: net.CIDRMask(16, 32)}, // 192.168.0.0/16
    {IP: net.ParseIP("::1"), Mask: net.CIDRMask(128, 128)},     // IPv6 loopback
  }

  for _, rng := range privateRanges {
    if rng.Contains(ip) {
      return false // Block requests to internal IP ranges
    }
  }
  return true
}
