package nats

var (
	address        string = "127.0.0.1:4222"
	defaultSubject string = "updates"
)

// URL returns the configured URL
func Address() string {
	return address
}

// DefaultSubject returns the configured default subject
func DefaultSubject() string {
	return defaultSubject
}

// SetURL sets the new URL
func SetAddress(newAddress string) {
	address = newAddress
}

// SetDefaultSubject sets the default subject
func SetDefaultSubject(newSubject string) {
	defaultSubject = newSubject
}
