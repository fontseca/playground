package playground

import (
  "encoding/json"
  "fmt"
  "github.com/google/uuid"
  "io"
  "regexp"
  "slices"
  "strconv"
  "strings"
)

// A collQueryParam is key-value representation of a query parameter in a collURL.
type collQueryParam struct {
  Key   string `json:"key"`
  Value string `json:"value"`
}

// A collURL contains the complete broken-down URL for this a collRequest.
type collURL struct {
  Raw      string           `json:"raw"`      // The string representation of the request URL, including the protocol, host, path, hash, query parameter(s) and path variable(s).
  Protocol string           `json:"protocol"` // The protocol associated with the request, E.g: 'http'.
  Host     []string         `json:"host"`     // The host for the URL, E.g: api.yourdomain.com. Can be stored as a string or as an array of strings.
  Path     []string         `json:"path"`     // The complete path of the current url, broken down into segments. A segment could be a string, or a path variable.
  Port     string           `json:"port"`     // The port number present in this URL.
  Query    []collQueryParam `json:"query"`    // An array of collQueryParam, which is basically the query string part of the URL, parsed into separate variables

}

// A collHeader represents a single HTTP Header.
type collHeader struct {
  Key      string `json:"key"`
  Value    string `json:"value"`
  Disabled bool   `json:"disabled"` //  If set to true, the current header should not be sent with requests.
}

type collURLEncodedParameter struct {
  Key   string `json:"key"`
  Value string `json:"value"`
}

// collBody contains all the data needed for a request body.
type collBody struct {
  Mode       string                    `json:"mode"` // The type of data associated with this request in this field. One of: raw, urlencoded, formdata, file or graphql.
  Raw        string                    `json:"raw"`
  URLEncoded []collURLEncodedParameter `json:"urlencoded"`
}

// A collRequest represents an HTTP request.
type collRequest struct {
  URL    collURL      `json:"url"`
  Method string       `json:"method"`
  Header []collHeader `json:"header"` // A representation for a list of headers.
  Body   *collBody    `json:"body"`   // Holds the data contained in the request body.
}

// collItem are entities which contain an actual HTTP request.
type collItem struct {
  ID      string       `json:"id"`   // A unique ID that is used to identify collections internally.
  Name    string       `json:"name"` // A human-readable identifier for the current item.
  Item    []collItem   `json:"item"` // If not empty, the collItem is a folder, which may contain many `collItem`s.
  Request *collRequest `json:"request"`
}

// A collVariable allows you to define a set of variables, that are a part of the collection.
type collVariable struct {
  ID       string `json:"id"`    // A unique user-defined value that identifies the variable within a collection.
  Key      string `json:"key"`   // A human friendly value that identifies the variable within a collection.
  Value    string `json:"value"` // The value that a variable holds in the coll struct. Ultimately, the variables will be replaced by this value, when say running a set of requests from a collection
  Type     string `json:"type"`  // Specifies the type of the variable: string, boolean, any or number.
  Name     string `json:"name"`  // Variable name.
  Disabled bool   `json:"disabled"`
}

// coll holds (some) fields of the Postman Collection Format v2.1.0.
type coll struct {
  Info struct {
    Name string `json:"name"` // A collection's friendly name is defined by this field.
  }
  Item     []collItem     `json:"item"`
  Variable []collVariable `json:"variable"`
}

var (
  regexpVariable = regexp.MustCompile(`\{\{([a-zA-Z0-9_]+)\}\}`)
  // newString holds a reference to uuid.NewString, but it is replaced by a mock function in testing,
  // so that deterministic behaviour is assured.
  newString = uuid.NewString
)

// writePair writes a JSON key-value pair to a strings.Builder.
// This is used to construct JSON objects for requests and responses.
func writePair(builder *strings.Builder, key, value string) {
  builder.WriteByte('"')
  builder.WriteString(key)
  builder.WriteString(`":`)
  builder.WriteString(strconv.Quote(fmt.Sprint("", value, "")))
  builder.WriteByte(',')
}

// walk recursively processes collItems and their sub-items (folders), generating
// HTML tree structures and JSON request arrays. It resolves URL variables using
// the provided map and generates a unique ID for each item.
func walk(variables map[string]collVariable, array *strings.Builder, dirtree *strings.Builder, fullItemName string, item []collItem) {
  for i := range slices.Values(item) {
    if len(i.Item) > 0 { /* folder */
      dirtree.WriteString(fmt.Sprintf(
        `<div class="item folder">`+
          "<span class=\"name\">%s</span>", i.Name))
      walk(variables, array, dirtree, fmt.Sprint(fullItemName, i.Name, " / "), i.Item)
      dirtree.WriteString("</div>")
    } else {
      resolvedUrl := regexpVariable.ReplaceAllStringFunc(i.Request.URL.Raw, func(match string) string {
        vname := match[2 : len(match)-2]

        if replacement, exists := variables[vname]; exists {
          return replacement.Value
        }

        return match
      })

      i.ID = newString()

      array.WriteByte('{')

      writePair(array, "id", i.ID)
      writePair(array, "name", i.Name)
      writePair(array, "full_name", fmt.Sprint(fullItemName, i.Name))

      if nil != i.Request {

        writePair(array, "request_method", i.Request.Method)

        array.WriteString(`"request_header":`)
        array.WriteByte('[')

        for u := range slices.Values(i.Request.Header) {
          array.WriteByte('{')
          writePair(array, "key", u.Key)
          writePair(array, "value", u.Value)
          array.WriteByte('}')
          array.WriteByte(',')
        }

        array.WriteByte(']')
        array.WriteByte(',')

        if nil != i.Request.Body {
          if "" != i.Request.Body.Mode {
            writePair(array, "request_body_mode", i.Request.Body.Mode)
          }

          if "" != i.Request.Body.Raw {
            writePair(array, "request_body_raw", i.Request.Body.Raw)
          }

          if len(i.Request.Body.URLEncoded) > 0 {
            array.WriteString(`"request_body_urlencoded":`)
            array.WriteByte('[')

            for u := range slices.Values(i.Request.Body.URLEncoded) {
              array.WriteByte('{')
              writePair(array, "key", u.Key)
              writePair(array, "value", u.Value)
              array.WriteByte('}')
              array.WriteByte(',')
            }

            array.WriteByte(']')
            array.WriteByte(',')
          }
        }

        writePair(array, "url_raw", i.Request.URL.Raw)
        writePair(array, "url_port", i.Request.URL.Port)
        writePair(array, "url_protocol", i.Request.URL.Protocol)
        array.WriteString(`"url_query":`)
        array.WriteByte('[')

        for u := range slices.Values(i.Request.URL.Query) {
          array.WriteByte('{')
          writePair(array, "key", u.Key)
          writePair(array, "value", u.Value)
          array.WriteByte('}')
          array.WriteByte(',')
        }

        array.WriteByte(']')
        array.WriteByte(',')

      }

      writePair(array, "url_resolved", resolvedUrl)
      array.WriteByte('}')
      array.WriteByte(',')

      dirtree.WriteString(fmt.Sprintf("<div class=\"item\"><span data-id=\"%s\" class=\"name\">%s</span></div>", i.ID, i.Name))
    }
  }
}

// parseColl parses a JSON input representing a collection file and converts it
// into a coll struct, returning any decoding errors encountered.
func parseColl(collfile io.Reader) (c *coll, err error) {
  decoder := json.NewDecoder(collfile)
  c = &coll{}

  if err = decoder.Decode(&c); nil != err {
    return nil, err
  }

  return c, nil
}

// collGen generates JavaScript and HTML snippets from a collection file, producing
// a JSON array of requests and an HTML directory tree of requests and folders.
// It uses variables for URL resolution within requests.
func collGen(collfile io.Reader) (collsrc, colldirtree string, err error) {
  c, err := parseColl(collfile)
  if nil != err {
    return "", "", err
  }

  var (
    requestsArrayBuilder = &strings.Builder{}
    dirtreeBuilder       = &strings.Builder{}
    variables            = make(map[string]collVariable)
  )

  for v := range slices.Values(c.Variable) {
    variables[v.Key] = v
  }

  dirtreeBuilder.WriteString("<header><h3>")
  dirtreeBuilder.WriteString(c.Info.Name)
  dirtreeBuilder.WriteString("</h3></header>")

  walk(variables, requestsArrayBuilder, dirtreeBuilder, "", c.Item)
  collsrc = fmt.Sprintf("<script>const requests = [%s];</script>", requestsArrayBuilder.String())
  return collsrc, dirtreeBuilder.String(), nil
}
