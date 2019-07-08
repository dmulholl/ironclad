package ironrpc


import (
    "errors"
    "net"
    "net/rpc"
    "time"
)


import (
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
func (client *CacheClient) GetPass(filename, nonce string) (string, error) {
    var password string
    data := GetData{ Filename: filename, Nonce: nonce }
    err := client.rpcc.Call("CacheServer.GetPass", data, &password)
    return password, err
}


// SetPass attempts to cache a password to the server.
func (client *CacheClient) SetPass(filename, password string) error {
    var notused bool
    data := SetData{ Filename: filename, Password: password }
    return client.rpcc.Call("CacheServer.SetPass", data, &notused)
}


// Close calls the underlying net/rpc.Client's Close() method.
func (client *CacheClient) Close() {
    if client.rpcc != nil {
        client.rpcc.Close()
    }
}
