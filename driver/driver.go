package driver

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	name    = "csi-vault"
	version = "1.0.0-rc2"
)

// Driver for Vault
type Driver struct {
	node     string
	endpoint string
	srv      *grpc.Server
	readyMu  sync.Mutex
	ready    bool
}

// New creates a new driver
func New(nodeID, endpoint string) *Driver {
	d := &Driver{
		node:     nodeID,
		endpoint: endpoint,
	}

	return d
}

// Run starts the driver
func (d *Driver) Run() error {
	u, err := url.Parse(d.endpoint)
	if err != nil {
		return fmt.Errorf("unable to parse address: %q", err)
	}

	addr := path.Join(u.Host, filepath.FromSlash(u.Path))
	if u.Host == "" {
		addr = filepath.FromSlash(u.Path)
	}

	// CSI plugins talk only over UNIX sockets currently
	if u.Scheme != "unix" {
		return fmt.Errorf("currently only unix domain sockets are supported, have: %s", u.Scheme)
	}

	if err := os.Remove(addr); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove unix domain socket file %s, error: %s", addr, err)
	}

	listener, err := net.Listen(u.Scheme, addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	// log response errors for better observability
	errHandler := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		return resp, err
	}

	d.srv = grpc.NewServer(grpc.UnaryInterceptor(errHandler))
	csi.RegisterIdentityServer(d.srv, d)
	csi.RegisterControllerServer(d.srv, d)
	csi.RegisterNodeServer(d.srv, d)

	d.ready = true // we're now ready to go!
	return d.srv.Serve(listener)
}

// Stop stops the plugin
func (d *Driver) Stop() {
	d.readyMu.Lock()
	d.ready = false
	d.readyMu.Unlock()
	d.srv.Stop()
}
