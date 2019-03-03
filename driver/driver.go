package driver

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/hashicorp/vault/api"
	"google.golang.org/grpc"
)

// Config is a vault csi driver's config
type Config struct {
	Endpoint string
	Node     string
}

// Driver implements the following CSI interfaces:
//
//   csi.IdentityServer
//   csi.ControllerServer
//   csi.NodeServer
//
type Driver struct {
	cfg     *Config
	client  *api.Client
	srv     *grpc.Server
	readyMu sync.Mutex
	ready   bool
}

// New returns a CSI plugin that contains the necessary gRPC
// interfaces to interact with Kubernetes over unix domain sockets for
// managaing DigitalOcean Block Storage
func New(cfg *Config) (*Driver, error) {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}

	d := &Driver{
		cfg:    cfg,
		client: client,
		srv:    grpc.NewServer(grpc.UnaryInterceptor(errorHandler)),
	}

	csi.RegisterIdentityServer(d.srv, d)
	csi.RegisterControllerServer(d.srv, d)
	csi.RegisterNodeServer(d.srv, d)
	return d, nil
}

// Run starts the CSI plugin by communication over the given endpoint
func (d *Driver) Run() error {
	addr, err := d.parseAddr()
	if err != nil {
		return err
	}

	return d.serve(addr)
}

// Stop stops the plugin
func (d *Driver) Stop() {
	d.readyMu.Lock()
	d.ready = false
	d.readyMu.Unlock()

	// TODO: log
	d.srv.Stop()
}

func (d *Driver) parseAddr() (string, error) {
	u, err := url.Parse(d.cfg.Endpoint)
	if err != nil {
		return "", fmt.Errorf("unable to parse address: %q", err)
	}

	addr := path.Join(u.Host, filepath.FromSlash(u.Path))
	if u.Host == "" {
		addr = filepath.FromSlash(u.Path)
	}

	// CSI plugins talk only over UNIX sockets currently
	if u.Scheme != "unix" {
		return "", fmt.Errorf("currently only unix domain sockets are supported, have: %s", u.Scheme)
	}

	if err = os.RemoveAll(addr); err != nil {
		return "", fmt.Errorf("failed to remove unix domain socket file %s, error: %s", addr, err)
	}

	return addr, nil
}

func (d *Driver) serve(addr string) error {
	listener, err := net.Listen("unix", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	d.ready = true
	// TODO: log
	return d.srv.Serve(listener)
}

func errorHandler(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	resp, err := handler(ctx, req)
	if err != nil {
		// TODO: log
	}
	return resp, err
}
