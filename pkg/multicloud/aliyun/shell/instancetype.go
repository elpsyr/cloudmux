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

	"yunion.io/x/cloudmux/pkg/multicloud/aliyun"
)

func init() {
	type InstanceTypeListOptions struct {
	}
	shellutils.R(&InstanceTypeListOptions{}, "instance-type-list", "List intance types", func(cli *aliyun.SRegion, args *InstanceTypeListOptions) error {
		instanceTypes, e := cli.GetInstanceTypes()
		if e != nil {
			return e
		}
		printList(instanceTypes, 0, 0, 0, []string{})
		return nil
	})

	type InstanceMatchOptions struct {
		CPU  int    `help:"CPU count"`
		MEM  int    `help:"Memory in MB"`
		GPU  int    `help:"GPU size"`
		Zone string `help:"Test in zone"`
	}
	shellutils.R(&InstanceMatchOptions{}, "instance-type-select", "Select matching instance types", func(cli *aliyun.SRegion, args *InstanceMatchOptions) error {
		instanceTypes, e := cli.GetMatchInstanceTypes(args.CPU, args.MEM, args.GPU, args.Zone)
		if e != nil {
			return e
		}
		printList(instanceTypes, 0, 0, 0, []string{})
		return nil
	})

	// 获取 instanceType spot 价格
	shellutils.R(&InstanceTypeListOptions{}, "zone-instance-spot-price", "", func(cli *aliyun.SRegion, args *InstanceTypeListOptions) error {
		cli.RegionId = "us-east-1"
		_, err := cli.GetDescribePrice("us-east-1a", "ecs.gn7i-c32g1.32xlarge", "SPOTPOSTPAID")

		if err != nil {
			return err
		}
		//printList(skus, 0, 0, 0, []string{})
		return nil
	})
}
