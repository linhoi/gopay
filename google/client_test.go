package google

import (
	"context"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"google.golang.org/api/androidpublisher/v3"
	"gopkg.in/yaml.v2"
)

func TestClient_GetVoidedPurchase(t *testing.T) {
	tests := []struct {
		name                string
		wantVoidedPurchases []*androidpublisher.VoidedPurchase
		wantErr             bool
		startTime           time.Time
		endTime             time.Time
	}{
		{
			name:      "test",
			startTime: time.Now().AddDate(0, 0, -1),
			endTime:   time.Now(),
		},
	}

	c, pkg, err := testClient()
	if err != nil {
		fmt.Println("err", err)
		return
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNextPageToken := ""
			for {
				gotVoidedPurchases, gotNextPageToken, err := c.GetVoidedPurchase(context.Background(), VoidedPurchase{
					PackageName: pkg,
					StartTime:   &tt.startTime,
					EndTime:     &tt.endTime,
					Token:       gotNextPageToken,
				})

				if (err != nil) != tt.wantErr {
					t.Errorf("GetVoidedPurchase() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				for _, v := range gotVoidedPurchases {
					fmt.Println(v.OrderId)
				}
				if gotNextPageToken == "" {
					break
				}
			}
		})
	}
}

type conf struct {
	PackageName string `yaml:"package_name"`
	Credentials string `yaml:"credentials"`
}

func testClient() (*Client, string, error) {
	bd, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		return nil, "", err
	}
	var c conf
	err = yaml.Unmarshal(bd, &c)
	if err != nil {
		return nil, "", err
	}

	client, err := NewClient([]byte(c.Credentials))
	if err != nil {
		return nil, "", err
	}

	return client, c.PackageName, nil
}
