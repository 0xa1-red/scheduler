package nats

var (
	url string = "127.0.0.1:4222"
)

// URL returns the configured URL
func URL() string {
	return url
}

// SetURL sets the new URL
func SetURL(newURL string) {
	url = newURL
}
