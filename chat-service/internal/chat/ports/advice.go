package ports

type AdvicePublisher interface {
	PublishAdviceRequest(threadID string) error
}
