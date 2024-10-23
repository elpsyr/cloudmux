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

package ecloudcfel

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"time"

	"yunion.io/x/jsonutils"
	"yunion.io/x/pkg/errors"

	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
)

type SHost struct {
	multicloud.SHostBase
	zone *SZone
}

func (h *SHost) GetId() string {
	return fmt.Sprintf("%s-%s", h.zone.region.client.cpcfg.Id, h.zone.GetId())
}

func (h *SHost) GetName() string {
	return fmt.Sprintf("%s-%s", h.zone.region.client.cpcfg.Name, h.zone.GetId())
}

func (h *SHost) GetGlobalId() string {
	return h.GetId()
}

func (h *SHost) GetStatus() string {
	return api.HOST_STATUS_RUNNING
}

func (h *SHost) Refresh() error {
	return nil
}

func (h *SHost) IsEmulated() bool {
	return true
}

func (h *SHost) GetIVMs() ([]cloudprovider.ICloudVM, error) {
	zoneRegion := h.zone.Region
	vms, err := h.zone.region.GetInstances(zoneRegion)
	if err != nil {
		return nil, errors.Wrap(err, "SHost.GetVMs")
	}
	ivms := make([]cloudprovider.ICloudVM, len(vms))
	for i := range vms {
		vms[i].host = h
		ivms[i] = &vms[i]
	}
	return ivms, nil
}

func (h *SHost) GetIVMById(id string) (cloudprovider.ICloudVM, error) {
	vm, err := h.zone.region.GetInstanceById(id)
	if err != nil {
		return nil, err
	}
	vm.host = h
	return vm, nil
}

func (h *SHost) GetIStorages() ([]cloudprovider.ICloudStorage, error) {
	return h.zone.GetIStorages()
}

func (h *SHost) GetIStorageById(id string) (cloudprovider.ICloudStorage, error) {
	return h.zone.GetIStorageById(id)
}

func (h *SHost) GetEnabled() bool {
	return true
}

func (h *SHost) GetHostStatus() string {
	return api.HOST_ONLINE
}

func (h *SHost) GetAccessIp() string {
	return ""
}

func (h *SHost) GetAccessMac() string {
	return ""
}

func (h *SHost) GetSysInfo() jsonutils.JSONObject {
	info := jsonutils.NewDict()
	info.Add(jsonutils.NewString(api.CLOUD_PROVIDER_ECLOUD), "manufacture")
	return info
}

func (h *SHost) GetSN() string {
	return ""
}

func (h *SHost) GetCpuCount() int {
	return 0
}

func (h *SHost) GetNodeCount() int8 {
	return 0
}

func (h *SHost) GetCpuDesc() string {
	return ""
}

func (h *SHost) GetCpuMhz() int {
	return 0
}

func (h *SHost) GetMemSizeMB() int {
	return 0
}

func (h *SHost) GetStorageSizeMB() int64 {
	return 0
}

func (h *SHost) GetStorageType() string {
	return api.DISK_TYPE_HYBRID
}

func (h *SHost) GetHostType() string {
	return api.HOST_TYPE_ECLOUD
}

func (h *SHost) GetIsMaintenance() bool {
	return false
}

func (h *SHost) GetVersion() string {
	return CLOUD_API_VERSION
}

type createResp struct {
	OrderId string `json:"orderId"`
}

type orderInfo struct {
	InstanceId string `json:"instanceId"`
}

func (h *SHost) CreateVM(desc *cloudprovider.SManagedVMCreateConfig) (cloudprovider.ICloudVM, error) {
	password, err := rsaEncryptPassword(desc.Password)
	if err != nil {
		return nil, err
	}
	params := map[string]interface{}{
		"region":      h.zone.Region,
		"billingType": "HOUR",
		"vmType":      "common",
		"cpu":         desc.Cpu,
		"ram":         desc.MemoryMB / 1024,
		// "disk":        40,
		"specsName": desc.InstanceType,
		"bootVolume": map[string]interface{}{
			"size":       desc.SysDisk.SizeGB,
			"volumeType": desc.SysDisk.StorageType,
		},
		"imageName": desc.ExternalImageId,
		"networks": map[string]interface{}{
			"networkId": desc.ExternalNetworkId,
		},
		"name":             desc.Name,
		"quantity":         1,
		"securityGroupIds": desc.ExternalSecgroupIds,
		"userData":         desc.UserData,
		"password":         password,
	}
	req := NewConsoleRequest(h.zone.region.ID, "/api/openapi-ecs/acl/v3/server/order", nil, jsonutils.Marshal(params))
	res, err := h.zone.region.client.doPost(req)
	if err != nil {
		return nil, err
	}
	var ret createResp
	if err := res.Unmarshal(&ret); err != nil {
		return nil, err
	}
	
	var ids []string
	for i := 1; i < 10; i++ {
		time.Sleep(6 * time.Second)
		ids, err = h.zone.region.getOrderInfo(ret.OrderId)
		if err != nil {
			fmt.Printf("getOrderInfo [%d] err: %v", i, err)
		}
		if len(ids) > 0 {
			break
		}
	}
	var ins []SInstance
	for _,id := range ids {
		ins,err = h.zone.region.getInstances(h.zone.Region,id)
		if len(ins) > 0 {
			break
		}
	}
	var vm SInstance
	if len(ins) != 0 {
		vm = ins[0]
	} else {
		vm.Id = ret.OrderId
	}
	return &vm, nil
}

func (r *SRegion) getOrderInfo(orderId string) ([]string, error) {
	query := map[string]string{
		"orderId": orderId,
	}
	req := NewConsoleRequest(r.ID, "/api/openapi-ecs/acl/v3/server/order/relation/info", query, nil)
	var order []orderInfo
	if err := r.client.doGet(context.Background(), req, &order); err != nil {
		return nil, err
	}
	if len(order) == 0 {
		return nil, fmt.Errorf("order info is empty, orderId:%s", orderId)
	}
	var res []string
	for i,val := range order {
		if len(val.InstanceId) == 0 {
			return nil,fmt.Errorf("[%d] instanceId is empty",i + 1)
		}
		res = append(res, val.InstanceId)
	}
	return res, nil
}

func (h *SHost) GetIHostNics() ([]cloudprovider.ICloudHostNetInterface, error) {
	return nil, cloudprovider.ErrNotImplemented
}

func (h *SRegion) GetVMs() ([]SInstance, error) {
	return nil, cloudprovider.ErrNotImplemented
}

func (h *SRegion) GetVMById(vmId string) (*SInstance, error) {
	return nil, cloudprovider.ErrNotImplemented
}

var ecloudPubKey = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC/VpRysi0bPRLS7sbgQDJHo1MAt9/bK+nwK5Pe
3z0/O4cH5I/8kFNYy4yFsLMM+zyFvVw9C4wzjHaRcmEuF3ziJMC9PD5ufUWgfO5nSGgZW1cmgjqn
hcWJ3i+Azj72RnhKQRCn9DgJduEC9MiKfbyTICGd6FXf9cxb21nkxI7vtwIDAQAB
-----END PUBLIC KEY-----
`

func rsaEncryptPassword(data string) (string, error) {
	if len(data) == 0 {
		return data, nil
	}
	// 解析 PEM 公钥
	block, _ := pem.Decode([]byte(ecloudPubKey))
	if block == nil || block.Type != "PUBLIC KEY" {
		fmt.Println("failed to decode PEM block containing public key")
		return "", fmt.Errorf("failed to decode PEM block containing public key")
	}

	// 解析 PKIX 格式的公钥
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		fmt.Println("failed to parse public key:", err)
		return "", err
	}
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, pub.(*rsa.PublicKey), []byte(data))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(cipherText), nil
}
