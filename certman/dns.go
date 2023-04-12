package certman

import (
	"context"
	"fmt"

	"github.com/libdns/libdns"
	"github.com/libdns/tencentcloud"
)

func createRecord(ctx context.Context, zone, value string) error {

	provider := &tencentcloud.Provider{
		SecretId:  "YOUR_Secret_ID",
		SecretKey: "YOUR_Secret_Key",
	}

	record, err := provider.SetRecords(ctx, zone, []libdns.Record{
		{
			Type:  "TXT",
			Name:  "_acme-challenge",
			Value: value,
		},
	})

	fmt.Println(record, err)

	return err

}
