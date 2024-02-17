package cache

import (
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"sync"
	"time"

	"github.com/dmulholl/ironclad/internal/config"
	"github.com/dmulholl/ironclad/internal/ironcrypt"
	"github.com/dmulholl/ironclad/internal/ironcrypt/aes"
)

var CacheTimeout time.Duration = 15 * time.Minute

type IsCachedData struct {
	Filename string
}

type GetPassData struct {
	Filename  string
	CachePass string
	Token     string
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

	// If the token matches, it validates that the caller has read access to $HOME.
	token, found, err := config.Get("token")
	if err != nil {
		return fmt.Errorf("GetPass(): %v", err)
	}
	if !found {
		return fmt.Errorf("GetPass(): token not found")
	}
	if token != data.Token {
		return fmt.Errorf("GetPass(): invalid token")
	}

	// Do we have a cache entry for the specified database file?
	entry, found := server.cache[data.Filename]
	if !found {
		return fmt.Errorf("GetPass(): filename not in cache")
	}

	// Use the cache password and salt to regenerate the encryption key.
	key := ironcrypt.Key(data.CachePass, entry.salt, 10000, aes.KeySize)

	// Attempt to decrypt the entry.
	plaintext, err := aes.Decrypt(entry.data, key)
	if err != nil {
		if data.CachePass != "" {
			delete(server.cache, data.Filename)
		}
		return fmt.Errorf("invalid cache password")
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
		data: ciphertext,
	}
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
	err = config.Set("address", address)
	if err != nil {
		return errors.New("Serve: cannot set address")
	}

	rpc.Accept(listener)
	return nil
}
