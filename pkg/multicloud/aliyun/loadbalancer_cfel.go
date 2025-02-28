package aliyun

import (
	"fmt"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
)

func (s *SLoadbalancer) CfelCreateILoadBalancerBackendGroup(bg *cloudprovider.SCfelLoadbalancerBackendGroup) (cloudprovider.ICloudLoadbalancerBackendGroup, error) {
	return nil, nil
}

// https://cloud.tencent.com/document/product/214/30680
func (s *SLoadbalancer) CfelModifyLoadBalancerAttributes(opts *cloudprovider.SCfelModifyLbAttributes) error {
	params := map[string]string{
		"LoadBalancerId": s.LoadBalancerId,
		"AutoPay":        "true",
	}

	if opts.Bandwidth > 0 {
		params["Bandwidth"] = fmt.Sprintf("%d", opts.Bandwidth)
	}

	if len(opts.InternetChargeType) > 0 {
		if opts.InternetChargeType == "bandwidth" {
			params["InternetChargeType"] = "paybybandwidth"
		} else {
			params["InternetChargeType"] = "paybytraffic"
		}
	}
	_, err := s.region.lbRequest("ModifyLoadBalancerInternetSpec ", params)
	return err
}

func (s *SLoadbalancer) CfelSetLoadBalancerSecurityGroups(sgs []string) error {
	return nil
}

func (s *SLoadbalancer) CfelUnSetLoadBalancerSecurityGroups(sgs []string) error {
	return nil
}
