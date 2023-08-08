// Copyright 2018 JDCLOUD.COM
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
//
// NOTE: This class is auto generated by the jdcloud code generator program.

package models


type DescribeTokenEnd struct {

    /* 数据源token (Optional) */
    AccessToken string `json:"accessToken"`

    /* ark节点 (Optional) */
    Ark string `json:"ark"`

    /*  (Optional) */
    CloudMonitorOption CloudMonitorOption `json:"cloudMonitorOption"`

    /* HawkeyeOption      *model.HawkeyeOption      `json:"hawkeyeOption,omitempty"`
DeepLogOption      *model.DeepLogOption      `json:"deepLogOption,omitempty"` (Optional) */
    CreateTime string `json:"createTime"`

    /* 数据源类型 (Optional) */
    DatasourceType string `json:"datasourceType"`

    /*  (Optional) */
    OrgId int64 `json:"orgId"`

    /*  (Optional) */
    UpdateTime string `json:"updateTime"`
}
