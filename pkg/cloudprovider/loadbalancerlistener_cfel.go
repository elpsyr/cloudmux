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

package cloudprovider

type SCfelLoadbalancerListenerRule struct {
	Name             string
	Domain           string
	Path             string
	BackendGroupId   string
	BackendGroupType string
	ListenerId       string

	HttpRequestRate       int
	HttpRequestRatePerSrc int

	Redirect       string
	RedirectCode   int
	RedirectScheme string
	RedirectHost   string
	RedirectPath   string
}

type SCfelLoadbalancerBackendGroup struct {
	Name           string
	GroupType      string
	Backends       []SLoadbalancerBackend
	LoadBalancerId string
}
