package youtube_sub_dl

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/beevik/etree"
)

type YoutubeDownloader struct {
	URL       string
	VideoID   string
	Languages map[string]string
}

type Line struct {
	Start time.Duration
	End   time.Duration
	Text  string
	Index int // zero-based
}

type Subtitle struct {
	Lines []Line
}

func NewYoutubeDownloader(URL string) (*YoutubeDownloader, error) {
	y := new(YoutubeDownloader)
	y.URL = URL
	var err error
	y.VideoID, err = parseVideoID(URL)
	if err != nil {
		return nil, err
	}
	y.Languages = make(map[string]string)
	if err := y.getAvailableLanguages(); err != nil {
		return nil, err
	}
	return y, nil
}

// Examples:
// - http://youtu.be/5MgBikgcWnY
// - http://www.youtube.com/watch?v=5MgBikgcWnY&feature=feed
// - http://www.youtube.com/embed/5MgBikgcWnY
// - http://www.youtube.com/v/5MgBikgcWnY?version=3&amp;hl=en_US
func parseVideoID(URL string) (string, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return "", err
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return "", fmt.Errorf("not supported scheme (http/https)")
	}

	switch u.Host {

	case "youtu.be":
		return u.Path[1:], nil

	case "youtube.com", "www.youtube.com":

		switch {

		case strings.HasPrefix(u.Path, "/watch"):
			return u.Query().Get("v"), nil

		case strings.HasPrefix(u.Path, "/v/"), strings.HasPrefix(u.Path, "/embed/"):
			return strings.Split(u.Path, "/")[2], nil

		default:
			return "", fmt.Errorf("path not correct")

		}

	default:
		return "", fmt.Errorf("host not correct")
	}
}

func (y *YoutubeDownloader) getAvailableLanguages() error {
	URL := fmt.Sprintf("http://www.youtube.com/api/timedtext?v=%s&type=list", y.VideoID)

	resp, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	doc := etree.NewDocument()
	if _, err := doc.ReadFrom(resp.Body); err != nil {
		return err
	}

	root := doc.SelectElement("transcript_list")
	for _, s := range root.SelectElements("track") {
		y.Languages[s.SelectAttr("lang_original").Value] = s.SelectAttr("lang_code").Value
		if s.SelectAttrValue("lang_default", "false") == "true" {
			y.Languages["default"] = s.SelectAttr("lang_code").Value
		}
	}

	return nil
}

func (y *YoutubeDownloader) Download() (*Subtitle, error) {
	doc, err := y.download()
	if err != nil {
		return nil, err
	}
	var subtitle Subtitle
	idx := 0
	for _, s := range doc.SelectElements("text") {
		start, err := strconv.ParseFloat(s.SelectAttr("start").Value, 64)
		if err != nil {
			return nil, err
		}
		dur, err := strconv.ParseFloat(s.SelectAttr("dur").Value, 64)
		if err != nil {
			return nil, err
		}
		var l Line
		l.Start = time.Duration(start * float64(time.Second))
		l.End = time.Duration((start + dur) * float64(time.Second))
		l.Text = strings.ReplaceAll(s.Text(), "&#39;", "'")
		l.Index = idx
		idx++
		subtitle.Lines = append(subtitle.Lines, l)
	}
	return &subtitle, nil
}

func (y *YoutubeDownloader) download() (*etree.Element, error) {
	downloadUrl := fmt.Sprintf("http://www.youtube.com/api/timedtext?v=%s&lang=%s", y.VideoID, y.Languages["default"])

	fmt.Println(downloadUrl)
	resp, err := http.Get(downloadUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc := etree.NewDocument()
	if _, err := doc.ReadFrom(resp.Body); err != nil {
		return nil, err
	}
	return doc.SelectElement("transcript"), nil

}
