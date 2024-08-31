package playground

import (
  "bytes"
  "encoding/json"
  "errors"
  "fmt"
  xhtml "golang.org/x/net/html"
  "html"
  "io"
  "log/slog"
  "regexp"
  "strings"
  "unicode/utf8"
)

type textFormatterImpl struct{}

// format trims left and right spaces from input and returns the actual content as is.
func (textFormatterImpl) format(input []byte, output io.Writer, _ string) {
  if len(input) == 0 {
    return
  }

  _, err := output.Write(bytes.TrimSpace(input))
  if nil != err {
    slog.Error(err.Error())
  }
}

type jsonFormatterImpl struct{}

// format implements a JSON formatter.
func (jsonFormatterImpl) format(input []byte, output io.Writer, indent string) {
  if len(input) == 0 {
    return
  }

  buffer := bytes.Buffer{}
  err := json.Indent(&buffer, input, "", indent)
  if nil != err {
    slog.Error(err.Error())
    return
  }

  _, err = output.Write(bytes.TrimSpace(buffer.Bytes()))
  if nil != err {
    slog.Error(err.Error())
  }
}

var (
  xmlReTag              = regexp.MustCompile(`<([/!]?)([^>]+?)(/?)>`)        // regexp xmlReTag matches an XML tag
  xmlReComment          = regexp.MustCompile(`(?s)(<!--)(.*?)(-->)`)         // regexp xmlReComment matches an XML comment
  xmlReXMLInsideComment = regexp.MustCompile(`<!--[^>]*?<[^>]+?[^-]*?-->`)   // regexp xmlReXMLInsideComment matches XML tags inside comments
  xmlReBlanksAround     = regexp.MustCompile(`\s*(<|/?>)\s*`)                // regexp xmlReBlanksAround matches blanks around XML tags
  xmlReBlanksInsideTags = regexp.MustCompile(`>([^<]*[\n\r\t]| {3,})[^<]*<`) // regexp xmlReBlanksInsideTags matches the content of an XML tags that contains repeated spaces or \r\n\t
  xmlReBlanks           = regexp.MustCompile(`\s{2,}`)                       // regexp xmlReBlanks matches the content of an XML tag that contains more than two spaces together
)

type xmlFormatterImpl struct{}

// format implements basic formatting for XML input by applying indentation, normalizing whitespace, handling
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

type htmlFormatterImpl struct{}

var (
  // htmlReStartingBlankBeforeText matches a string that starts with leading blank followed by a non-blank character
  htmlReStartingBlankBeforeText = regexp.MustCompile(`^\s+\S`)

  // htmlReTrailingBlankAfterText matches a string that ends with a non-blank character followed by trailing blank
  htmlReTrailingBlankAfterText = regexp.MustCompile(`\S\s+$`)

  // inlineTags is a dictionary that contains tags that typically do not start on a new line and only occupy the space
  // required by their content.
  inlineTags = map[string]struct{}{
    "a": {}, "b": {}, "i": {}, "em": {}, "strong": {}, "code": {}, "span": {},
    "ins": {}, "big": {}, "small": {}, "tt": {}, "abbr": {}, "acronym": {},
    "cite": {}, "dfn": {}, "kbd": {}, "samp": {}, "var": {}, "bdo": {},
    "map": {}, "q": {}, "sub": {}, "sup": {},
  }

  // selfClosingTags is a dictionary that contains tags that do not have a closing counterpart.
  selfClosingTags = map[string]struct{}{
    "input": {}, "link": {}, "meta": {}, "hr": {}, "img": {}, "br": {}, "area": {},
    "base": {}, "col": {}, "param": {}, "command": {}, "embed": {}, "keygen": {},
    "source": {}, "track": {}, "wbr": {},
  }
)

// isInlineTag checks if the given HTML tag is an inline tag.
func isInlineTag(tag []byte) bool {
  _, found := inlineTags[string(tag)]
  return found
}

// isSelfClosingTag checks if the given HTML tag is a self-closing tag.
func isSelfClosingTag(tag []byte) bool {
  _, found := selfClosingTags[string(tag)]
  return found
}

// adjustTextIndentation adjusts the indentation of a block of text. It first determines the minimum indentation level
// (leading spaces) present in the non-empty lines of the text; then it removes this common indentation and re-indents
// the text with a specified depth of spaces.
func adjustTextIndentation(text []byte, depth int) []byte {
  minIndent := 1000

  for _, line := range bytes.Split(text, []byte{'\n'}) {
    if trimmed := bytes.TrimLeft(line, " "); len(trimmed) > 0 {
      indent := len(line) - len(trimmed)
      if indent < minIndent {
        minIndent = indent
      }
    }
  }

  if minIndent == 1000 {
    minIndent = 0
  }

  re := regexp.MustCompile(fmt.Sprintf(`\n\s{%d}`, minIndent))
  return re.ReplaceAllLiteral(text, append([]byte{'\n'}, bytes.Repeat([]byte("  "), depth)...))
}

// format formats the HTML input with proper indentation.
func (htmlFormatterImpl) format(input []byte, output io.Writer, indent string) {
  var (
    reader           = bytes.NewReader(input)
    tokenizer        = xhtml.NewTokenizer(reader)
    depth            = 0
    prevType         = xhtml.ErrorToken
    tagName          []byte
    prvName          []byte
    longText         = false
    skipFirstNewline = true
  )

  writeIndent := func(depth int) {
    if skipFirstNewline {
      skipFirstNewline = false
    } else {
      output.Write([]byte("\n"))
    }

    output.Write(bytes.Repeat([]byte(indent), depth))
  }

  for {
    tokenType := tokenizer.Next()

    if tokenType != xhtml.TextToken {
      prvName = tagName
      tagName, _ = tokenizer.TagName()
    }

    switch tokenType {
    case xhtml.ErrorToken:
      if errors.Is(tokenizer.Err(), io.EOF) {
        return
      }

      slog.Error(tokenizer.Err().Error())
      return

    case xhtml.StartTagToken:
      if !(isInlineTag(tagName) && prevType == xhtml.TextToken) {
        writeIndent(depth)
      }

      output.Write(tokenizer.Raw())
      if !isSelfClosingTag(tagName) {
        depth++
      }

    case xhtml.SelfClosingTagToken, xhtml.CommentToken, xhtml.DoctypeToken:
      writeIndent(depth)
      output.Write(tokenizer.Raw())

    case xhtml.EndTagToken:
      if depth > 0 {
        depth--
      }

      if !bytes.Equal(prvName, tagName) || prevType == xhtml.SelfClosingTagToken || prevType == xhtml.CommentToken ||
        prevType == xhtml.DoctypeToken || (prevType == xhtml.TextToken && longText) {
        writeIndent(depth)
      }

      output.Write(tokenizer.Raw())

    case xhtml.TextToken:
      t := bytes.Replace(tokenizer.Raw(), []byte{'\t'}, []byte(indent), -1)
      text := bytes.Trim(t, "\n\r ")

      if htmlReTrailingBlankAfterText.Match(t) {
        text = append(text, ' ')
      }

      longText = false
      if len(text) > 0 {
        if bytes.Contains(text, []byte{'\n'}) {
          if !(prevType == xhtml.EndTagToken && isInlineTag(tagName)) {
            writeIndent(depth)
          } else if htmlReStartingBlankBeforeText.Match(t) {
            text = append([]byte{' '}, text...)
          }

          output.Write(adjustTextIndentation(text, depth))
          longText = true
        } else {
          if utf8.RuneCount(text) > 80 || prevType != xhtml.StartTagToken {
            if !(prevType == xhtml.EndTagToken && isInlineTag(tagName)) {
              writeIndent(depth)
              longText = true
            } else if htmlReStartingBlankBeforeText.Match(t) {
              text = append([]byte{' '}, text...)
            }
          }
          output.Write(text)
        }
      }
    }

    prevType = tokenType
  }
}
