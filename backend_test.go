package playground

import (
  "testing"
)

func TestSplitMediaType(t *testing.T) {
  tests := [...][3]string{
    {"", "", ""},
    {"abc", "", ""},
    {" \t\n\t ", "", ""},
    {";", "", ""},
    {"a//b", "", ""},
    {"application/problem+json", "application", "problem+json"},
    {"application/sql; charset=utf-8", "application", "sql"},
    {"application/vnd.api+json", "application", "vnd.api+json"},
    {"application/x-stuff; title*=us-ascii'en-us'This%20is%20%2A%2A%2Afun%2A%2A%2A", "application", "x-stuff"},
    {"image/avif", "image", "avif"},
    {"text/prs.prop.logic", "text", "prs.prop.logic"},
    {"image/svg+xml", "image", "svg+xml"},
    {"text/html;\tcharset=utf-8", "text", "html"},
    {"unregistered/prs.unregistered", "unregistered", "prs.unregistered"},
  }

  for _, test := range tests {
    contentType, contentSubtype := splitMediaType(test[0])
    expectedContentType := test[1]
    expectedContentSubtype := test[2]
    if contentType != expectedContentType || contentSubtype != expectedContentSubtype {
      t.Errorf("splitMediaType(%q) = (%q, %q), want (%q, %q)",
        test[0], contentType, contentSubtype, expectedContentType, expectedContentSubtype)
    }
  }
}
