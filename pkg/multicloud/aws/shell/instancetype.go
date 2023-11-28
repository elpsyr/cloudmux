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

	"yunion.io/x/cloudmux/pkg/multicloud/aws"
)

func init() {
	type InstanceTypeListOptions struct {
	}
	shellutils.R(&InstanceTypeListOptions{}, "instance-type-list", "List intance types", func(cli *aws.SRegion, args *InstanceTypeListOptions) error {
		skus, err := cli.GetInstanceTypes()
		if err != nil {
			return err
		}
		printList(skus, 0, 0, 0, []string{})
		return nil
	})

	type SkuListOptions struct {
		Arch      string
		NextToken string
	}
	shellutils.R(&SkuListOptions{}, "sku-list", "List intance types", func(cli *aws.SRegion, args *SkuListOptions) error {
		skus, _, err := cli.DescribeInstanceTypes(args.Arch, args.NextToken)
		if err != nil {
			return err
		}
		printList(skus, 0, 0, 0, []string{})
		return nil
	})

	shellutils.R(&SkuListOptions{}, "sku-list-all", "List all intance types", func(cli *aws.SRegion, args *SkuListOptions) error {
		skus, err := cli.DescribeInstanceTypesAll()
		if err != nil {
			return err
		}
		printList(skus, 0, 0, 0, []string{})
		return nil
	})

	shellutils.R(&SkuListOptions{}, "sku-offering-list-all", "List all intance types offering", func(cli *aws.SRegion, args *SkuListOptions) error {
		skus, err := cli.DescribeInstanceTypeOfferingsAll()
		if err != nil {
			return err
		}
		printList(skus, 0, 0, 0, []string{})
		return nil
	})

	// 获取带 zone的sku 列表
	shellutils.R(&SkuListOptions{}, "zone-offering-sku-list", "List all intance types offering", func(cli *aws.SRegion, args *SkuListOptions) error {
		skus, err := cli.GetISkus()
		if err != nil {
			return err
		}
		printList(skus, 0, 0, 0, []string{})
		return nil
	})
}
