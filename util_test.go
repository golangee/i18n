package i18n

import "testing"

func Test_guessLocaleFromFilename(t *testing.T) {
	tests := []struct {
		name string
		args string
		want string
	}{
		{"1", "bla-de-DE.xml", "de-DE"},
		{"2", "bla-de_DE.XMl", "de_DE"},
		{"3", "bla-en.xml", "en"},
		{"4", "bla-de-DE.toml", "de-DE"},
		{"4", "ignore-strings-de-DE_broken.xml", "de-DE"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := guessLocaleFromFilename(tt.args); got != tt.want {
				t.Errorf("guessLocaleFromFilename() = %v, want %v", got, tt.want)
			}
		})
	}
}
