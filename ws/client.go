package ws

import (
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

// ClientConfig holds configuration for a WebSocket client.
type ClientConfig struct {
	// BaseURL is the WebSocket endpoint (e.g., "wss://fstream.binance.com/ws").
	BaseURL string

	// PingInterval controls heartbeat frequency. 0 disables pings.
	PingInterval time.Duration

	// PingMessage is the raw ping payload. If nil, pings are disabled.
	PingMessage func() ([]byte, error)

	// ReadBufferSize and WriteBufferSize for the WebSocket dialer.
	ReadBufferSize  int
	WriteBufferSize int

	// ReconnectDelay is the pause before reconnecting after a read error.
	ReconnectDelay time.Duration

	// Label is used in log messages (e.g., "Binance").
	Label string
}

// DefaultClientConfig returns sensible defaults.
func DefaultClientConfig(label, baseURL string) ClientConfig {
	return ClientConfig{
		BaseURL:         baseURL,
		ReadBufferSize:  16 * 1024,
		WriteBufferSize: 16 * 1024,
		ReconnectDelay:  2 * time.Second,
		Label:           label,
	}
}

// DispatchFunc is called for each incoming WebSocket message.
type DispatchFunc func(msg []byte)

// ResubscribeFunc is called after a reconnect to restore stream subscriptions.
type ResubscribeFunc func(client *BaseClient) error

// BaseClient manages a single WebSocket connection with automatic reconnection,
// optional heartbeat, and message dispatch.
type BaseClient struct {
	mu     sync.Mutex
	conn   *websocket.Conn
	config ClientConfig
	symbol string

	stopChan chan struct{}
	stopped  bool

	dispatch    DispatchFunc
	resubscribe ResubscribeFunc
}

// NewBaseClient creates a new WebSocket client.
func NewBaseClient(symbol string, cfg ClientConfig) *BaseClient {
	return &BaseClient{
		symbol:   symbol,
		config:   cfg,
		stopChan: make(chan struct{}),
	}
}

// SetDispatch sets the message handler.
func (c *BaseClient) SetDispatch(fn DispatchFunc) { c.dispatch = fn }

// SetResubscribe sets the function called after reconnection.
func (c *BaseClient) SetResubscribe(fn ResubscribeFunc) { c.resubscribe = fn }

// Symbol returns the symbol this client is connected for.
func (c *BaseClient) Symbol() string { return c.symbol }

// Connect establishes the WebSocket connection.
func (c *BaseClient) Connect() error {
	dialer := *websocket.DefaultDialer
	dialer.EnableCompression = true
	dialer.ReadBufferSize = c.config.ReadBufferSize
	dialer.WriteBufferSize = c.config.WriteBufferSize

	conn, _, err := dialer.Dial(c.config.BaseURL, nil)
	if err != nil {
		return fmt.Errorf("connect %s: %w", c.config.Label, err)
	}
	c.mu.Lock()
	c.conn = conn
	c.mu.Unlock()
	return nil
}

// Close closes the underlying connection.
func (c *BaseClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

// Stop signals the listen loop to exit.
func (c *BaseClient) Stop() {
	select {
	case <-c.stopChan:
	default:
		close(c.stopChan)
	}
}

// WriteMessage sends a text message on the WebSocket.
func (c *BaseClient) WriteMessage(data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn == nil {
		return fmt.Errorf("not connected")
	}
	return c.conn.WriteMessage(websocket.TextMessage, data)
}

// ListenLoop reads messages and dispatches them. Blocks until Stop is called
// or an unrecoverable error occurs. Automatically reconnects on read errors.
func (c *BaseClient) ListenLoop() {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("[%s:%s] listen.panic: %v", c.config.Label, c.symbol, r)
			c.Stop()
			c.Close()
		}
	}()

	// Start heartbeat if configured
	if c.config.PingInterval > 0 && c.config.PingMessage != nil {
		ticker := time.NewTicker(c.config.PingInterval)
		defer ticker.Stop()
		go func() {
			for {
				select {
				case <-c.stopChan:
					return
				case <-ticker.C:
					msg, err := c.config.PingMessage()
					if err != nil {
						log.Warnf("[%s:%s] ping.marshal: %v", c.config.Label, c.symbol, err)
						continue
					}
					if err := c.WriteMessage(msg); err != nil {
						log.Warnf("[%s:%s] ping.send: %v", c.config.Label, c.symbol, err)
					}
				}
			}
		}()
	}

	for {
		select {
		case <-c.stopChan:
			log.Infof("[%s:%s] listen.stop", c.config.Label, c.symbol)
			c.Close()
			return
		default:
			c.mu.Lock()
			conn := c.conn
			c.mu.Unlock()

			if conn == nil {
				time.Sleep(100 * time.Millisecond)
				continue
			}

			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Warnf("[%s:%s] read.error: %v", c.config.Label, c.symbol, err)
				time.Sleep(c.config.ReconnectDelay)
				if rErr := c.reconnect(); rErr != nil {
					log.Errorf("[%s:%s] reconnect.failed: %v", c.config.Label, c.symbol, rErr)
				}
				continue
			}

			if c.dispatch != nil {
				c.dispatch(msg)
			}
		}
	}
}

func (c *BaseClient) reconnect() error {
	c.Close()
	time.Sleep(1 * time.Second)

	if err := c.Connect(); err != nil {
		return err
	}

	if c.resubscribe != nil {
		return c.resubscribe(c)
	}
	return nil
}
