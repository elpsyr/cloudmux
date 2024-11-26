package aws

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type AttrValues struct {
	Value string
}

func (self *SRegion) GetAttributeValue(serviceCode string, attribute string) ([]AttrValues, error) {
	params := map[string]interface{}{
		"ServiceCode":   serviceCode,
		"AttributeName": attribute,
		"MaxResults":    100,
	}

	var nextToken string

	if len(nextToken) > 0 {
		params["NextToken"] = nextToken
	}

	var result []AttrValues

	for {
		ret := struct {
			FormatVersion   string
			NextToken       string `json:"NextToken"`
			AttributeValues []AttrValues
		}{}
		err := self.priceRequest("GetAttributeValues", params, &ret)
		// rr, _ := json.Marshal(ret)
		// fmt.Println(string(rr))
		if err != nil {
			continue
		}
		result = append(result, ret.AttributeValues...)
		if ret.NextToken == "" {
			break
		}
	}
	rr, _ := json.Marshal(result)
	ioutil.WriteFile("./"+attribute+".json", rr, os.ModeAppend)

	return result, nil
}
