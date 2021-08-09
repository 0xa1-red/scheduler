package database

var (
	path string = "/tmp/badger"
)

func Path() string {
	return path
}

func SetPath(newPath string) {
	path = newPath
}
