package shell

import (
	"yunion.io/x/pkg/util/shellutils"

	"yunion.io/x/cloudmux/pkg/multicloud/huawei"
)

func init() {
	type InstanceMatchOptions struct {
		Zone string `help:"Test in zone"`
	}

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

	// 售卖情况
	type InstanceTypeStatus struct {
	}
	// yunion.io/x/cloudmux/cmd/huaweicli --region cn-east-3 --project projectId --access-key xxx  --secret xxx   instance-type-status
	shellutils.R(&InstanceTypeStatus{}, "instance-type-status", "Get instance type status", func(cli *huawei.SRegion, args *InstanceTypeStatus) error {
		instanceTypes, e := cli.GetInstanceTypeStatus("cn-east-3a", "ac7.12xlarge.2")
		instanceTypes, e = cli.GetInstanceTypeStatus("cn-east-3c", "ac7.12xlarge.2")
		instanceTypes, e = cli.GetInstanceTypeStatus("cn-east-3d", "ac7.12xlarge.2")
		if e != nil {
			return e
		}
		printObject(struct {
			Status string
		}{Status: instanceTypes})
		return nil
	})

	type InstanceTypePrice struct {
	}

	// 按量付费价格
	// 1. 查询 instanceTypes yunion.io/x/cloudmux/cmd/huaweicli --region cn-north-9  --access-key xxx  --secret xxx   region-instance-types  ===> cn-north-9a ac6.12xlarge.2
	// 2. 查询 instanceType 价格 yunion.io/x/cloudmux/cmd/huaweicli --region cn-north-9 --project ff0580add98c440b978079ee4cf3eef3 --access-key xxx  --secret xxx   instance-type-post-price
	shellutils.R(&InstanceTypePrice{}, "instance-type-post-price", "Get instance type price", func(cli *huawei.SRegion, args *InstanceTypePrice) error {
		//instanceTypes, e := cli.GetPostPaidPrice("cn-east-3a", "ac7.12xlarge.2")
		//instanceTypes, e := cli.GetPostPaidPrice("cn-east-3c", "ac7.12xlarge.2")
		instanceTypePrice, e := cli.GetPostPaidPrice("cn-east-3c", "ac7.12xlarge.2")
		if e != nil {
			return e
		}
		printObject(struct {
			Price float64
		}{
			Price: instanceTypePrice,
		})
		return nil
	})

	// 包年包月价格
	// 1. 查询 instanceTypes yunion.io/x/cloudmux/cmd/huaweicli --region cn-north-9  --access-key xxx  --secret xxx   region-instance-types  ===> cn-north-9a ac6.12xlarge.2
	// 2. 查询 instanceType 价格 yunion.io/x/cloudmux/cmd/huaweicli --region cn-north-9 --project ff0580add98c440b978079ee4cf3eef3 --access-key xxx  --secret xxx   instance-type-post-price
	shellutils.R(&InstanceTypePrice{}, "instance-type-prepaid-price", "Get instance type price", func(cli *huawei.SRegion, args *InstanceTypePrice) error {
		//instanceTypes, e := cli.GetPostPaidPrice("cn-east-3a", "ac7.12xlarge.2")
		//instanceTypes, e := cli.GetPostPaidPrice("cn-east-3c", "ac7.12xlarge.2")
		instanceTypePrice, e := cli.GetPrePaidPrice("cn-east-3c", "ac7.12xlarge.2")
		if e != nil {
			return e
		}
		printObject(struct {
			Price float64
		}{
			Price: instanceTypePrice,
		})
		return nil
	})
}
