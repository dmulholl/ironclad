package ironrpc

import (
	"errors"
	"net"
	"net/rpc"
	"os"
	"sync"
	"time"

	"github.com/dmulholl/ironclad/ironconfig"
	"github.com/dmulholl/ironclad/ironcrypt"
	"github.com/dmulholl/ironclad/ironcrypt/aes"
)

var CacheTimeout time.Duration = 15 * time.Minute

type IsCachedData struct {
	Filename string
}

type GetPassData struct {
	Filename  string
	CachePass string
}

type SetPassData struct {
	Filename   string
	MasterPass string
	CachePass  string
}

type CacheEntry struct {
	salt []byte
	data []byte
}

type CacheServer struct {
	cache      map[string]CacheEntry
	mutex      *sync.Mutex
	lastaccess time.Time
}

// newServer returns an initialized RPC password server.
func newServer() *CacheServer {
	return &CacheServer{
		cache:      make(map[string]CacheEntry),
		mutex:      &sync.Mutex{},
		lastaccess: time.Now(),
	}
}

// Contains method exposed by the RPC server.
func (server *CacheServer) IsCached(data IsCachedData, result *bool) error {
	server.mutex.Lock()
	defer server.mutex.Unlock()

	if _, found := server.cache[data.Filename]; found {
		*result = true
		return nil
	}

	*result = false
	return nil
}

// GetPass method exposed by the RPC server.
func (server *CacheServer) GetPass(data GetPassData, password *string) error {
	server.mutex.Lock()
	defer server.mutex.Unlock()

	// Do we have a cache entry for the specified database file?
	entry, found := server.cache[data.Filename]
	if !found {
		return errors.New("GetPass: filename not in cache")
	}

	// Use the cache password and salt to regenerate the encryption key.
	key := ironcrypt.Key(data.CachePass, entry.salt, 10000, aes.KeySize)

	// Attempt to decrypt the entry. Delete the entry from the cache on failure.
	plaintext, err := aes.Decrypt(entry.data, key)
	if err != nil {
		delete(server.cache, data.Filename)
		return errors.New("GetPass: decryption failure")
	}

	*password = string(plaintext)
	server.lastaccess = time.Now()
	return nil
}

// SetPass method exposed by the RPC server.
func (server *CacheServer) SetPass(data SetPassData, notused *bool) error {
	server.mutex.Lock()
	defer server.mutex.Unlock()

	// Generate a random salt.
	salt, err := ironcrypt.RandBytes(32)
	if err != nil {
		return errors.New("SetPass: cannot generate random salt")
	}

	// Generate an encryption key from the cache password.
	key := ironcrypt.Key(data.CachePass, salt, 10000, aes.KeySize)

	// Encrypt the database password using the cache password.
	ciphertext, err := aes.Encrypt([]byte(data.MasterPass), key)
	if err != nil {
		return errors.New("Set Pass: cannot encrypt master password")
	}

	server.cache[data.Filename] = CacheEntry{
		salt: salt,
		data: ciphertext}
	return nil
}

// Automatically shuts the server down after the specified duration.
func (server *CacheServer) timeout() {
	for {
		server.mutex.Lock()
		if time.Since(server.lastaccess) > CacheTimeout {
			os.Exit(0)
		}
		server.mutex.Unlock()
		time.Sleep(time.Second)
	}
}

// Serve launches a new RPC password server and blocks until the server
// automatically shuts itself down via the timeout() function.
func Serve() error {
	server := newServer()
	go server.timeout()

	err := rpc.Register(server)
	if err != nil {
		return err
	}

	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return err
	}
	defer listener.Close()

	address := listener.Addr().String()
	err = ironconfig.Set("address", address)
	if err != nil {
		return errors.New("Serve: cannot set address")
	}

	rpc.Accept(listener)
	return nil
}
