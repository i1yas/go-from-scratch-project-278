package application

// Service contains dependencies for application layer
type Service struct {
	links        linkStore
	shortcodeGen shortcodeGenerator
}

type shortcodeGenerator interface {
	Generate() string
}

// NewService takes dependencies and creates Service
func NewService(
	links linkStore,
	shortcodeGen shortcodeGenerator,
) *Service {
	return &Service{
		links:        links,
		shortcodeGen: shortcodeGen,
	}
}
