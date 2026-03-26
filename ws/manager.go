package ws

import (
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"
)

// ClientFactory creates a new BaseClient for a symbol and connects it.
type ClientFactory func(symbol string) (*BaseClient, error)

// BaseManager manages per-symbol workers and WebSocket clients.
// It provides the subscribe/unsubscribe API and handles worker/client lifecycle.
type BaseManager struct {
	mu            sync.RWMutex
	clients       map[string]*BaseClient
	workers       map[string]*SymbolWorker
	activeStreams map[string]map[StreamType]bool
	workerConfig  WorkerConfig
	clientFactory ClientFactory
	label         string
}

// NewBaseManager creates a new manager with the given config.
func NewBaseManager(label string, wcfg WorkerConfig, factory ClientFactory) *BaseManager {
	return &BaseManager{
		clients:       make(map[string]*BaseClient),
		workers:       make(map[string]*SymbolWorker),
		activeStreams: make(map[string]map[StreamType]bool),
		workerConfig:  wcfg,
		clientFactory: factory,
		label:         label,
	}
}

// GetOrCreateWorker returns an existing worker or creates a new one.
func (m *BaseManager) GetOrCreateWorker(symbol string) *SymbolWorker {
	m.mu.RLock()
	w, ok := m.workers[symbol]
	m.mu.RUnlock()
	if ok {
		return w
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	if w, ok := m.workers[symbol]; ok {
		return w
	}

	w = NewSymbolWorker(symbol, m.workerConfig)
	m.workers[symbol] = w
	return w
}

// GetOrCreateClient returns an existing client or creates, connects, and starts one.
func (m *BaseManager) GetOrCreateClient(symbol string) (*BaseClient, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if c, ok := m.clients[symbol]; ok {
		return c, nil
	}

	c, err := m.clientFactory(symbol)
	if err != nil {
		return nil, fmt.Errorf("create client %s: %w", symbol, err)
	}

	m.clients[symbol] = c
	if m.activeStreams[symbol] == nil {
		m.activeStreams[symbol] = make(map[StreamType]bool)
	}

	go c.ListenLoop()
	return c, nil
}

// ActivateStream marks a stream as active (idempotent).
func (m *BaseManager) ActivateStream(symbol string, st StreamType) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.activeStreams[symbol] == nil {
		m.activeStreams[symbol] = make(map[StreamType]bool)
	}
	m.activeStreams[symbol][st] = true
}

// DeactivateStream removes a stream. If no streams remain, the client and worker are stopped.
func (m *BaseManager) DeactivateStream(symbol string, st StreamType) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.activeStreams[symbol] == nil {
		return
	}
	delete(m.activeStreams[symbol], st)

	if len(m.activeStreams[symbol]) == 0 {
		if c, ok := m.clients[symbol]; ok {
			c.Stop()
			c.Close()
			delete(m.clients, symbol)
		}
		delete(m.activeStreams, symbol)

		if w, ok := m.workers[symbol]; ok {
			w.Stop()
			delete(m.workers, symbol)
		}
	}
}

// GetActiveStreams returns a copy of active streams for a symbol.
func (m *BaseManager) GetActiveStreams(symbol string) map[StreamType]bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make(map[StreamType]bool)
	for k, v := range m.activeStreams[symbol] {
		out[k] = v
	}
	return out
}

// GetSnapshot returns the snapshot for a symbol, or nil.
func (m *BaseManager) GetSnapshot(symbol string) *Snapshot {
	m.mu.RLock()
	w, ok := m.workers[symbol]
	m.mu.RUnlock()
	if !ok {
		return nil
	}
	return w.GetSnapshot()
}

// GetStatus returns a summary of active symbols, streams, and metrics.
func (m *BaseManager) GetStatus() map[string]any {
	m.mu.RLock()
	defer m.mu.RUnlock()

	symbols := make([]string, 0, len(m.clients))
	for s := range m.clients {
		symbols = append(symbols, s)
	}

	metrics := make(map[string]map[string]int64)
	for s, w := range m.workers {
		metrics[s] = w.GetMetrics()
	}

	return map[string]any{
		"label":   m.label,
		"symbols": symbols,
		"metrics": metrics,
	}
}

// Shutdown stops all clients and workers.
func (m *BaseManager) Shutdown() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, c := range m.clients {
		c.Stop()
		c.Close()
	}
	for _, w := range m.workers {
		w.Stop()
	}
	m.clients = make(map[string]*BaseClient)
	m.workers = make(map[string]*SymbolWorker)
	m.activeStreams = make(map[string]map[StreamType]bool)

	log.Infof("[%s] Manager shut down", m.label)
}
