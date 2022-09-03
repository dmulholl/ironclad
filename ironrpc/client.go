package ironrpc

import (
	"errors"
	"net"
	"net/rpc"
	"time"

	"github.com/dmulholl/ironclad/ironconfig"
)

var ClientTimeout = 100 * time.Millisecond

type CacheClient struct {
	rpcc *rpc.Client
}

// NewClient returns an initialized RPC client.
func NewClient() (*CacheClient, error) {
	address, found, err := ironconfig.Get("address")
	if err != nil {
		return nil, errors.New("NewClient: cannot read config file")
	}
	if !found {
		return nil, errors.New("NewClient: address not found")
	}

	conn, err := net.DialTimeout("tcp", address, ClientTimeout)
	if err != nil {
		return nil, err
	}

	return &CacheClient{rpcc: rpc.NewClient(conn)}, nil
}

// GetPass attempts to fetch a cached password from the server.
func (client *CacheClient) GetPass(filename, cachepass, token string) (string, error) {
	data := GetPassData{
		Filename:  filename,
		CachePass: cachepass,
		Token:     token,
	}
	var masterpass string
	err := client.rpcc.Call("CacheServer.GetPass", data, &masterpass)
	return masterpass, err
}

// SetPass attempts to cache a password to the server.
func (client *CacheClient) SetPass(filename, masterpass, cachepass string) error {
	data := SetPassData{
		Filename:   filename,
		MasterPass: masterpass,
		CachePass:  cachepass,
	}
	var notused bool
	return client.rpcc.Call("CacheServer.SetPass", data, &notused)
}

// IsCached checks if the server has a cache entry for the specified database file.
func (client *CacheClient) IsCached(filename string) bool {
	var found bool
	data := IsCachedData{Filename: filename}
	err := client.rpcc.Call("CacheServer.IsCached", data, &found)
	if err != nil {
		return false
	}
	return found
}

// Close calls the underlying net/rpc.Client's Close() method.
func (client *CacheClient) Close() {
	if client.rpcc != nil {
		client.rpcc.Close()
	}
}
