package driver

import (
	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/phogolabs/log"

	csicommon "github.com/kubernetes-csi/drivers/pkg/csi-common"
)

const (
	driverName = "csi-vault"
	version    = "1.0.0-rc2"
)

// Driver for Vault
type Driver struct {
	csiDriver *csicommon.CSIDriver
	endpoint  string

	ids *csicommon.DefaultIdentityServer
	ns  *nodeServer

	cap   []*csi.VolumeCapability_AccessMode
	cscap []*csi.ControllerServiceCapability
}

// New creates a new driver
func New(nodeID, endpoint string) *Driver {
	log.Infof("Driver: %v version: %v", driverName, version)

	d := &Driver{}

	d.endpoint = endpoint

	csiDriver := csicommon.NewCSIDriver(driverName, version, nodeID)
	csiDriver.AddVolumeCapabilityAccessModes([]csi.VolumeCapability_AccessMode_Mode{csi.VolumeCapability_AccessMode_MULTI_NODE_READER_ONLY})
	// NFS plugin does not support ControllerServiceCapability now.
	// If support is added, it should set to appropriate
	// ControllerServiceCapability RPC types.
	csiDriver.AddControllerServiceCapabilities([]csi.ControllerServiceCapability_RPC_Type{csi.ControllerServiceCapability_RPC_UNKNOWN})

	d.csiDriver = csiDriver

	return d
}

// Run starts the driver
func (d *Driver) Run() {
	s := csicommon.NewNonBlockingGRPCServer()
	s.Start(d.endpoint,
		csicommon.NewDefaultIdentityServer(d.csiDriver),
		// NFS plugin has not implemented ControllerServer.
		nil,
		&nodeServer{
			server: csicommon.NewDefaultNodeServer(d.csiDriver),
		},
	)
	s.Wait()
}
