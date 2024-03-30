package gapi

import (
	"fmt"
	"sync"

	"github.com/Streamfair/streamfair_idp/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type PoolConfig struct {
	MaxOpenConnection     int
	MaxIdleConnection     int
	ConnectionQueueLength int
	Address               string
	ConfigOptions         []grpc.DialOption
}

type ConnectionPool struct {
	mu                  sync.Mutex
	address             string
	configOptions       []grpc.DialOption
	maxOpenConnection   int
	maxIdleConnection   int
	numOfOpenConnection int
	connectionQueue     chan *grpc.ClientConn
	idleConnections     map[string]*grpc.ClientConn
}

func NewClientPool(config *PoolConfig) *ConnectionPool {
	clientPool := &ConnectionPool{
		address:             config.Address,
		configOptions:       config.ConfigOptions,
		maxOpenConnection:   config.MaxOpenConnection,
		maxIdleConnection:   config.MaxIdleConnection,
		numOfOpenConnection: 0,
		connectionQueue:     make(chan *grpc.ClientConn, config.ConnectionQueueLength),
		idleConnections:     make(map[string]*grpc.ClientConn),
	}
	go clientPool.handleConnectionQueue()
	return clientPool
}

func (cp *ConnectionPool) handleConnectionQueue() {
	for {
		select {
		case conn := <-cp.connectionQueue:
			cp.mu.Lock()
			if cp.numOfOpenConnection > cp.maxOpenConnection {
				cp.numOfOpenConnection--
				conn.Close()
			} else {
				cp.idleConnections[cp.address] = conn
			}
			cp.mu.Unlock()
		}
	}
}

func (cp *ConnectionPool) GetConn(address string) (*grpc.ClientConn, error) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	if conn, ok := cp.idleConnections[address]; ok {
		delete(cp.idleConnections, address)
		return conn, nil
	}

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

	cp.numOfOpenConnection++
	return conn, nil
}

func (cp *ConnectionPool) ReleaseConn(conn *grpc.ClientConn) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	if cp.numOfOpenConnection <= cp.maxIdleConnection {
		cp.idleConnections[cp.address] = conn
	} else {
		cp.numOfOpenConnection--
		conn.Close()
	}
}
