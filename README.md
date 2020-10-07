# Youtube Subtitle Downloader

This is a simple youtube subtitle downloader. Any contribution is welcomed!

## Usage

```go
func main(){
    y, err := NewYoutubeDownloader(url)
	if err != nil {
		log.Fatal(err)
    }
    subtitle, err := y.Download()
    if err != nil {
        log.Fatal(err)
    }
    for _, line := range subtitle.Lines {
        // fmt.Println(line.Start, line.End, line.Text)
    }
}
```