package shortcut

import (
	"log/slog"
	"sync"
)

const (
	blockChanBufferSize = 10 // Increased buffer size
)

// Service struct
type Service struct {
	shortcutChan chan *KeyBindingAction
	matcher      *Matcher
	wg           sync.WaitGroup
	stopChan     chan struct{}
	blockChan    chan bool // Channel to communicate blocking decisions
}

// NewService creates a new KeybindingService
func NewService(shortcutChan chan *KeyBindingAction, matcher *Matcher) *Service {
	return &Service{
		shortcutChan: shortcutChan,
		matcher:      matcher,
		stopChan:     make(chan struct{}),
		blockChan:    make(chan bool, blockChanBufferSize), // Larger buffer

	}
}

// RegisterKeyBindingActions registers key binding actions
func (s *Service) RegisterKeyBindingActions(bindings ...KeyBindingAction) *Service {
	s.matcher.AddBindings(bindings...)
	return s
}

// Start starts the keybinding service to listen for events
func (s *Service) Start() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for {
			select {
			case <-s.stopChan:
				return
			case binding := <-s.shortcutChan:
				go func() {
					slog.Info("executing action", "binding", binding.Binding.PrettyString())
					if err := binding.Execute(); err != nil {
						slog.Error("failed to execute action", "error", err)
					}
				}()
			}
		}
	}()
}

// GetBlockDecision returns the blocking decision
func (s *Service) GetBlockDecision() bool {
	select {
	case shouldBlock := <-s.blockChan:
		return shouldBlock
	default:
		return false
	}
}

func (s *Service) GetShortcutChan() chan *KeyBindingAction {
	return s.shortcutChan
}

// Stop stops the keybinding service
func (s *Service) Stop() {
	close(s.stopChan)
	s.wg.Wait()
}


func (s *Service) Match(event KeyEvent) (*KeyBindingAction, bool) {
	return s.matcher.Match(event)
}

