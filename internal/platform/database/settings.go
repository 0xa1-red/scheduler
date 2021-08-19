package database

var (
	backend DBKind = KindRedis
)

func Backend() DBKind {
	return backend
}

func SetBackend(newBackend string) {
	kind := DBKind(newBackend)
	switch kind {
	case KindRedis:
		backend = kind
	}
}
