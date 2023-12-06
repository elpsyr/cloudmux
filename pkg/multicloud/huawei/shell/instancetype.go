// Copyright 2019 Yunion
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package shell

import (
	"yunion.io/x/pkg/util/shellutils"

	"yunion.io/x/cloudmux/pkg/multicloud/huawei"
)

func init() {
	type InstanceMatchOptions struct {
		CPU  int    `help:"CPU count"`
		MEM  int    `help:"Memory in MB"`
		Zone string `help:"Test in zone"`
	}
	shellutils.R(&InstanceMatchOptions{}, "instance-type-select", "Select matching instance types", func(cli *huawei.SRegion, args *InstanceMatchOptions) error {
		instanceTypes, e := cli.GetMatchInstanceTypes(args.CPU, args.MEM, args.Zone)
		if e != nil {
			return e
		}
		printList(instanceTypes, 0, 0, 0, []string{})
		return nil
	})

	type InstanceTypeGetOptions struct {
		Zone string `help:"Test in zone"`
	}
	// yunion.io/x/cloudmux/cmd/huaweicli --access-key xxx  --secret xxx  subaccount-list
	// yunion.io/x/cloudmux/cmd/huaweicli --region cn-east-3 --project projectId --access-key xxx  --secret xxx   instance-types
	shellutils.R(&InstanceTypeGetOptions{Zone: "cn-east-3a"}, "instance-types", "Get all instance types", func(cli *huawei.SRegion, args *InstanceTypeGetOptions) error {
		instanceTypes, e := cli.GetInstanceTypes(args.Zone)
		if e != nil {
			return e
		}
		printList(instanceTypes, 0, 0, 0, []string{})
		return nil
	})

	type InstanceTypeRegionGetOptions struct {
	}
	// yunion.io/x/cloudmux/cmd/huaweicli --access-key xxx  --secret xxx  subaccount-list
	// yunion.io/x/cloudmux/cmd/huaweicli --region cn-east-3 --project projectId --access-key xxx  --secret xxx   region-instance-types
	shellutils.R(&InstanceTypeRegionGetOptions{}, "region-instance-types", "Get region all instance types", func(cli *huawei.SRegion, args *InstanceTypeRegionGetOptions) error {
		instanceTypes, e := cli.GetRegionInstanceTypes()
		if e != nil {
			return e
		}
		printList(instanceTypes, 0, 0, 0, []string{})
		return nil
	})

	type InstanceTypeStatus struct {
	}
	// yunion.io/x/cloudmux/cmd/huaweicli --access-key xxx  --secret xxx  subaccount-list
	// yunion.io/x/cloudmux/cmd/huaweicli --region cn-east-3 --project projectId --access-key xxx  --secret xxx   region-instance-types
	shellutils.R(&InstanceTypeStatus{}, "instance-type-status", "Get instance type status", func(cli *huawei.SRegion, args *InstanceTypeStatus) error {
		instanceTypes, e := cli.GetInstanceTypeStatus("cn-east-3a", "ac7.12xlarge.2")
		instanceTypes, e = cli.GetInstanceTypeStatus("cn-east-3c", "ac7.12xlarge.2")
		instanceTypes, e = cli.GetInstanceTypeStatus("cn-east-3d", "ac7.12xlarge.2")
		if e != nil {
			return e
		}
		printList(instanceTypes, 0, 0, 0, []string{})
		return nil
	})

	type InstanceTypePrice struct {
	}
	// yunion.io/x/cloudmux/cmd/huaweicli --access-key xxx  --secret xxx  subaccount-list
	// yunion.io/x/cloudmux/cmd/huaweicli --region cn-east-3 --project projectId --access-key xxx  --secret xxx   region-instance-types
	shellutils.R(&InstanceTypePrice{}, "instance-type-post-price", "Get instance type price", func(cli *huawei.SRegion, args *InstanceTypePrice) error {
		instanceTypes, e := cli.GetPostPaidPrice("cn-east-3a", "ac7.12xlarge.2")
		instanceTypes, e = cli.GetPostPaidPrice("cn-east-3c", "ac7.12xlarge.2")
		instanceTypes, e = cli.GetPostPaidPrice("cn-east-3d", "ac7.12xlarge.2")
		if e != nil {
			return e
		}
		printList(instanceTypes, 0, 0, 0, []string{})
		return nil
	})
}
