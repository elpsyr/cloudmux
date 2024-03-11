package aws

import "yunion.io/x/cloudmux/pkg/cloudprovider"

// Verify that *SRegion implements ICfelCloudRegion
var _ cloudprovider.ICfelCloudRegion = (*SRegion)(nil)
