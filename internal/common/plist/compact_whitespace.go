package plist

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
)

// structuralPlistElements are the plist container elements whose direct
// whitespace-only child text nodes are pure formatting: every plist parser
// discards them. Whitespace inside any other element (<string>, <key>,
// <data>, <integer>, …) is element content and must never be touched.
var structuralPlistElements = map[string]bool{
	"plist": true,
	"dict":  true,
	"array": true,
}

// CompactStructuralWhitespace removes whitespace-only text that sits between
// tags inside plist structural elements (<plist>, <dict>, <array>) and in the
// document prolog/epilog, collapsing a pretty-printed plist payload onto a
// single line. Everything else passes through byte-for-byte: leaf element
// content (<string>, <key>, <data>, …), comments, CDATA sections,
// entity/character references and dict key order are left untouched — the
// document is never re-serialised, only whitespace-only byte ranges of the
// original input are removed.
//
// Why this exists: the Jamf Pro Classic API's server-side plist parser
// materialises the whitespace text nodes between sibling <array> tags as
// phantom empty <array/> entries in the stored plist on write. For a
// com.apple.homescreenlayout payload that inserts a blank leading home-screen
// page, pushing all real content onto page 2 (the "empty page 1" symptom).
// Plist semantics ignore inter-element whitespace inside containers, so
// stripping it sidesteps the whole bug class; the wire payload needs no
// readability. The approach mirrors the fix in the Jamf-maintained
// jamfplatform provider (Jamf-Concepts/terraform-provider-jamfplatform@8126b1b).
//
// On malformed XML the input is returned unchanged alongside the error, so
// callers can fall back to sending the original bytes and let the server
// report the malformation with its own error.
func CompactStructuralWhitespace(raw []byte) ([]byte, error) {
	dec := xml.NewDecoder(bytes.NewReader(raw))
	// The decoder knows only the five predefined XML entities by default, and
	// it never fetches the DTD referenced by the plist DOCTYPE. Without this a
	// payload containing a named entity (&nbsp;, &copy;, …) fails to tokenise,
	// compaction bails out, and the profile silently ships uncompacted — i.e.
	// the phantom <array/> bug reappears with only a [WARN] to show for it.
	dec.Entity = xml.HTMLEntity

	var stack []string
	type span struct{ start, end int64 }
	var cuts []span

	for {
		start := dec.InputOffset()
		tok, err := dec.Token()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return raw, fmt.Errorf("tokenising plist XML: %w", err)
		}
		end := dec.InputOffset()

		switch t := tok.(type) {
		case xml.StartElement:
			stack = append(stack, t.Name.Local)
		case xml.EndElement:
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
		case xml.CharData:
			// An empty stack means the prolog or epilog — text there is
			// insignificant, so it is compacted like any structural context.
			structural := len(stack) == 0 || structuralPlistElements[stack[len(stack)-1]]
			if !structural {
				continue
			}
			if start < end && isXMLWhitespace(raw[start:end]) {
				cuts = append(cuts, span{start: start, end: end})
			}
		}
	}

	if len(cuts) == 0 {
		return raw, nil
	}

	out := make([]byte, 0, len(raw))
	var prev int64
	for _, c := range cuts {
		out = append(out, raw[prev:c.start]...)
		prev = c.end
	}
	out = append(out, raw[prev:]...)
	return out, nil
}

// isXMLWhitespace reports whether every byte is one of the four XML whitespace
// characters (space, tab, carriage return, line feed).
func isXMLWhitespace(b []byte) bool {
	for _, c := range b {
		switch c {
		case ' ', '\t', '\r', '\n':
		default:
			return false
		}
	}
	return true
}
