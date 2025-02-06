package service

import (
	"fmt"
	"log/slog"
	"time"
	"wincuts/app"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"
)

const (
	serviceName        = "WinCuts"
	serviceDisplayName = "WinCuts Virtual Desktop Manager"
	serviceDescription = "Manages virtual desktops and provides keyboard shortcuts for desktop switching"
)

type winCutsService struct {
	stopChan chan struct{}
}

// Execute implements the Windows service interface
func (s *winCutsService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown
	changes <- svc.Status{State: svc.StartPending}

	s.stopChan = make(chan struct{})

	// Start the application in a goroutine
	go func() {
		if err := app.Run(); err != nil {
			slog.Error("application error", "error", err)
		}
	}()

	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

	// Wait for stop signal
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				changes <- svc.Status{State: svc.StopPending}
				close(s.stopChan)
				return false, 0
			default:
				slog.Error("unexpected control request", "command", c)
			}
		}
	}
}

// RunAsService starts WinCuts as a Windows service
func RunAsService(isDebug bool) error {
	var err error
	if !isDebug {
		// Set up event logging
		elog, err := eventlog.Open(serviceName)
		if err != nil {
			return fmt.Errorf("failed to open event log: %w", err)
		}
		defer elog.Close()

		elog.Info(1, fmt.Sprintf("starting %s service", serviceName))
	}

	run := svc.Run
	if isDebug {
		run = debug.Run
	}

	err = run(serviceName, &winCutsService{})
	if err != nil {
		return fmt.Errorf("service failed: %w", err)
	}

	return nil
}

// IsWindowsService determines if the program is running as a Windows service
func IsWindowsService() (bool, error) {
	return svc.IsWindowsService()
}

// InstallService installs WinCuts as a Windows service
func InstallService(exePath string) error {
	m, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to service manager: %w", err)
	}
	defer m.Disconnect()

	s, err := m.OpenService(serviceName)
	if err == nil {
		s.Close()
		return fmt.Errorf("service %s already exists", serviceName)
	}

	config := mgr.Config{
		DisplayName: serviceDisplayName,
		Description: serviceDescription,
		StartType:   mgr.StartAutomatic,
	}

	s, err = m.CreateService(serviceName, exePath, config, "-service")
	if err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}
	defer s.Close()

	// Set up event logging
	err = eventlog.InstallAsEventCreate(serviceName, eventlog.Error|eventlog.Warning|eventlog.Info)
	if err != nil {
		s.Delete()
		return fmt.Errorf("failed to set up event logging: %w", err)
	}

	return nil
}

// UninstallService removes the WinCuts Windows service
func UninstallService() error {
	m, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to service manager: %w", err)
	}
	defer m.Disconnect()

	s, err := m.OpenService(serviceName)
	if err != nil {
		return fmt.Errorf("service %s is not installed", serviceName)
	}
	defer s.Close()

	// Stop the service if it's running
	status, err := s.Query()
	if err != nil {
		return fmt.Errorf("failed to query service status: %w", err)
	}

	if status.State != svc.Stopped {
		// Request service stop
		_, err := s.Control(svc.Stop)
		if err != nil {
			return fmt.Errorf("failed to stop service: %w", err)
		}

		// Wait for the service to stop
		timeout := time.Now().Add(10 * time.Second)
		for status.State != svc.Stopped {
			if time.Now().After(timeout) {
				return fmt.Errorf("timed out waiting for service to stop")
			}
			time.Sleep(300 * time.Millisecond)
			status, err = s.Query()
			if err != nil {
				return fmt.Errorf("failed to query service status: %w", err)
			}
		}
	}

	// Remove the service
	err = s.Delete()
	if err != nil {
		return fmt.Errorf("failed to delete service: %w", err)
	}

	// Remove event logging
	err = eventlog.Remove(serviceName)
	if err != nil {
		return fmt.Errorf("failed to remove event logging: %w", err)
	}

	return nil
}
