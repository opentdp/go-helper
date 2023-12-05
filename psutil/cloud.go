package psutil

import (
	"strings"

	"github.com/opentdp/go-helper/request"
)

// 云实例 Id

func CloudInstanceId() string {

	var url string
	var mid string

	// alibaba
	url = "http://100.100.100.200/latest/meta-data/instance-id"
	mid = request.TimingGet(url, request.H{}, 3)
	if mid != "" {
		return strings.TrimSpace(mid)
	}

	// tencent
	url = "http://metadata.tencentyun.com/latest/meta-data/instance-id"
	mid = request.TimingGet(url, request.H{}, 3)
	if mid != "" {
		return strings.TrimSpace(mid)
	}

	// aws baidu huawei
	url = "http://169.254.169.254/latest/meta-data/instance-id"
	mid = request.TimingGet(url, request.H{}, 3)
	if mid != "" {
		return strings.TrimSpace(mid)
	}

	// azure
	url = "http://169.254.169.254/metadata/instance/compute/vmId?api-version=2021-01-01"
	mid = request.TimingGet(url, request.H{"Metadata": "true"}, 3)
	if mid != "" {
		return strings.TrimSpace(mid)
	}

	// google
	url = "http://metadata.google.internal/computeMetadata/v1/instance/id"
	mid = request.TimingGet(url, request.H{"Metadata-Flavor": "Google"}, 3)
	if mid != "" {
		return strings.TrimSpace(mid)
	}

	// digitalocean
	url = "http://169.254.169.254/metadata/v1/id"
	mid = request.TimingGet(url, request.H{}, 3)
	if mid != "" {
		return strings.TrimSpace(mid)
	}

	return mid

}
