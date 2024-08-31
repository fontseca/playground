package playground

import (
  "bytes"
  "html"
  "io"
  "log/slog"
  "regexp"
  "strings"
)

var (
  xmlReTag              = regexp.MustCompile(`<([/!]?)([^>]+?)(/?)>`)        // regexp xmlReTag matches an XML tag
  xmlReComment          = regexp.MustCompile(`(?s)(<!--)(.*?)(-->)`)         // regexp xmlReComment matches an XML comment
  xmlReXMLInsideComment = regexp.MustCompile(`<!--[^>]*?<[^>]+?[^-]*?-->`)   // regexp xmlReXMLInsideComment matches XML tags inside comments
  xmlReBlanksAround     = regexp.MustCompile(`\s*(<|/?>)\s*`)                // regexp xmlReBlanksAround matches blanks around XML tags
  xmlReBlanksInsideTags = regexp.MustCompile(`>([^<]*[\n\r\t]| {3,})[^<]*<`) // regexp xmlReBlanksInsideTags matches the content of an XML tags that contains repeated spaces or \r\n\t
  xmlReBlanks           = regexp.MustCompile(`\s{2,}`)                       // regexp xmlReBlanks matches the content of an XML tag that contains more than two spaces together
)

type xmlFormatterImpl struct{}

// formatXML implements basic formatting for XML input by applying indentation, normalizing whitespace, handling
// comments and XML tags inside comments.
func (xmlFormatterImpl) format(input []byte, output io.Writer, indent string) {
  if nil == input || len(input) < 2 {
    return
  }

  var needsUnescape bool

  // Comments might contain further XML code. In that's the case, we want to escape
  // that code to avoid formatting.
  out := xmlReXMLInsideComment.ReplaceAllFunc(input, func(comment []byte) []byte {
    needsUnescape = true
    submatches := xmlReComment.FindSubmatch(comment)
    b := bytes.Buffer{}
    b.Grow(len(comment))
    b.Write(submatches[1])                                  // <!--
    b.WriteString(html.EscapeString(string(submatches[2]))) // ... (which includes XML code)
    b.Write(submatches[3])                                  // -->
    return b.Bytes()
  })

  out = xmlReBlanksAround.ReplaceAll(out, []byte("$1"))
  out = xmlReBlanksInsideTags.ReplaceAllFunc(out, func(m []byte) []byte {
    return xmlReBlanks.ReplaceAll(m, []byte(" "))
  })

  out = xmlReTag.ReplaceAllFunc(out, xmlTagReplacer(indent))

  if needsUnescape {
    // restore the original comment escaped content
    out = xmlReComment.ReplaceAllFunc(out, func(comment []byte) []byte {
      submatches := xmlReComment.FindSubmatch(comment)
      b := bytes.Buffer{}
      b.Grow(len(comment))
      b.Write(submatches[1])                                    // <!--
      b.WriteString(html.UnescapeString(string(submatches[2]))) // ... (which include XML code)
      b.Write(submatches[3])                                    // -->
      return b.Bytes()
    })
  }

  _, err := output.Write(out[1:])
  if nil != err {
    slog.Error(err.Error())
  }
}

// xmlTagReplacer returns a closure that processes XML tag slices, applying indentation and
// whitespace normalization based on the type of XML tag provided.
func xmlTagReplacer(indent string) func([]byte) []byte {
  var (
    depth       = 0
    wasEndOfTag = true
    buffer      = bytes.Buffer{}
  )

  buffer.Grow(64)

  return func(tag []byte) []byte {
    // XML declaration
    if bytes.HasPrefix(tag, []byte("<?xml")) {
      return append([]byte("\n"), xmlReBlanks.ReplaceAll(tag, []byte(" "))...)
    }

    defer buffer.Reset()
    buffer.Write([]byte("\n"))

    switch {
    default: // start element
      wasEndOfTag = false
      buffer.WriteString(strings.Repeat(indent, depth))
      buffer.Write(xmlReBlanks.ReplaceAll(tag, []byte(" ")))
      depth++
      return buffer.Bytes()

    case bytes.HasSuffix(tag, []byte("/>")): // empty element
      wasEndOfTag = true
      buffer.WriteString(strings.Repeat(indent, depth))
      buffer.Write(tag)
      return buffer.Bytes()

    case bytes.HasPrefix(tag, []byte("<!")): // comment or doctype
      if bytes.HasPrefix(bytes.ToLower(tag), []byte("<!doctype")) {
        tag = xmlReBlanks.ReplaceAll(tag, []byte(" "))
      }

      wasEndOfTag = true
      buffer.WriteString(strings.Repeat(indent, depth))
      buffer.Write(tag)
      return buffer.Bytes()

    case bytes.HasPrefix(tag, []byte("</")): // end element
      depth--
      if wasEndOfTag {
        buffer.WriteString(strings.Repeat(indent, depth))
        buffer.Write(tag)
        return buffer.Bytes()
      }
      wasEndOfTag = true
      return tag
    }
  }
}
