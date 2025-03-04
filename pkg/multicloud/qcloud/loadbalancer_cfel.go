package qcloud

import (
	"fmt"
	"time"

	sdkerrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/jsonutils"
	"yunion.io/x/pkg/errors"
)

// https://cloud.tencent.com/document/api/214/30692
// 这个接口只能创建按量付费，不能创建包年包月
func (self *SRegion) CreateILoadBalancer(opts *cloudprovider.SLoadbalancerCreateOptions) (cloudprovider.ICloudLoadbalancer, error) {
	params := map[string]string{
		"LoadBalancerName": opts.Name,
		"VpcId":            opts.VpcId,
	}

	LoadBalancerType := "INTERNAL"
	if opts.AddressType == api.LB_ADDR_TYPE_INTERNET { // 公网 实例计费方式有按量计费和包年包月
		LoadBalancerType = "OPEN"
		switch opts.InstanceChargeType {
		case "PrePaid": // 包年包月
			bps := opts.EgressMbps
			if bps == 0 {
				bps = 200
			}
			period := opts.BillingCycle.Count
			if opts.BillingCycle.Unit == "Year" {
				period *= 12
			}
			params["LBChargeType"] = "PREPAID"
			params["LBChargePrepaid.Period"] = fmt.Sprintf("%d", period)

			params["InternetAccessible.InternetChargeType"] = "BANDWIDTH_PREPAID"
			params["InternetAccessible.InternetMaxBandwidthOut"] = fmt.Sprintf("%d", bps)

		default: // 按量付费
			bps := opts.EgressMbps
			if bps == 0 {
				bps = 200
			}
			if opts.ChargeType == api.LB_CHARGE_TYPE_BY_BANDWIDTH { // 按带宽
				params["InternetAccessible.InternetChargeType"] = "BANDWIDTH_POSTPAID_BY_HOUR"
				params["InternetAccessible.InternetMaxBandwidthOut"] = fmt.Sprintf("%d", bps)
			} else { // 按流量
				params["InternetAccessible.InternetChargeType"] = "TRAFFIC_POSTPAID_BY_HOUR"
				params["InternetAccessible.InternetMaxBandwidthOut"] = fmt.Sprintf("%d", bps)
			}

		}
	} else { // 内网 只有按量付费 LoadbalancerSpec 有值即非共享型有效
		if opts.EgressMbps > 0 {
			params["InternetAccessible.InternetMaxBandwidthOut"] = fmt.Sprintf("%d", opts.EgressMbps)
		}
	}

	params["LoadBalancerType"] = LoadBalancerType

	if len(opts.ProjectId) > 0 {
		params["ProjectId"] = opts.ProjectId
	}

	if len(opts.LoadbalancerSpec) > 0 {
		params["SlaType"] = opts.LoadbalancerSpec
	}

	if opts.AddressType != api.LB_ADDR_TYPE_INTERNET {
		params["SubnetId"] = opts.NetworkIds[0]
	} else {
		// 公网类型ELB可支持多可用区
		if len(opts.ZoneId) > 0 {
			if len(opts.SlaveZoneId) > 0 {
				params["MasterZoneId"] = opts.ZoneId
			} else {
				params["ZoneId"] = opts.ZoneId
			}
		}
	}
	i := 0
	for k, v := range opts.Tags {
		params[fmt.Sprintf("Tags.%d.TagKey", i)] = k
		params[fmt.Sprintf("Tags.%d.TagValue", i)] = v
		i++
	}

	resp, err := func() (jsonutils.JSONObject, error) {
		_resp, err := self.clbRequest("CreateLoadBalancer", params)
		if err != nil {
			// 兼容不支持指定zone的账号
			if e, ok := err.(*sdkerrors.TencentCloudSDKError); ok && e.Code == "InvalidParameterValue" {
				delete(params, "ZoneId")
				delete(params, "MasterZoneId")
				return self.clbRequest("CreateLoadBalancer", params)
			}
		}
		return _resp, err
	}()
	if err != nil {
		return nil, errors.Wrapf(err, "CreateLoadBalancer")
	}

	ret := struct {
		RequestId       string
		LoadBalancerIds []string
	}{}
	err = resp.Unmarshal(&ret)
	if err != nil {
		return nil, errors.Wrapf(err, "resp.Unmarshal")
	}
	if len(ret.RequestId) == 0 || len(ret.LoadBalancerIds) != 1 {
		return nil, errors.Wrapf(cloudprovider.ErrNotFound, resp.String())
	}
	err = self.WaitLBTaskSuccess(ret.RequestId, 5*time.Second, time.Minute*1)
	if err != nil {
		return nil, errors.Wrapf(err, "WaitLBTaskSuccess")
	}
	return self.GetLoadbalancer(ret.LoadBalancerIds[0])
}

func (s *SLoadbalancer) CfelCreateILoadBalancerBackendGroup(bg *cloudprovider.SCfelLoadbalancerBackendGroup) (cloudprovider.ICloudLoadbalancerBackendGroup, error) {
	return nil, nil
}

// https://cloud.tencent.com/document/product/214/30680
func (s *SLoadbalancer) CfelModifyLoadBalancerAttributes(opts *cloudprovider.SCfelModifyLbAttributes) error {
	params := map[string]string{
		"LoadBalancerId": s.LoadBalancerId,
	}
	if len(opts.LoadbalancerName) > 0 {
		params["LoadBalancerName"] = opts.LoadbalancerName
	}

	if opts.Bandwidth > 0 {
		// 修改成相同的会报错
		if opts.Bandwidth == s.NetworkAttributes.InternetMaxBandwidthOut {
			return nil
		}
		// 不支持接口改网络计费方式
		params["InternetChargeInfo.InternetMaxBandwidthOut"] = fmt.Sprintf("%d", opts.Bandwidth)
		if opts.ChargeType == "PrePaid" {
			params["InternetChargeInfo.InternetChargeType"] = "BANDWIDTH_PREPAID"
		} else if len(opts.InternetChargeType) > 0 {
			if opts.InternetChargeType == "traffic" {
				params["InternetChargeInfo.InternetChargeType"] = "TRAFFIC_POSTPAID_BY_HOUR"
			} else {
				params["InternetChargeInfo.InternetChargeType"] = "BANDWIDTH_POSTPAID_BY_HOUR"
			}
		}
	} else {
		params["LoadBalancerPassToTarget"] = fmt.Sprintf("%v", opts.LoadBalancerPassToTarget)
	}

	resp, err := s.region.clbRequest("ModifyLoadBalancerAttributes", params)
	if err != nil {
		return err
	}
	requestId, err := resp.GetString("RequestId")
	if err != nil {
		return err
	}
	return s.region.WaitLBTaskSuccess(requestId, 5*time.Second, 60*time.Second)
}

func (s *SLoadbalancer) CfelSetLoadBalancerSecurityGroups(sgs []string) error {
	params := map[string]string{
		"LoadBalancerId": s.LoadBalancerId,
		// "SecurityGroups.0":           sgs[0],
	}
	for i := range sgs {
		params[fmt.Sprintf("SecurityGroups.%d", i)] = sgs[i]
	}
	resp, err := s.region.clbRequest("SetLoadBalancerSecurityGroups", params)
	if err != nil {
		return err
	}
	requestId, err := resp.GetString("RequestId")
	if err != nil {
		return err
	}
	return s.region.WaitLBTaskSuccess(requestId, 5*time.Second, 60*time.Second)
}

func (s *SLoadbalancer) CfelUnSetLoadBalancerSecurityGroups(sgs []string) error {
	params := map[string]string{
		"LoadBalancerIds.0": s.LoadBalancerId,
		"SecurityGroup":     sgs[0],
		"OperationType":     "DEL",
	}

	resp, err := s.region.clbRequest("SetSecurityGroupForLoadbalancers", params)
	if err != nil {
		return err
	}
	requestId, err := resp.GetString("RequestId")
	if err != nil {
		return err
	}
	return s.region.WaitLBTaskSuccess(requestId, 5*time.Second, 60*time.Second)
}
