package driver

import (
	"github.com/container-storage-interface/spec/lib/go/csi"
)

var (
	volumeAccessModes = []*csi.VolumeCapability_AccessMode{
		&csi.VolumeCapability_AccessMode{
			Mode: csi.VolumeCapability_AccessMode_MULTI_NODE_READER_ONLY,
		},
	}
)

func hasCapabilities(caps []*csi.VolumeCapability) bool {
	hasSupport := func(mode csi.VolumeCapability_AccessMode_Mode) bool {
		for _, m := range volumeAccessModes {
			if mode == m.Mode {
				return true
			}
		}
		return false
	}

	supported := false

	for _, cap := range caps {
		if hasSupport(cap.AccessMode.Mode) {
			supported = true
		} else {
			// we need to make sure all capabilities are supported. Revert back
			// in case we have a cap that is supported, but is invalidated now
			supported = false
		}
	}

	return supported
}
