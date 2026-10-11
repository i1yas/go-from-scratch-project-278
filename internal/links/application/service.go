package application

import "log/slog"

// Service contains dependencies for application layer
type Service struct {
	links        linkStore
	visits       visitsStore
	shortcodeGen shortcodeGenerator
	logger       *slog.Logger
}

type shortcodeGenerator interface {
	Generate() (string, error)
}

// NewService takes dependencies and creates Service
func NewService(
	links linkStore,
	visits visitsStore,
	shortcodeGen shortcodeGenerator,
	logger *slog.Logger,
) *Service {
	return &Service{
		links:        links,
		visits:       visits,
		shortcodeGen: shortcodeGen,
		logger:       logger,
	}
}
