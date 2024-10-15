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

package huawei

import (
	"context"
	"fmt"
	"strings"
	"time"

	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
)

// https://console.huaweicloud.com/apiexplorer/#/openapi/ECS/doc?api=NovaRebootServer
func (self *SRegion) RebootVM(instanceId string, isForce bool) error {
	stopType := "SOFT"
	if isForce {
		stopType = "HARD"
	}
	params := map[string]interface{}{
		"reboot": map[string]string{
			"type": stopType,
		},
	}
	_, err := self.post(SERVICE_ECS_V2_1, fmt.Sprintf("servers/%s/action", instanceId), params)
	if err != nil && strings.Contains(err.Error(), "empty response") {
		err = nil
	}
	return err
}

func (self *SInstance) GetIHostId() string {
	return fmt.Sprintf("%s-%s", "", self.OSEXTAZAvailabilityZone)
}

func (self *SInstance) RebootVM(_ context.Context) error {
	err := self.host.zone.region.RebootVM(self.GetId(), false)
	if err != nil {
		return err
	}
	return cloudprovider.WaitStatus(self, api.VM_RUNNING, 10*time.Second, 300*time.Second) // 5mintues
}
