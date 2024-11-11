package aws

import (
	"fmt"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/pkg/errors"
)

// Verify that *SRegion implements ICfelCloudRegion
var _ cloudprovider.ICfelCloudRegion = (*SRegion)(nil)

func (self *SRegion) GetInstanceMatchImage(string) ([]cloudprovider.ICloudImage, error) {
	return self.GetICfelCloudImage(false)
}

func (self *SRegion) GetICfelCloudImage(withUserMeta bool) ([]cloudprovider.ICloudImage, error) {
	images, err := self.getPublicImages(ImageOwnerSystem, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "GetImages")
	}
	ret := []cloudprovider.ICloudImage{}
	for i := 0; i < len(images); i += 1 {
		// images[i].storageCache = self
		ret = append(ret, &images[i])
	}
	return ret, nil
}

func (self *SRegion) GetICfelCloudImageById(id string) (cloudprovider.ICloudImage, error) {
	return nil,nil
}

func (self *SRegion) getPublicImages(owners []TImageOwnerType, ownerIds []string) ([]SImage, error) {
	params := map[string]string{}
	idx := 1

	params[fmt.Sprintf("Filter.%d.Name", idx)] = "is-public"
	params[fmt.Sprintf("Filter.%d.Value.1", idx)] = "true"
	idx++

	params[fmt.Sprintf("Filter.%d.Name", idx)] = "image-type"
	params[fmt.Sprintf("Filter.%d.Value.1", idx)] = "machine"
	idx++

	params[fmt.Sprintf("Filter.%d.Name", idx)] = "architecture"
	params[fmt.Sprintf("Filter.%d.Value.1", idx)] = "x86_64"
	idx++

	// if len(owners) > 0 || len(ownerIds) > 0 {
	// 	for i, owner := range imageOwnerTypes2Strings(owners, ownerIds) {
	// 		params[fmt.Sprintf("Owner.%d", i+1)] = string(owner)
	// 	}
	// }
	params[fmt.Sprintf("Owner.%d", 1)] = "amazon"

	ret := []SImage{}
	for {
		part := struct {
			ImagesSet []SImage `xml:"imagesSet>item"`
			NextToken string   `xml:"nextToken"`
		}{}
		err := self.ec2Request("DescribeImages", params, &part)
		if err != nil {
			return nil, errors.Wrapf(err, "DescribeImages")
		}
		ret = append(ret, part.ImagesSet...)

		if len(part.ImagesSet) == 0 || len(part.NextToken) == 0 {
			break
		}
		params["NextToken"] = part.NextToken
	}
	return ret, nil
}
