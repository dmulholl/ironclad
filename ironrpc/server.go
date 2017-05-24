/*
    Package ironrpc implements password-caching between application runs.
*/
package ironrpc


import (
    "net"
    "net/rpc"
    "time"
    "sync"
    "errors"
    "os"
)


// Error returned when an invalid token is presented to the server.
var TokenError = errors.New("invalid token")


// Duration after which the server will automatically shut down.
var ServerTimeout = 60 * time.Minute


// TokenPair instances are used internally by the RPC implementation.
type TokenPair struct {
    Token string
    Password string
}


// RPC password server.
type Server struct {
    password string
    token string
    mutex *sync.Mutex
    lastaccess time.Time
}


// NewServer returns an initialized RPC password server.
func NewServer() *Server {
    return &Server{
        mutex: &sync.Mutex{},
        lastaccess: time.Now(),
    }
}


// Get method exposed by the RPC server.
func (server *Server) Get(token string, password *string) error {
    server.mutex.Lock()
    defer server.mutex.Unlock()

    if token != server.token {
        time.Sleep(time.Second)
        return TokenError
    }

    *password = server.password
    server.lastaccess = time.Now()
    return nil
}


// Set method exposed by the RPC server.
func (server *Server) Set(pair TokenPair, ok *bool) error {
    server.mutex.Lock()
    defer server.mutex.Unlock()

    server.token = pair.Token
    server.password = pair.Password

    *ok = true
    server.lastaccess = time.Now()
    return nil
}


// Automatically shuts the server down after the specified duration.
func (server *Server) timeout() {
    for {
        server.mutex.Lock()
        if time.Since(server.lastaccess) > ServerTimeout {
            os.Exit(0)
        }
        server.mutex.Unlock()
        time.Sleep(time.Second)
    }
}


// Serve launches a new RPC password server and blocks until the server
// automatically shuts itself down.
func Serve(address string) error {
    server := NewServer()
    go server.timeout()

    err := rpc.Register(server)
    if err != nil {
        return err
    }

    listener, err := net.Listen("tcp", address)
    if err != nil {
        return err
    }

    rpc.Accept(listener)
    return nil
}
