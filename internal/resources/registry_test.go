package resources

import "testing"

func TestStringToBool(t *testing.T) {
	if !stringToBool("1") {
		t.Error("expected true for '1'")
	}
	if stringToBool("0") {
		t.Error("expected false for '0'")
	}
	if stringToBool("") {
		t.Error("expected false for empty string")
	}
}

func TestStringToInt64(t *testing.T) {
	cases := []struct {
		in   string
		want int64
	}{
		{"42", 42},
		{"0", 0},
		{"-1", -1},
		{"", 0},
		{"abc", 0},
		{" 7 ", 7},
	}
	for _, c := range cases {
		if got := stringToInt64(c.in); got != c.want {
			t.Errorf("stringToInt64(%q) = %d, want %d", c.in, got, c.want)
		}
	}
}

func TestStringToInt64Default(t *testing.T) {
	cases := []struct {
		in   string
		def  int64
		want int64
	}{
		{"", -1, -1},
		{"abc", -1, -1},
		{"0", -1, 0},
		{"42", -1, 42},
	}
	for _, c := range cases {
		if got := stringToInt64Default(c.in, c.def); got != c.want {
			t.Errorf("stringToInt64Default(%q, %d) = %d, want %d", c.in, c.def, got, c.want)
		}
	}
}

func TestStringToFloat64Default(t *testing.T) {
	cases := []struct {
		in   string
		def  float64
		want float64
	}{
		{"", -1, -1},
		{"abc", -1, -1},
		{"0", -1, 0},
		{"1.5", -1, 1.5},
	}
	for _, c := range cases {
		if got := stringToFloat64Default(c.in, c.def); got != c.want {
			t.Errorf("stringToFloat64Default(%q, %v) = %v, want %v", c.in, c.def, got, c.want)
		}
	}
}

func TestHclString(t *testing.T) {
	got := hclString(`foo "bar" \baz`)
	want := `"foo \"bar\" \\baz"`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestHclStringOrNull(t *testing.T) {
	if hclStringOrNull("") != "null" {
		t.Error("empty string should produce null")
	}
	if hclStringOrNull("x") != `"x"` {
		t.Error("non-empty string should be quoted")
	}
}

func TestHclBool(t *testing.T) {
	if hclBool(true) != "true" || hclBool(false) != "false" {
		t.Error("hclBool returned unexpected value")
	}
}

func TestHclSet(t *testing.T) {
	got := hclSet([]string{"b", "a", "", "c"})
	want := `["a", "b", "c"]`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestHclSetEmpty(t *testing.T) {
	if got := hclSet(nil); got != "[]" {
		t.Errorf("empty set got %q, want []", got)
	}
}

func TestHclInt(t *testing.T) {
	if hclInt(42) != "42" {
		t.Error("hclInt(42) should return \"42\"")
	}
}

func TestSplitNL(t *testing.T) {
	got := splitNL("a\nb\n\nc")
	if len(got) != 3 || got[0] != "a" || got[1] != "b" || got[2] != "c" {
		t.Errorf("unexpected splitNL result: %v", got)
	}
}

// pendingUpstreamTypes are resource types that exist in the
// vijayanms/terraform-provider-opnsense fork but haven't been merged into
// upstream browningluke/terraform-provider-opnsense yet — validTFTypes is
// generated from upstream only, so these need a hand-maintained exception
// until the corresponding PR lands. Remove an entry once 'task
// update-allowlist' picks it up from upstream on its own.
var pendingUpstreamTypes = map[string]bool{
	// vijayanms/terraform-provider-opnsense@aa00917: "Add opnsense_cron_job resource and data source"
	"opnsense_cron_job": true,
}

func TestRegistryTFTypesAreValid(t *testing.T) {
	for _, def := range Registry {
		if !validTFTypes[def.TFType] && !pendingUpstreamTypes[def.TFType] {
			t.Errorf("resource %q: TFType %q is not a real terraform-provider-opnsense resource type", def.Filename, def.TFType)
		}
	}
}

func TestRegistryNoDuplicates(t *testing.T) {
	seenType := map[string]bool{}
	seenFile := map[string]bool{}
	for _, def := range Registry {
		if seenType[def.TFType] {
			t.Errorf("duplicate TFType %q", def.TFType)
		}
		seenType[def.TFType] = true
		if seenFile[def.Filename] {
			t.Errorf("duplicate Filename %q", def.Filename)
		}
		seenFile[def.Filename] = true
	}
}

func TestRegistryEntriesComplete(t *testing.T) {
	for _, def := range Registry {
		if def.Filename == "" {
			t.Errorf("resource %q: empty Filename", def.TFType)
		}
		if def.Fetch == nil {
			t.Errorf("resource %q: nil Fetch func", def.TFType)
		}
	}
}
