package aws

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type ServiceCode struct {
	AttributeNames []string
	ServiceCode    string
}

func (self *SRegion) DescirbeServices(serviceCode string) ([]SInstanceType, error) {

	var result []ServiceCode
	for {
		ret := struct {
			FormatVersion string
			NextToken     string `json:"NextToken"`
			Services      []ServiceCode
		}{}
		params := map[string]interface{}{}
		if len(ret.NextToken) > 0 {
			params["NextToken"] = ret.NextToken
		}
		err := self.priceRequest("DescribeServices", params, &ret)
		if len(ret.Services) > 0 {
			result = append(result, ret.Services...)
		}
		if ret.NextToken == "" {
			break
		}
		if err != nil {
			continue
		}
	}
	rr, _ := json.Marshal(result)
	// fmt.Println(string(rr))
	if len(rr) > 0 {
		ioutil.WriteFile("./serviceCode.json", rr, os.ModeAppend)
	}

	return nil, nil
}
