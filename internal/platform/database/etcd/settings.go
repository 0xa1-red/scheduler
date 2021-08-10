package etcd

import "time"

var (
	etcdEndpoints   []string      = []string{"localhost:2379"}
	etcdDialTimeout time.Duration = 5 * time.Second
	etcdUsername    string        = "development"
)

func EtcdEndpoints() []string {
	return etcdEndpoints
}

func EtcdDialTimeout() time.Duration {
	return etcdDialTimeout
}

func EtcdUsername() string {
	return etcdUsername
}

func SetEtcdEndpoints(endpoints ...string) {
	if len(endpoints) > 0 {
		etcdEndpoints = endpoints
	}
}

func SetEtcdDialTimeout(duration string) {
	if timeout, err := time.ParseDuration(duration); err == nil {
		etcdDialTimeout = timeout
	}
}

func SetEtcdUsername(newUsername string) {
	etcdUsername = newUsername
}
