package main

import (
	"sync"
	"time"
)

type XRayCoreLogsHandler struct {
	core     *XRayCore
	callback func(string)
	interval time.Duration
	active   bool
	mu       sync.Mutex
	wg       sync.WaitGroup
}

func NewXRayCoreLogsHandler(core *XRayCore, callback func(string), interval time.Duration) *XRayCoreLogsHandler {
	handler := &XRayCoreLogsHandler{
		core:     core,
		callback: callback,
		interval: interval,
		active:   true,
	}
	handler.wg.Add(1)
	go handler.cast()
	return handler
}

func (h *XRayCoreLogsHandler) Stop() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.active = false
	h.wg.Wait()
}

func (h *XRayCoreLogsHandler) cast() {
	defer h.wg.Done()
	ticker := time.NewTicker(h.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			h.mu.Lock()
			if !h.active {
				h.mu.Unlock()
				return
			}
			// Replace with actual log fetching logic
			log := "Sample log message"
			h.callback(log)
			h.mu.Unlock()
		}
	}
}
