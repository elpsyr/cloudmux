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

package cucloudcfel

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"moul.io/http2curl/v2"
	"yunion.io/x/jsonutils"
	"yunion.io/x/log"
	"yunion.io/x/pkg/errors"
	"yunion.io/x/pkg/gotypes"
	"yunion.io/x/pkg/util/httputils"
	"yunion.io/x/s3cli"

	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
)

const (
	CLOUD_PROVIDER_CUCLOUD_CN = "联通云"
	CLOUD_PROVIDER_CUCLOUD    = "ChinaUnionCfel"
	CUCLOUD_DEFAULT_REGION    = "cn-langfang-2"
)

type ChinaUnionClientConfig struct {
	cpcfg           cloudprovider.ProviderConfig
	accessKeyId     string
	accessKeySecret string

	mainUserName string
	userName     string
	password     string

	debug bool
}

type SChinaUnionClient struct {
	*ChinaUnionClientConfig

	client *http.Client
	lock   sync.Mutex
	ctx    context.Context

	regions   []SRegion
	ownerId   string
	token     string
	userId    int64
	accountId int64
}

func NewChinaUnionClientConfig(accessKeyId, accessKeySecret string) *ChinaUnionClientConfig {
	cfg := &ChinaUnionClientConfig{
		accessKeyId:     accessKeyId,
		accessKeySecret: accessKeySecret,
	}
	return cfg
}

func (self *ChinaUnionClientConfig) Debug(debug bool) *ChinaUnionClientConfig {
	self.debug = debug
	return self
}

func (self *ChinaUnionClientConfig) CloudproviderConfig(cpcfg cloudprovider.ProviderConfig) *ChinaUnionClientConfig {
	self.cpcfg = cpcfg
	return self
}

func NewChinaUnionClient(cfg *ChinaUnionClientConfig) (*SChinaUnionClient, error) {
	client := &SChinaUnionClient{
		ChinaUnionClientConfig: cfg,
		ctx:                    context.Background(),
	}
	client.ctx = context.WithValue(client.ctx, "time", time.Now())
	var err error
	client.regions, err = client.GetRegions()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (self *SChinaUnionClient) GetRegions() ([]SRegion, error) {
	if len(self.regions) > 0 {
		return self.regions, nil
	}
	resp, err := self.list("/instance/v1/product/cloudregions", nil)
	if err != nil {
		return nil, err
	}
	var ret []SRegion
	err = resp.Unmarshal(&ret, "list")
	if err != nil {
		return nil, err
	}
	self.regions = []SRegion{}
	for i := range ret {
		ret[i].client = self
		self.regions = append(self.regions, ret[i])
	}
	return self.regions, nil
}

func (self *SChinaUnionClient) GetRegion(id string) (*SRegion, error) {
	regions, err := self.GetRegions()
	if err != nil {
		return nil, err
	}
	for i := range regions {
		if regions[i].GetId() == id || regions[i].GetGlobalId() == id {
			regions[i].client = self
			return &regions[i], nil
		}
	}
	return nil, cloudprovider.ErrNotFound
}

func (self *SChinaUnionClient) getUrl(resource string) string {
	return fmt.Sprintf("https://gateway.cucloud.cn/%s", strings.TrimPrefix(resource, "/"))
}

func (self *SChinaUnionClient) SetToken(token *AccountInfo) {
	self.token = token.Token
	self.userId = token.UserId
	self.accountId = token.AccountId
}

type loginResp struct {
	Data struct {
		MainUserId    int64  `json:"mainUserId"`
		LoginUserID   int64  `json:"loginUserId"`
		LoginUserName string `json:"loginUserName"`
		PublicKey     string `json:"publicKey"`
		Token         string `json:"token"`
		UserType      int64  `json:"userType"`
	} `json:"data"`
	Msg    string `json:"msg"`
	Status string `json:"status"`
}

type AccountInfo struct {
	Token     string
	UserId    int64
	AccountId int64
	Timestamp time.Time
}

func (self *SChinaUnionClient) CheckAuth() error {

	res, err := self.getWithToken("iam/iam-portal/uc/v1/checkAuth", nil)
	if err != nil {
		return err
	}
	var ret loginResp
	if err = res.Unmarshal(&ret); err != nil {
		return err
	}
	if ret.Status != "200" {
		return fmt.Errorf(ret.Msg)
	}
	return nil
}

func (self *SChinaUnionClient) Login() (*AccountInfo, error) {
	// 测试先读本地
	var err error
	var tokenfile string
	if self.debug {
		tokenfile = "/root/project/puhui-uci/_output/token"
		tt, err := ioutil.ReadFile(tokenfile)
		if err != nil {
			log.Warningf("write token file err:%v", err)
		} else {
			arr := strings.Split(string(tt), ",")
			if len(arr) == 3 {
				userId, _ := strconv.Atoi(arr[1])
				accountId, _ := strconv.Atoi(arr[2])
				acc := &AccountInfo{Token: arr[0], UserId: int64(userId), AccountId: int64(accountId), Timestamp: time.Now()}
				self.SetToken(acc)
				if err = self.CheckAuth(); err == nil {
					return acc, err
				}
			}

		}
	}

	// iam用户登录 用iamUserName字段加密
	// 	https://gateway.cucloud.cn/iam/iam-portal/uc/v1/iam/login
	// {
	//   "userName": "zhaeng",
	//   "iamUserName": "cfel01",
	//   "password": "boEkdA7MWWBkeIJA4GWHrA==",
	//   "currentTimeMillis": "1735610890333"
	// }
	// 主账号登录用这个url iam/iam-portal/uc/v1/portal/login

	var timestamp = fmt.Sprintf("%v", time.Now().Unix())
	self.mainUserName = "zhaeng"
	self.userName = "cfel01" // iam 用户
	self.password = "Cfel@12345678"
	pwd, err := encryptPwd(self.userName, self.password, timestamp)
	if err != nil {
		return nil, err
	}
	params := map[string]interface{}{
		"userName": self.mainUserName, // 主账号
		"password": pwd,
		// "loginMode":         "0", //主账号登录需要
		"currentTimeMillis": timestamp,
		"iamUserName":       self.userName, // 不是 iam用户登录不需要
	}
	res, err := self.postWithToken("iam/iam-portal/uc/v1/iam/login", params)
	if err != nil {
		return nil, err
	}
	var ret loginResp
	if err = res.Unmarshal(&ret); err != nil {
		return nil, err
	}
	if ret.Status != "200" {
		return nil, fmt.Errorf(ret.Msg)
	}
	// 本地使用
	if self.debug {
		ioutil.WriteFile(tokenfile, []byte(fmt.Sprintf("%s,%d,%d", ret.Data.Token, ret.Data.LoginUserID, ret.Data.MainUserId)), os.ModeAppend)
	}

	account := &AccountInfo{Token: ret.Data.Token, UserId: ret.Data.LoginUserID, AccountId: ret.Data.MainUserId, Timestamp: time.Now()}
	self.SetToken(account)
	return account, nil
}

func (self *SChinaUnionClient) GetToken() string {
	return self.token
}

func (cli *SChinaUnionClient) getDefaultClient() *http.Client {
	cli.lock.Lock()
	defer cli.lock.Unlock()
	if !gotypes.IsNil(cli.client) {
		return cli.client
	}
	cli.client = httputils.GetTimeoutClient(60 * time.Second)
	httputils.SetClientProxyFunc(cli.client, cli.cpcfg.ProxyFunc)
	ts, _ := cli.client.Transport.(*http.Transport)
	ts.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	cli.client.Transport = cloudprovider.GetCheckTransport(ts, func(req *http.Request) (func(resp *http.Response) error, error) {
		if cli.cpcfg.ReadOnly {
			if req.Method == "GET" {
				return nil, nil
			}
			return nil, errors.Wrapf(cloudprovider.ErrAccountReadOnly, "%s %s", req.Method, req.URL.Path)
		}
		return nil, nil
	})
	return cli.client
}

type sChinaUnionError struct {
	StatusCode int `json:"statusCode"`
	Status     string
	Code       string
	Message    string
}

func (self *sChinaUnionError) Error() string {
	return jsonutils.Marshal(self).String()
}

func (self *sChinaUnionError) ParseErrorFromJsonResponse(statusCode int, status string, body jsonutils.JSONObject) error {
	if body != nil {
		body.Unmarshal(self)
	}
	self.StatusCode = statusCode
	return self
}

func (self *SChinaUnionClient) sign(req *http.Request) (string, error) {
	keys := []string{}
	keyMap := map[string]string{}
	for k := range req.Header {
		key, ok := map[string]string{
			"Accesskey":   "accessKey",
			"Algorithm":   "algorithm",
			"Requesttime": "requestTime",
		}[k]
		if ok {
			keys = append(keys, key)
			keyMap[key] = req.Header.Get(k)
		}
	}
	params, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		return "", errors.Wrapf(err, "ParseQuery")
	}
	for k := range params {
		keys = append(keys, k)
		keyMap[k] = params.Get(k)
	}
	if req.Method == string(httputils.POST) || req.Method == string(httputils.PUT) || req.Method == string(httputils.DELETE) {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return "", errors.Wrapf(err, "read body")
		}
		req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		obj, err := jsonutils.Parse(body)
		if err != nil {
			return "", errors.Wrapf(err, "params req body")
		}
		objMap, err := obj.GetMap()
		if err != nil {
			return "", errors.Wrapf(err, "req body map")
		}
		for k := range objMap {
			keys = append(keys, k)
			keyMap[k], _ = objMap[k].GetString()
		}
	}
	sort.Strings(keys)
	signStrs := []string{}
	for _, k := range keys {
		signStrs = append(signStrs, fmt.Sprintf(`%s="%s"`, k, keyMap[k]))
	}

	hasher := hmac.New(sha256.New, []byte(self.accessKeySecret))
	hasher.Write([]byte(strings.Join(signStrs, "&")))
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func (self *SChinaUnionClient) Do(req *http.Request) (*http.Response, error) {
	client := self.getDefaultClient()

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("algorithm", "HmacSHA256")
	req.Header.Set("requestTime", fmt.Sprintf("%d", time.Now().UTC().UnixMilli()))
	req.Header.Set("accessKey", self.accessKeyId)

	signature, err := self.sign(req)
	if err != nil {
		return nil, errors.Wrapf(err, "sign")
	}

	req.Header.Set("sign", signature)
	if self.debug {
		curlCmd, _ := http2curl.GetCurlCommand(req)
		fmt.Println(curlCmd)
	}
	return client.Do(req)
}

func (self *SChinaUnionClient) list(resource string, params map[string]interface{}) (jsonutils.JSONObject, error) {
	return self.request(httputils.GET, resource, params)
}

func (self *SChinaUnionClient) get(resource string, params map[string]interface{}) (jsonutils.JSONObject, error) {
	return self.request(httputils.GET, resource, params)
}

func (self *SChinaUnionClient) delete(resource string, params map[string]interface{}) error {
	_, err := self.request(httputils.DELETE, resource, params)
	return err
}

func (self *SChinaUnionClient) post(resource string, params map[string]interface{}) (jsonutils.JSONObject, error) {
	return self.request(httputils.POST, resource, params)
}

func (self *SChinaUnionClient) request(method httputils.THttpMethod, resource string, params map[string]interface{}) (jsonutils.JSONObject, error) {
	uri := self.getUrl(resource)
	if params == nil {
		params = map[string]interface{}{}
	}
	var body jsonutils.JSONObject = jsonutils.NewDict()
	switch method {
	case httputils.GET:
		values := url.Values{}
		for k, v := range params {
			values.Set(k, v.(string))
		}
		if len(values) > 0 {
			uri = fmt.Sprintf("%s?%s", uri, values.Encode())
		}
	case httputils.POST, httputils.PUT, httputils.DELETE:
		body = jsonutils.Marshal(params)
	}
	req := httputils.NewJsonRequest(method, uri, body)
	bErr := &sChinaUnionError{}
	client := httputils.NewJsonClient(self)
	_, resp, err := client.Send(self.ctx, req, bErr, self.debug)
	if err != nil {
		return nil, err
	}
	if gotypes.IsNil(resp) {
		return nil, fmt.Errorf("empty response")
	}
	code, _ := resp.GetString("code")
	if code != "200" {
		return nil, errors.Errorf(resp.String())
	}
	res, err := resp.GetIgnoreCases("result")
	if err != nil {
		return resp, nil
	}
	return res, nil
}

func (self *SChinaUnionClient) getWithToken(resource string, params map[string]interface{}) (jsonutils.JSONObject, error) {
	return self.requestWithToken(httputils.GET, resource, params)
}

func (self *SChinaUnionClient) postWithToken(resource string, params map[string]interface{}) (jsonutils.JSONObject, error) {
	return self.requestWithToken(httputils.POST, resource, params)
}

func (self *SChinaUnionClient) requestWithToken(method httputils.THttpMethod, resource string, params map[string]interface{}) (jsonutils.JSONObject, error) {
	uri := self.getUrl(resource)
	if params == nil {
		params = map[string]interface{}{}
	}
	var body jsonutils.JSONObject = jsonutils.NewDict()
	switch method {
	case httputils.GET:
		values := url.Values{}
		for k, v := range params {
			values.Set(k, v.(string))
		}
		if len(values) > 0 {
			uri = fmt.Sprintf("%s?%s", uri, values.Encode())
		}
	case httputils.POST:
		body = jsonutils.Marshal(params)
	}
	req := httputils.NewJsonRequest(method, uri, body)
	if self.token != "" {
		req.GetHeader().Add("access_token", self.token)
	}
	bErr := &sChinaUnionError{}
	client := httputils.NewJsonClient(self)
	_, resp, err := client.Send(self.ctx, req, bErr, self.debug)
	if err != nil {
		return nil, err
	}
	if gotypes.IsNil(resp) {
		return nil, fmt.Errorf("empty response")
	}
	code, err := resp.GetString("code")
	if errors.Cause(err) == jsonutils.ErrJsonDictKeyNotFound {
		code, err = resp.GetString("status")
		if errors.Cause(err) == jsonutils.ErrJsonDictKeyNotFound {
			return nil, errors.Errorf(resp.String())
		}
	}
	if code != "200" {
		return nil, errors.Errorf(resp.String())
	}
	return resp, nil
}

func (self *SChinaUnionClient) GetSubAccounts() ([]cloudprovider.SSubAccount, error) {
	subAccount := cloudprovider.SSubAccount{}
	subAccount.Id = self.GetAccountId()
	subAccount.Name = self.cpcfg.Name
	subAccount.Account = self.accessKeyId
	subAccount.HealthStatus = api.CLOUD_PROVIDER_HEALTH_NORMAL
	return []cloudprovider.SSubAccount{subAccount}, nil
}

func (self *SChinaUnionClient) getOwnerId() (string, error) {
	if len(self.ownerId) > 0 {
		return self.ownerId, nil
	}
	client, err := self.getS3Client()
	if err != nil {
		return "", err
	}
	buckets, err := client.ListBuckets()
	if err != nil {
		return "", err
	}
	self.ownerId = buckets.Owner.ID
	return self.ownerId, nil
}

func (self *SChinaUnionClient) getS3Client() (*s3cli.Client, error) {
	client, err := s3cli.New("obs-helf.cucloud.cn", self.accessKeyId, self.accessKeySecret, true, self.debug)

	tr := httputils.GetTransport(true)
	tr.Proxy = self.cpcfg.ProxyFunc
	return client, err
}

func (self *SChinaUnionClient) GetAccountId() string {
	ownerId, _ := self.getOwnerId()
	return ownerId
}

type CashBalance struct {
	CashBalance float64
}

// 接口不可用
func (self *SChinaUnionClient) QueryBalance() (*CashBalance, error) {
	ret := &CashBalance{}
	resp, err := self.post("bill-manage-console/bill/manage/balance/queryAvailableBalanceDetail", nil)
	if err != nil {
		return nil, err
	}
	err = resp.Unmarshal(ret)
	if err != nil {
		return nil, errors.Wrapf(err, "resp.Unmarshal")
	}
	return ret, nil
}

func (self *SChinaUnionClient) GetCapabilities() []string {
	caps := []string{
		cloudprovider.CLOUD_CAPABILITY_COMPUTE + cloudprovider.READ_ONLY_SUFFIX,
	}
	return caps
}
