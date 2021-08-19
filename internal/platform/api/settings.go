package api

var (
	address string = "0.0.0.0:8080"
)

// Address returns the configured address
func Address() string {
	return address
}

// SetAddress sets a new address
func SetAddress(newAddress string) {
	address = newAddress
}
