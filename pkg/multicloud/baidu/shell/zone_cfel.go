package shell

import (
	"fmt"
	"yunion.io/x/cloudmux/pkg/multicloud/baidu"
	"yunion.io/x/pkg/util/shellutils"
)

func init() {
	type ZoneListOptions struct {
	}
	// yunion.io/x/cloudmux/cmd/baiducli  --region-id bj zone-list
	shellutils.R(&ZoneListOptions{}, "zone-list", "list zones", func(cli *baidu.SRegion, args *ZoneListOptions) error {
		zones, err := cli.GetIZones()
		if err != nil {
			return err
		}
		printList(zones, 0, 0, 0, []string{})
		for _, i := range zones {
			fmt.Println(i.GetName())
		}
		return nil
	})

}
