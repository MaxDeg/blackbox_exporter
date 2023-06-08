// Copyright 2016 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build windows
// +build windows

package main

import (
	"os"

	"github.com/go-kit/log/level"

	"github.com/prometheus/blackbox_exporter/exporter"
	"golang.org/x/sys/windows/svc"
)

type WindowsExporterService struct {
	stopCh chan<- bool
}

func main() {
	isService, err := svc.IsWindowsService()
	if err != nil {
		level.Error(e.Logger).Log("err", err)
		os.Exit(1)
	}

	stopCh := make(chan bool)
	if isService {
		go func() {
			err = svc.Run("Script Exporter", &WindowsExporterService{stopCh: ch})
			if err != nil {
				level.Error(e.Logger).Log("msg", "Failed to start service", "err", err)
				os.Exit(1)
			}
		}()
	}

	go func() {
		exporter.Run()
	}()

	for {
		if <-stopCh {
			level.Info(e.Logger).Log("msg", "Shutting down Script Exporter")
			break
		}
	}

}

func (s *WindowsExporterService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown
	changes <- svc.Status{State: svc.StartPending}
	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
loop:
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				s.stopCh <- true
				break loop
			default:
				log.Fatalf(fmt.Sprintf("unexpected control request #%d", c))
			}
		}
	}
	changes <- svc.Status{State: svc.StopPending}
	return
}
