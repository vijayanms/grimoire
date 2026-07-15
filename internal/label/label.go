package label

import (
	"regexp"
	"strings"
	"unicode"
)

var (
	nonAlphanumRE = regexp.MustCompile(`[^a-z0-9_]+`)
	multiUnderRE  = regexp.MustCompile(`_+`)
)

func sanitizeLabel(s string) string {
	s = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return unicode.ToLower(r)
		}
		return '_'
	}, s)
	s = multiUnderRE.ReplaceAllString(s, "_")
	s = strings.Trim(s, "_")
	if s == "" {
		s = "resource"
	}
	if unicode.IsDigit(rune(s[0])) {
		s = "_" + s
	}
	if len(s) > 63 {
		s = s[:63]
	}
	return s
}

// Tracker tracks labels per resource type to avoid collisions.
type Tracker struct {
	seen map[string]int
}

func New() *Tracker {
	return &Tracker{seen: make(map[string]int)}
}

// Derive returns a unique label for a resource given candidate name, description, and uuid.
func (t *Tracker) Derive(name, description, uuid string) string {
	var base string
	switch {
	case name != "":
		base = sanitizeLabel(name)
	case description != "":
		d := description
		if len(d) > 20 {
			d = d[:20]
		}
		base = sanitizeLabel(d)
	default:
		if len(uuid) >= 8 {
			base = sanitizeLabel(uuid[:8])
		} else {
			base = sanitizeLabel(uuid)
		}
	}

	if _, exists := t.seen[base]; !exists {
		t.seen[base] = 1
		return base
	}
	t.seen[base]++
	return base + "_" + strings.Repeat("", 0) + string(rune('0'+t.seen[base]))
}

var _ = nonAlphanumRE // silence unused warning
