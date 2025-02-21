//go:build windows

// Package systray provides system tray functionality for displaying the current virtual desktop.
package systray

import (
	"context"
	"log/slog"
	"sync"
	"wincuts/config"
)

// Service manages the system tray icon and updates
type Service struct {
	icon    *Icon
	current int
	mu      sync.RWMutex
	ctx     context.Context
	cancel  context.CancelFunc
}

// NewService creates a new system tray service
func NewService(cfg config.TrayIconConfig) (*Service, error) {
	icon, err := New(cfg)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	svc := &Service{
		icon:   icon,
		ctx:    ctx,
		cancel: cancel,
	}

	// Set initial desktop number
	if err := svc.UpdateDesktop(1); err != nil {
		slog.Error("failed to set initial desktop number", "error", err)
	}

	return svc, nil
}

// UpdateDesktop updates the displayed desktop number
func (s *Service) UpdateDesktop(num int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.current == num {
		return nil
	}

	if err := s.icon.UpdateText(num); err != nil {
		return err
	}

	s.current = num
	return nil
}

// Stop cleans up resources and removes the system tray icon
func (s *Service) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.cancel()
	return s.icon.Close()
}
