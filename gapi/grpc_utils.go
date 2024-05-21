package gapi

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Streamfair/streamfair_idp/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// PoolConfig holds the configuration for the gRPC connection pool.
type PoolConfig struct {
	MaxOpenConnection     int               // Maximum number of open connections allowed.
	MaxIdleConnection     int               // Maximum number of idle connections to keep in the pool.
	ConnectionQueueLength int               // Length of the queue for managing connections.
	Address               string            // Address of the gRPC service.
	ConfigOptions         []grpc.DialOption // Additional gRPC dial options.
	IdleTimeout           time.Duration     // Duration after which idle connections are closed.
}

// ConnectionPool manages the gRPC connections.
type ConnectionPool struct {
	mu                        sync.Mutex
	configOptions             []grpc.DialOption
	maxOpenConnection         int
	maxIdleConnection         int
	numOfOpenConnection       int
	connectionQueue           chan *grpc.ClientConn
	idleConnections           map[string]map[*grpc.ClientConn]struct{}
	lastUsed                  map[*grpc.ClientConn]time.Time
	atomicNumOfOpenConnection uint32 // Use atomic type for safe updates
}

// NewClientPool creates a new connection pool with the given configuration.
func NewClientPool(config *PoolConfig) *ConnectionPool {
	clientPool := &ConnectionPool{
		configOptions:       config.ConfigOptions,
		maxOpenConnection:   config.MaxOpenConnection,
		maxIdleConnection:   config.MaxIdleConnection,
		numOfOpenConnection: 0,
		connectionQueue:     make(chan *grpc.ClientConn, config.ConnectionQueueLength),
		idleConnections:     make(map[string]map[*grpc.ClientConn]struct{}), // Initialize as map of maps.
		lastUsed:            make(map[*grpc.ClientConn]time.Time),
	}
	go clientPool.handleConnectionQueue()                    // Start goroutine to manage connection queue.
	go clientPool.cleanupIdleConnections(config.IdleTimeout) // Start goroutine to cleanup idle connections.
	return clientPool
}

// handleConnectionQueue is a goroutine that manages connections returned to the pool.
func (cp *ConnectionPool) handleConnectionQueue() {
	for {
		select {
		case conn := <-cp.connectionQueue:
			cp.mu.Lock()
			address := conn.Target() // Get the address of the connection.
			if _, ok := cp.idleConnections[address]; !ok {
				cp.idleConnections[address] = make(map[*grpc.ClientConn]struct{})
			}
			if atomic.LoadUint32(&cp.atomicNumOfOpenConnection) > uint32(cp.maxOpenConnection) {
				atomic.AddUint32(&cp.atomicNumOfOpenConnection, ^uint32(0))
				conn.Close() // Close the connection if the limit is exceeded.
			} else {
				cp.idleConnections[address][conn] = struct{}{} // Add the connection to the idle connections map.
			}
			cp.mu.Unlock()
		}
	}
}

// cleanupIdleConnections is a goroutine that periodically checks and closes idle connections.
func (cp *ConnectionPool) cleanupIdleConnections(timeout time.Duration) {
	for {
		time.Sleep(timeout)
		cp.mu.Lock()
		for address, conns := range cp.idleConnections {
			for conn := range conns {
				if time.Since(cp.lastUsed[conn]) > timeout {
					conn.Close() // Close the connection if it's idle for longer than the timeout.
					delete(cp.idleConnections[address], conn)
					delete(cp.lastUsed, conn)
					cp.numOfOpenConnection--
				}
			}
		}
		cp.mu.Unlock()
	}
}

// GetConn obtains a connection from the pool. If no idle connections are available, it creates a new connection.
func (cp *ConnectionPool) GetConn(address string) (*grpc.ClientConn, error) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	if conns, ok := cp.idleConnections[address]; ok && len(conns) > 0 {
		for conn := range conns {
			delete(cp.idleConnections[address], conn)
			cp.lastUsed[conn] = time.Now() // Update the last used time of the connection.
			return conn, nil
		}
	}

	// Load TLS configuration and dial the gRPC service.
	config := util.GetConfigService().GetConfig()
	tlsConfig, err := LoadTLSConfigWithTrustedCerts(config.CertPem, config.KeyPem, config.CaCertPem)
	if err != nil {
		return nil, fmt.Errorf("failed to load TLS config for 'GetConn': %w", err)
	}

	creds := credentials.NewTLS(tlsConfig)
	conn, err := grpc.Dial(address, append(cp.configOptions, grpc.WithTransportCredentials(creds))...)
	if err != nil {
		return nil, err
	}

	atomic.AddUint32(&cp.atomicNumOfOpenConnection, 1) // Safely update the number of open connections

	return conn, nil
}

// ReleaseConn releases a connection back to the pool. If the pool is full, it closes the connection.
func (cp *ConnectionPool) ReleaseConn(conn *grpc.ClientConn) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	address := conn.Target() // Get the address of the connection.
	if _, ok := cp.idleConnections[address]; !ok {
		cp.idleConnections[address] = make(map[*grpc.ClientConn]struct{})
	}

	if len(cp.idleConnections[address]) < cp.maxIdleConnection {
		cp.idleConnections[address][conn] = struct{}{} // Add the connection to the idle connections map.
		cp.lastUsed[conn] = time.Now()                 // Update the last used time of the connection.
	} else {
		cp.numOfOpenConnection-- // Decrement the number of open connections.
		conn.Close()             // Close the connection if the pool is full.
	}
}
