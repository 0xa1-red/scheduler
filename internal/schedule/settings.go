package schedule

var (
	defaultTopic string = "updates"
)

func DefaultTopic() string {
	return defaultTopic
}

func SetDefaultTopic(newTopic string) {
	defaultTopic = newTopic
}
