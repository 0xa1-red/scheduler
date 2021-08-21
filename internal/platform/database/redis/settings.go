package redis

import "time"

var (
	address  string        = "localhost:6379"
	password string        = "test"
	database int           = 0
	retries  int           = 3
	interval time.Duration = 1 * time.Second
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

// Retries returns the configured number of retries
func Retries() int {
	return retries
}

// Interval returns the configured delay between retries
func Interval() time.Duration {
	return interval
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

// SetRetries sets the new number of retries
func SetRetries(newRetries int) {
	if newRetries == 0 {
		newRetries++
	}
	retries = newRetries
}

// SetInterval sets the new delay between retries
func SetInterval(newInterval string) {
	if d, err := time.ParseDuration(newInterval); err == nil {
		interval = d
	}
}
