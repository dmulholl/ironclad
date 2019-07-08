package ironrpc


import (
    "encoding/base64"
    "net"
    "net/rpc"
    "time"
    "sync"
    "errors"
    "os"
)


import (
    "github.com/dmulholl/ironclad/ironconfig"
    "github.com/dmulholl/ironclad/ironcrypt"
)


var CacheTimeout = 15 * time.Minute


type GetData struct {
    Filename string
    Nonce string
}


type SetData struct {
    Filename string
    Password string
}


type CacheServer struct {
    cache map[string]string
    mutex *sync.Mutex
    lastaccess time.Time
}


// NewServer returns an initialized RPC password server.
func NewServer() *CacheServer {
    return &CacheServer{
        cache: make(map[string]string),
        mutex: &sync.Mutex{},
        lastaccess: time.Now(),
    }
}


// GetPass method exposed by the RPC server.
func (server *CacheServer) GetPass(data GetData, password *string) error {
    server.mutex.Lock()
    defer server.mutex.Unlock()

    nonce, found, err := ironconfig.Get("nonce")
    if err != nil {
        return errors.New("GetPass: cannot read config file")
    }
    if !found {
        return errors.New("GetPass: nonce not found in file")
    }
    if nonce != data.Nonce {
        time.Sleep(time.Second)
        return errors.New("GetPass: invalid nonce")
    }

    cachedpass, found := server.cache[data.Filename]
    if found {
        *password = cachedpass
        server.lastaccess = time.Now()
        return nil
    } else {
        return errors.New("GetPass: filename not in cache")
    }
}


// SetPass method exposed by the RPC server.
func (server *CacheServer) SetPass(data SetData, notused *bool) error {
    server.mutex.Lock()
    defer server.mutex.Unlock()

    bytes, err := ironcrypt.RandBytes(32)
    if err != nil {
        return errors.New("SetPass: cannot generate random bytes")
    }

    nonce := base64.StdEncoding.EncodeToString(bytes)
    err = ironconfig.Set("nonce", nonce)
    if err != nil {
        return errors.New("SetPass: cannot set nonce")
    }

    server.cache[data.Filename] = data.Password
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
    server := NewServer()
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
