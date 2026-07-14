package label

import (
	"strings"
	"testing"
)

func TestSanitizeLabel(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"My Alias", "my_alias"},
		{"foo-bar.baz", "foo_bar_baz"},
		{"  leading", "leading"},
		{"UPPER", "upper"},
		{"__double__under__", "double_under"},
		{"", "resource"},
		{strings.Repeat("a", 70), strings.Repeat("a", 63)},
	}
	for _, c := range cases {
		got := sanitizeLabel(c.in)
		if got != c.want {
			t.Errorf("sanitizeLabel(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestLabelTrackerDerive(t *testing.T) {
	tr := New()

	// name takes priority
	if got := tr.Derive("MyAlias", "desc", "uuid-1234"); got != "myalias" {
		t.Errorf("got %q, want myalias", got)
	}

	// collision appends suffix
	got2 := tr.Derive("MyAlias", "desc", "uuid-5678")
	if !strings.HasPrefix(got2, "myalias_") {
		t.Errorf("expected collision suffix, got %q", got2)
	}

	// no name — falls back to description
	tr2 := New()
	if got := tr2.Derive("", "my description", "uuid-0000"); got != "my_description" {
		t.Errorf("got %q, want my_description", got)
	}

	// no name, no description — falls back to uuid prefix
	tr3 := New()
	if got := tr3.Derive("", "", "abcdef12-rest"); got != "abcdef12" {
		t.Errorf("got %q, want abcdef12", got)
	}
}
