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


type DescribeTemplatesResponseEnd struct {

    /* 总页数 (Optional) */
    NumberPages int64 `json:"numberPages"`

    /* 总记录数 (Optional) */
    NumberRecords int64 `json:"numberRecords"`

    /* 当前页码 (Optional) */
    PageNumber int64 `json:"pageNumber"`

    /* 分页大小 (Optional) */
    PageSize int64 `json:"pageSize"`

    /* 当查询用户自定义模板时，表示该用户目前已有的自定义模板总数量;当查询默认模板时，表示该用户目前已有的默认模板总数量 (Optional) */
    TemplateCount int64 `json:"templateCount"`

    /* 模板列表 (Optional) */
    TemplateList []TemplateVo `json:"templateList"`
}
