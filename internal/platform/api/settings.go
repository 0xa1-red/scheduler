package api

var (
	address string = "0.0.0.0:8080"
)

func Address() string {
	return address
}

func SetAddress(newAddress string) {
	address = newAddress
}
