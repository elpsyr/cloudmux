package shell

import (
	"yunion.io/x/cloudmux/pkg/multicloud/baidu"
	"yunion.io/x/pkg/util/shellutils"
)

func init() {
	type InstanceTypeListOptions struct {
	}
	// yunion.io/x/cloudmux/cmd/baiducli  --access-key-id xxx   --access-key-secret  xxx    --region-id bj instance-type-list
	shellutils.R(&InstanceTypeListOptions{}, "instance-type-list", "list zones", func(cli *baidu.SRegion, args *InstanceTypeListOptions) error {
		flavors, err := cli.GetISkus()
		if err != nil {
			return err
		}
		printList(flavors, 0, 0, 0, []string{})

		return nil
	})

	type InstanceTypePriceOptions struct {
	}
	// yunion.io/x/cloudmux/cmd/baiducli  --access-key-id xxx   --access-key-secret  xxx    --region-id bj instance-type-post-price
	shellutils.R(&InstanceTypePriceOptions{}, "instance-type-post-price", "instance-type post price", func(cli *baidu.SRegion, args *InstanceTypePriceOptions) error {
		price, err := cli.GetPostPaidPrice("cn-bj-d", "bcc.g5.c2m8")
		if err != nil {
			return err
		}
		printObject(struct {
			Price float64
		}{
			Price: price,
		})

		// 售罄
		price, err = cli.GetPostPaidPrice("cn-bj-b", "bcc.g5.c2m8")
		if err != nil {
			return err
		}
		printObject(struct {
			Price float64
		}{
			Price: price,
		})

		return nil
	})

	// yunion.io/x/cloudmux/cmd/baiducli  --access-key-id xxx   --access-key-secret  xxx    --region-id bj instance-type-post-price
	shellutils.R(&InstanceTypePriceOptions{}, "instance-type-pre-price", "instance-type pre-paid price", func(cli *baidu.SRegion, args *InstanceTypePriceOptions) error {
		price, err := cli.GetPrePaidPrice("cn-bj-d", "bcc.g5.c2m8")
		if err != nil {
			return err
		}
		printObject(struct {
			Price float64
		}{
			Price: price,
		})

		return nil
	})

}
