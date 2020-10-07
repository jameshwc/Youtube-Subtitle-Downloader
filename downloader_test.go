package youtube_sub_dl

import (
	"reflect"
	"testing"
)

var (
	sampleLanguages = map[string]string{
		"Deutsch":            "de",
		"English":            "en",
		"Español":            "es",
		"Español (España)":   "es-ES",
		"Français":           "fr",
		"Hrvatski":           "hr",
		"Italiano":           "it",
		"Magyar":             "hu",
		"Polski":             "pl",
		"Português (Brasil)": "pt-BR",
		"Slovenčina":         "sk",
		"Tiếng Việt":         "vi",
		"Türkçe":             "tr",
		"default":            "en",
		"Русский":            "ru",
		"Српски":             "sr",
		"Українська":         "uk",
		"עברית":              "iw",
		"العربية":            "ar",
		"فارسی":              "fa",
		"বাংলা":              "bn",
		"ไทย":                "th",
		"မြန်မာ":             "my",
		"中文（简体）":             "zh-CN",
		"中文（繁體）":             "zh-TW",
		"日本語":                "ja",
	}
	sampleAndDefault = sampleLanguages
)

func Test_parseVideoID(t *testing.T) {
	type args struct {
		URL string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"http short url", args{URL: "http://youtu.be/5MgBikgcWnY"}, "5MgBikgcWnY", false},
		{"https short url", args{URL: "https://youtu.be/5MgBikgcWnY"}, "5MgBikgcWnY", false},
		{"http long url", args{URL: "http://www.youtube.com/watch?v=5MgBikgcWnY&feature=feed"}, "5MgBikgcWnY", false},
		{"https long url", args{URL: "http://www.youtube.com/watch?v=5MgBikgcWnY&feature=feed"}, "5MgBikgcWnY", false},
		{"no protocol url", args{URL: "www.youtube.com/watch?v=5MgBikgcWnY"}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseVideoID(tt.args.URL)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseVideoID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseVideoID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewYoutubeDownloader(t *testing.T) {
	type args struct {
		URL string
	}
	tests := []struct {
		name    string
		args    args
		want    *YoutubeDownloader
		wantErr bool
	}{
		{"http short url", args{URL: "http://youtu.be/5MgBikgcWnY"}, &YoutubeDownloader{"http://youtu.be/5MgBikgcWnY", "5MgBikgcWnY", sampleLanguages}, false},
		{"https short url", args{URL: "https://youtu.be/5MgBikgcWnY"}, &YoutubeDownloader{"https://youtu.be/5MgBikgcWnY", "5MgBikgcWnY", sampleLanguages}, false},
		{"http long url", args{URL: "http://www.youtube.com/watch?v=5MgBikgcWnY&feature=feed"}, &YoutubeDownloader{"http://www.youtube.com/watch?v=5MgBikgcWnY&feature=feed", "5MgBikgcWnY", sampleLanguages}, false},
		{"https long url", args{URL: "https://www.youtube.com/watch?v=5MgBikgcWnY&feature=feed"}, &YoutubeDownloader{"https://www.youtube.com/watch?v=5MgBikgcWnY&feature=feed", "5MgBikgcWnY", sampleLanguages}, false},
		{"no protocol url", args{URL: "www.youtube.com/watch?v=5MgBikgcWnY&feature=feed"}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewYoutubeDownloader(tt.args.URL)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewYoutubeDownloader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewYoutubeDownloader() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestYoutubeDownloader_getAvailableLanguages(t *testing.T) {
	type fields struct {
		URL       string
		VideoID   string
		Languages map[string]string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"http short url", fields{"http://youtu.be/5MgBikgcWnY", "5MgBikgcWnY", sampleLanguages}, false},
		{"https short url", fields{"https://youtu.be/5MgBikgcWnY", "5MgBikgcWnY", sampleLanguages}, false},
		{"http long url", fields{"http://www.youtube.com/watch?v=5MgBikgcWnY&feature=feed", "5MgBikgcWnY", sampleLanguages}, false},
		{"https long url", fields{"https://www.youtube.com/watch?v=5MgBikgcWnY&feature=feed", "5MgBikgcWnY", sampleLanguages}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			y := &YoutubeDownloader{
				URL:       tt.fields.URL,
				VideoID:   tt.fields.VideoID,
				Languages: tt.fields.Languages,
			}
			if err := y.getAvailableLanguages(); (err != nil) != tt.wantErr {
				t.Errorf("YoutubeDownloader.getAvailableLanguages() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
