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

	// 获取 instanceType spot 价格
	// --access-key xxx --secret xxx --region-id me-east-1 zone-instance-spot-price
	shellutils.R(&InstanceTypeListOptions{}, "zone-instance-spot-price", "", func(cli *aliyun.SRegion, args *InstanceTypeListOptions) error {
		cli.RegionId = "us-east-1"
		price, err := cli.GetDescribePrice("us-east-1a", "ecs.gn7i-c32g1.32xlarge", "SPOTPOSTPAID")

		if err != nil {
			return err
		}
		printObject(price)
		return nil
	})

	// 获取 instanceType post 价格
	// --region cn-wuhan-lr   --access-key xxx  --secret xxx  zone-instance-post-price
	shellutils.R(&InstanceTypeListOptions{}, "zone-instance-post-price", "", func(cli *aliyun.SRegion, args *InstanceTypeListOptions) error {
		cli.RegionId = "cn-wuhan-lr"
		//_, err := cli.GetDescribePrice("cn-wuhan-lr", "ecs.r7.2xlarge", "PostPaid")
		_, err := cli.GetDescribePrice("cn-wuhan-lr", "ecs.t1.small", "PostPaid")

		if err != nil {
			return err
		}
		//printList(skus, 0, 0, 0, []string{})
		return nil
	})
}
