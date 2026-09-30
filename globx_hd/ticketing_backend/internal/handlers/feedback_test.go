package handlers

import "testing"

func TestValidateFeedbackImage(t *testing.T) {
	png := []byte("\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR")
	jpg := []byte("\xff\xd8\xff\xe0\x00\x10JFIF\x00")
	webp := []byte("RIFF\x00\x00\x00\x00WEBPVP8 ")
	svg := []byte(`<svg xmlns="http://www.w3.org/2000/svg"><script>alert(1)</script></svg>`)

	for name, tc := range map[string]struct {
		size    int64
		head    []byte
		wantExt string
	}{
		"png":            {1024, png, ".png"},
		"jpeg":           {1024, jpg, ".jpg"},
		"webp":           {1024, webp, ".webp"},
		"exactly 5MB":    {5 << 20, png, ".png"},
		"over 5MB":       {5<<20 + 1, png, ""},
		"svg rejected":   {100, svg, ""},
		"pdf rejected":   {100, []byte("%PDF-1.7"), ""},
		"text as .png":   {100, []byte("hello"), ""},
		"empty rejected": {0, nil, ""},
	} {
		ext, err := validateFeedbackImage(tc.size, tc.head)
		if tc.wantExt == "" && err == nil {
			t.Errorf("%s: expected rejection, got %q", name, ext)
		}
		if tc.wantExt != "" && (err != nil || ext != tc.wantExt) {
			t.Errorf("%s: want %q, got %q err=%v", name, tc.wantExt, ext, err)
		}
	}
}
