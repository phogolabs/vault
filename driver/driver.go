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
	"github.com/davecgh/go-spew/spew"
	"github.com/kubernetes-csi/csi-lib-utils/protosanitizer"
	"github.com/phogolabs/log"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	name    = "csi-vault"
	version = "1.0.0"
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
	scheme, addr, err := endpoint(d.endpoint)
	if err != nil {
		return err
	}

	listener, err := net.Listen(scheme, addr)
	if err != nil {
		return err
	}

	d.srv = grpc.NewServer(grpc.UnaryInterceptor(logger))

	csi.RegisterIdentityServer(d.srv, d)
	csi.RegisterControllerServer(d.srv, d)
	csi.RegisterNodeServer(d.srv, d)

	d.ready = true
	return d.srv.Serve(listener)
}

// Stop stops the plugin
func (d *Driver) Stop() {
	d.readyMu.Lock()
	d.ready = false
	d.readyMu.Unlock()
	d.srv.Stop()
}

func endpoint(endpoint string) (string, string, error) {
	uri, err := parse(endpoint)
	if err != nil {
		return "", "", err
	}

	addr := path.Join(uri.Host, filepath.FromSlash(uri.Path))

	if uri.Scheme == "unix" {
		if err := os.Remove(addr); err != nil && !os.IsNotExist(err) {
			return "", "", fmt.Errorf("failed to remove unix domain socket file %s, error: %s", addr, err)
		}
	}

	return uri.Scheme, addr, nil
}

func parse(endpoint string) (*url.URL, error) {
	uri, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	if uri.Scheme == "unix" || uri.Scheme == "tcp" {
		return uri, nil
	}

	return nil, fmt.Errorf("invalid endpoint: %v", endpoint)
}

func logger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	//TODO: remove
	spew.Dump(req)

	logger := log.WithFields(log.Fields{
		"method":       info.FullMethod,
		"request_body": protosanitizer.StripSecrets(req),
	})

	logger.Info("call")

	resp, err := handler(ctx, req)
	if err != nil {
		logger.WithError(err).Error("failed")
	} else {
		logger.WithField("response_body", protosanitizer.StripSecrets(resp)).Info("succeeded")

		//TODO: remove
		spew.Dump(resp)
	}

	return resp, err
}
