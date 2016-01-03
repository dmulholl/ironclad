package ironrpc


import (
    "net"
    "net/rpc"
    "time"
)


// Timeout duration for client connections.
var ClientTimeout = 100 * time.Millisecond


// RPC client for connecting to the password server.
type Client struct {
    connection *rpc.Client
}


// NewClient returns an initialized RPC client.
func NewClient(address string) (*Client, error) {
    connection, err := net.DialTimeout("tcp", address, ClientTimeout)
    if err != nil {
        return nil, err
    }
    return &Client{connection: rpc.NewClient(connection)}, nil
}


// Get attempts to fetch a cached password from the server.
func (client *Client) Get(token string) (password string, err error) {
    err = client.connection.Call("Server.Get", token, &password)
    return password, err
}


// Set passes a token and password pair to the server.
func (client *Client) Set(token, password string) (ok bool, err error) {
    pair := TokenPair{ Token: token, Password: password }
    err = client.connection.Call("Server.Set", pair, &ok)
    return ok, err
}
