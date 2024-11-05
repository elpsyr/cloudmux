package huawei

import (
	"fmt"
	"net/url"
	"strings"

	"yunion.io/x/jsonutils"
	"yunion.io/x/pkg/errors"
	"yunion.io/x/pkg/util/netutils"
	"yunion.io/x/pkg/util/secrules"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
)

type SCfelLoadbalancerSku struct {
	Name string
	Id   string
	Type string
}

// GetID implements cloudprovider.ICfelLoadbalancerSku.
func (s *SCfelLoadbalancerSku) GetID() string {
	return s.Id
}

// GetName implements cloudprovider.ICfelLoadbalancerSku.
func (s *SCfelLoadbalancerSku) GetName() string {
	return s.Name
}

// GetType implements cloudprovider.ICfelLoadbalancerSku.
func (s *SCfelLoadbalancerSku) GetType() string {
	return s.Type
}

var _ cloudprovider.ICfelLoadbalancerSku = (*SCfelLoadbalancerSku)(nil)

func (region *SRegion) GetLoadbalancerSkus() ([]cloudprovider.ICfelLoadbalancerSku, error) {
	ret := jsonutils.NewArray()
	query := url.Values{}
	var res []SCfelLoadbalancerSku
	for {
		resp, err := region.list(SERVICE_ELB, "elb/flavors", query)
		if err != nil {
			return nil, err
		}
		arr, err := resp.GetArray("flavors")
		if err != nil {
			return nil, err
		}
		ret.Add(arr...)
		marker, _ := resp.GetString("page_info", "next_marker")
		if len(marker) == 0 {
			break
		}
		query.Set("marker", marker)
	}
	err := ret.Unmarshal(&res)
	if err != nil {
		return nil, err
	}
	var result []cloudprovider.ICfelLoadbalancerSku
	for i := range res {
		result = append(result, &res[i])
	}
	return result, nil
}

// GetICfelCloudImage 获取华为云镜像
func (self *SRegion) GetICfelCloudImage(withUserMeta bool) ([]cloudprovider.ICloudImage, error) {
	images := make([]SImage, 0)
	images, err := self.GetImages("", "", "gold", "")
	if err != nil {
		return nil, errors.Wrapf(err, "GetImages")
	}

	ret := []cloudprovider.ICloudImage{}
	for i := range images {
		images[i].storageCache = self.getStoragecache()
		ret = append(ret, &images[i])
	}
	return ret, nil
}

func (self *SRegion) GetInstanceOSExtraSpecs(id string) (*OSExtraSpecs, error) {
	resp, err := self.list(SERVICE_ECS_V2_1, fmt.Sprintf("flavors/%s/os-extra_specs", id), nil)
	if err != nil {
		return nil, errors.Wrapf(err, "get flavors os-extra_specs")
	}
	ret := OSExtraSpecs{}
	err = resp.Unmarshal(&ret, "extra_specs")
	if err != nil {
		return nil, errors.Wrapf(err, "Unmarshal")
	}
	return &ret, nil
}

func (self *SRegion) GetInstanceMatchImage(instanceType string) ([]cloudprovider.ICloudImage, error) {
	specs, err := self.GetInstanceOSExtraSpecs(instanceType)
	if err != nil {
		return nil, errors.Wrapf(err, "GetInstanceOSExtraSpecs")
	}

	var isArm bool
	if specs.CfelOSExtraSpecs.EcsInstanceArchitecture == "arm" ||
		specs.CfelOSExtraSpecs.EcsInstanceArchitecture == "arm64" {
		isArm = true
	} else {
		isArm = false
	}

	images, err := self.GetImages("", "", "gold", "")
	if err != nil {
		return nil, errors.Wrapf(err, "GetImages")
	}

	var ret []cloudprovider.ICloudImage
	for i := range images {
		if isArm == (strings.ToLower(images[i].SupportArm) == "true") {
			images[i].storageCache = self.getStoragecache()
			ret = append(ret, &images[i])
		}
	}
	return ret, nil
}

func (self *SRegion) DryCreateSecurityGroupRule(groupId string, opts *cloudprovider.SecurityGroupRuleCreateOptions) error {
	rule := map[string]interface{}{
		"security_group_id": groupId,
		"description":       opts.Desc,
		"direction":         "ingress",
		"ethertype":         "IPv4",
		"protocol":          strings.ToLower(opts.Protocol),
		"action":            "allow",
		"priority":          opts.Priority,
	}
	// 没有protocol表示支持所有协议，不支持空字符串
	if rule["protocol"] == "all" || rule["protocol"] == "any" {
		delete(rule, "protocol")
	}
	if len(opts.CIDR) > 0 {
		rule["remote_ip_prefix"] = opts.CIDR
		if _, err := netutils.NewIPV6Prefix(opts.CIDR); err == nil {
			rule["ethertype"] = "IPv6"
		}
	}
	if opts.Action == secrules.SecurityRuleDeny {
		rule["action"] = "deny"
	}
	if opts.Protocol == secrules.PROTO_ANY {
		delete(rule, "protocol")
	}
	if len(opts.Ports) > 0 {
		rule["multiport"] = opts.Ports
	}
	if opts.Direction == secrules.DIR_OUT {
		rule["direction"] = "egress"
	}
	params := map[string]interface{}{
		"dry_run":             true,
		"security_group_rule": rule,
	}
	_, err := self.post(SERVICE_VPC_V3, "vpc/security-group-rules", params)
	if err != nil {
		return errors.Wrapf(err, "dry create rule")
	}

	return nil
}
