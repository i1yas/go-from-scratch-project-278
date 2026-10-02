package application

// Service contains dependencies for application layer
type Service struct {
	links        linkStore
	visits       visitsStore
	shortcodeGen shortcodeGenerator
}

type shortcodeGenerator interface {
	Generate() (string, error)
}

// NewService takes dependencies and creates Service
func NewService(
	links linkStore,
	visits visitsStore,
	shortcodeGen shortcodeGenerator,
) *Service {
	return &Service{
		links:        links,
		visits:       visits,
		shortcodeGen: shortcodeGen,
	}
}
