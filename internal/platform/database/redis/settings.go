package redis

var (
	address  string = "localhost:6379"
	password string = "test"
	database int    = 0
)

// Address returns the configured Redis address
func Address() string {
	return address
}

// Password returns the configured Redis password
func Password() string {
	return password
}

// Database returns the configured Redis database index
func Database() int {
	return database
}

// SetAddress sets the Redis address
func SetAddress(newAddress string) {
	address = newAddress
}

// SetPassword sets the Redis password
func SetPassword(newPassword string) {
	password = newPassword
}

// SetDatabase sets the Redis database index
func SetDatabase(newDB int) {
	database = newDB
}
