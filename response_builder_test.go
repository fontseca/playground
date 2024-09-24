package playground

import (
  "bytes"
  "errors"
  "maps"
  "net/http"
  "reflect"
  "slices"
  "strings"
  "testing"
)

func TestNewResponseBuilder(t *testing.T) {
  r := newResponseBuilder()

  if 0 != r.body.Len() {
    t.Errorf("r.body.Len()=%d, want %d", r.body.Len(), 0)
  }

  if !bytes.Equal([]byte(""), r.startLine) {
    t.Errorf("startLine is not empty")
  }

  if r.header == nil {
    t.Errorf("header map is uninitialized")
  }
}

func TestResponseBuilder_SetStartLine(t *testing.T) {
  r := newResponseBuilder()
  r.SetStartLine("HTTP/1.1", "404 Not Found")
  if !bytes.Equal([]byte("HTTP/1.1 404 Not Found"), r.startLine) {
    t.Errorf("expected: %s, got: %s", "HTTP/1.1 404 Not Found", r.startLine)
  }
}

func TestResponseBuilder_SetHeader(t *testing.T) {
  r := newResponseBuilder()
  header := http.Header{}
  header.Add("one", "1")
  header.Add("two", "2")
  r.SetHeaders(header)

  if !reflect.DeepEqual(r.header, header) {
    t.Errorf("expected: %v, got: %v", header, r.header)
  }
}

func TestResponseBuilder_DefaultHeaders(t *testing.T) {
  r := newResponseBuilder()
  header := http.Header{}
  header.Add("one", "1")
  header.Add("two", "2")
  r.SetHeaders(header)
  r.DefaultHeaders()

  keys := []string{"Content-Length", "Content-Type", "Date", "Server"}

  for k := range slices.Values(keys) {
    if _, ok := r.header[k]; !ok {
      t.Errorf("missing default header: %s", k)
    }
  }

  for k := range maps.Keys(r.header) {
    if !slices.Contains(keys, k) {
      t.Errorf("superfluous default header: %s", k)
    }
  }
}

func TestResponseBuilder_Write(t *testing.T) {
  r := newResponseBuilder()
  _, _ = r.Write([]byte("\n"))
  _, _ = r.Write([]byte("line1"))
  _, _ = r.Write([]byte("\n"))
  _, _ = r.Write([]byte("line2"))
  _, _ = r.Write([]byte("\n"))

  expected := `
line1
line2
`

  if expected != r.body.String() {
    t.Errorf("expected:"+
      "\n---%s---\n"+
      "got:"+
      "\n---%s---\n", expected, r.body.String())
  }
}

func TestResponseBuilder_WriteError(t *testing.T) {
  r := newResponseBuilder()
  _, _ = r.Write([]byte("line1\nline2\nline3"))

  expected := "error message"
  r.WriteError(errors.New(expected))

  if expected != r.body.String() {
    t.Errorf("expected:"+
      "\n---%s---\n"+
      "got:"+
      "\n---%s---\n", expected, r.body.String())
  }
}

func makeBuilderForTest() (r *responseBuilder, expected string) {
  r = newResponseBuilder()
  h := http.Header{}
  h.Set("Content-Type", "text/plain; charset=utf-8")
  h.Set("Date", "Mon, 09 Sep 2024 13:39:17 CST")
  h.Set("Server", "fontseca.dev/playground (v1.0)")
  r.SetHeaders(h)

  _, _ = r.Write([]byte(`Lorem ipsum dolor sit amet, consectetur adipiscing
elit. Proin gravida sed eros vel posuere. In hac
habitasse platea dictumst. Nunc facilisis nunc eget
egestas malesuada. Nulla mattis arcu vel volutpat
auctor. Phasellus id dui vitae erat malesuada faucibus
sed quis diam. Phasellus finibus, augue nec luctus.`))

  expected = `HTTP/1.0 200 OK
Content-Type: text/plain; charset=utf-8
Date: Mon, 09 Sep 2024 13:39:17 CST
Server: fontseca.dev/playground (v1.0)

Lorem ipsum dolor sit amet, consectetur adipiscing
elit. Proin gravida sed eros vel posuere. In hac
habitasse platea dictumst. Nunc facilisis nunc eget
egestas malesuada. Nulla mattis arcu vel volutpat
auctor. Phasellus id dui vitae erat malesuada faucibus
sed quis diam. Phasellus finibus, augue nec luctus.`

  return r, expected
}

func TestResponseBuilder_build(t *testing.T) {
  r, expected := makeBuilderForTest()
  buffer := r.build()

  if expected != buffer.String() {
    t.Errorf("\nexpected:"+
      "\n---\n%s\n---\n"+
      "\ngot:"+
      "\n---\n%s\n---\n", expected, buffer.String())
  }

  r = newResponseBuilder()
  h := http.Header{}
  h.Set("Content-Type", "text/plain; charset=utf-8")
  h.Set("Date", "Mon, 09 Sep 2024 13:39:17 CST")
  h.Set("Server", "fontseca.dev/playground (v1.0)")
  h.Add("X-Value", "1")
  h.Add("X-Value", "2")
  h.Add("X-Value", "3")
  h.Add("X-Value", "4")
  r.SetHeaders(h)

  r.Write([]byte("Lorem ipsum dolor sit amet, consectetur adipiscing"))

  expected = `HTTP/1.0 200 OK
Content-Type: text/plain; charset=utf-8
Date: Mon, 09 Sep 2024 13:39:17 CST
Server: fontseca.dev/playground (v1.0)
X-Value: 1
X-Value: 2
X-Value: 3
X-Value: 4

Lorem ipsum dolor sit amet, consectetur adipiscing`

  buffer = r.build()

  if expected != buffer.String() {
    t.Errorf("\nexpected:"+
      "\n---\n%s\n---\n"+
      "\ngot:"+
      "\n---\n%s\n---\n", expected, buffer.String())
  }

}

func TestResponseBuilder_Bytes(t *testing.T) {
  r, expected := makeBuilderForTest()
  got := r.Bytes()

  if !bytes.Equal([]byte(expected), got) {
    t.Errorf("\nexpected:"+
      "\n---\n%s\n---\n"+
      "\ngot:"+
      "\n---\n%s\n---\n", expected, got)
  }
}

func TestResponseBuilder_String(t *testing.T) {
  r, expected := makeBuilderForTest()
  got := r.String()

  if strings.Compare(expected, got) != 0 {
    t.Errorf("\nexpected:"+
      "\n---\n%s\n---\n"+
      "\ngot:"+
      "\n---\n%s\n---\n", expected, got)
  }
}
