package aws

import (
	"fmt"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/pkg/errors"
)

// Verify that *SRegion implements ICfelCloudRegion
var _ cloudprovider.ICfelCloudRegion = (*SRegion)(nil)

func (self *SRegion) GetInstanceMatchImage(instanceType string) ([]cloudprovider.ICloudImage, error) {
	sku, err := self.CfelGetInstanceType(instanceType)
	if err != nil {
		return nil, err
	}
	images, err := self.cfelGetImages("", []TImageOwnerType{ImageOwnerTypeSystem}, nil, "", "", nil, "", sku.ProcessorInfo.SupportedArchitectures)
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

func (self *SRegion) CfelGetInstanceType(name string) (*Sku, error) {
	params := map[string]string{
		"InstanceType.1": name,
	}
	ret := struct {
		InstanceTypeSet []Sku  `xml:"instanceTypeSet>item"`
		NextToken       string `xml:"nextToken"`
	}{}
	err := self.ec2Request("DescribeInstanceTypes", params, &ret)
	if err != nil {
		return nil, err
	}
	for i := range ret.InstanceTypeSet {
		if ret.InstanceTypeSet[i].InstanceType == name {
			return &ret.InstanceTypeSet[i], nil
		}
	}
	return nil, errors.Wrapf(cloudprovider.ErrNotFound, name)
}

func (self *SRegion) GetICfelCloudImage(withUserMeta bool) ([]cloudprovider.ICloudImage, error) {
	return nil, cloudprovider.ErrNotImplemented
}

func (self *SRegion) GetICfelCloudImageById(id string) (cloudprovider.ICloudImage, error) {
	return nil, nil
}

func (self *SRegion) cfelGetImages(status ImageStatusType, owners []TImageOwnerType, imageId []string, name string, virtualizationType string, ownerIds []string, volumeType string, arch []string) ([]SImage, error) {
	params := map[string]string{}
	idx := 1

	if len(status) > 0 {
		params[fmt.Sprintf("Filter.%d.Name", idx)] = "state"
		params[fmt.Sprintf("Filter.%d.Value.1", idx)] = string(status)
		idx++
	}

	if len(name) > 0 {
		params[fmt.Sprintf("Filter.%d.Name", idx)] = "name"
		params[fmt.Sprintf("Filter.%d.Value.1", idx)] = name
		idx++
	}

	if len(virtualizationType) > 0 {
		params[fmt.Sprintf("Filter.%d.Name", idx)] = "virtualization-type"
		params[fmt.Sprintf("Filter.%d.Value.1", idx)] = virtualizationType
		idx++
	}

	if len(volumeType) > 0 {
		params[fmt.Sprintf("Filter.%d.Name", idx)] = "block-device-mapping.volume-type"
		params[fmt.Sprintf("Filter.%d.Value.1", idx)] = volumeType
		idx++
	}

	if len(arch) > 0 {
		params[fmt.Sprintf("Filter.%d.Name", idx)] = "architecture"
		for i, v := range arch {
			params[fmt.Sprintf("Filter.%d.Value.%d", idx, i+1)] = v
		}
		idx++
	}
	
	if len(owners) > 0 || len(ownerIds) > 0 {
		for i, owner := range imageOwnerTypes2Strings(owners, ownerIds) {
			params[fmt.Sprintf("Owner.%d", i+1)] = string(owner)
		}
	}

	for i, id := range imageId {
		params[fmt.Sprintf("ImageId.%d", i+1)] = id
	}

	params[fmt.Sprintf("Filter.%d.Name", idx)] = "image-type"
	params[fmt.Sprintf("Filter.%d.Value.1", idx)] = "machine"
	idx++

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

	noVersionImages := make([]SImage, 0)
	versionedImages := make(map[string][]SImage)
	for i := range ret {
		key := fmt.Sprintf("%s%s", getImageOSDist(ret[i]), getImageOSVersion(ret[i]))
		if len(key) == 0 {
			noVersionImages = append(noVersionImages, ret[i])
			continue
		}
		if _, ok := versionedImages[key]; !ok {
			versionedImages[key] = make([]SImage, 0)
		}
		versionedImages[key] = append(versionedImages[key], ret[i])
	}
	for key := range versionedImages {
		noVersionImages = append(noVersionImages, getLatestImage(versionedImages[key]))
	}
	return noVersionImages, nil
}
