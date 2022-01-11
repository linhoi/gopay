//+build ignore

package paypal

import (
	"context"
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"
	"time"

	"github.com/linhoi/kit/log"
	"gopkg.in/yaml.v3"
)

func TestMain(m *testing.M) {
	// Config Just for test
	b, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		b, err = ioutil.ReadFile("config.temp.yaml")
		if err != nil {
			fmt.Println("please set you paypal config in config.temp.yaml")
			fmt.Println(err)
			return
		}

	}
	var config Config
	err = yaml.Unmarshal(b, &config)
	if err != nil {
		fmt.Println(err)
		return
	}

	testClient = NewClient(config)

	m.Run()
}

var testClient *Client

func TestClient_GetAccessToken(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name:    "get token success",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := testClient
			got, err := c.GetAccessToken(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAccessToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == "" {
				t.Errorf("GetAccessToken() got = %v, want %v", got, tt.want)
			}
			t.Logf("got token: %v", got)
		})
	}
}

func TestClient_ListDisputes(t *testing.T) {
	tests := []struct {
		name             string
		wantDisputeItems []DisputeItem
		wantErr          bool
	}{
		{
			name:    "test with success",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := testClient
			gotDisputeItems, err := c.ListDisputes(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("ListDisputes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDisputeItems, tt.wantDisputeItems) {
				t.Logf("ListDisputes() gotDisputeItems = %v, want %v", gotDisputeItems, tt.wantDisputeItems)
			}
		})
	}
}

func TestClient_ListOnePageDisputes(t *testing.T) {
	tests := []struct {
		name             string
		wantDisputeItems []DisputeItem
		wantErr          bool
	}{
		{
			name:    "test with success",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := testClient
			gotDisputeItems, _, err := c.ListOnePageDisputes(context.Background(), "")
			if (err != nil) != tt.wantErr {
				t.Errorf("ListDisputes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDisputeItems, tt.wantDisputeItems) {
				t.Logf("ListDisputes() gotDisputeItems = %v, want %v", gotDisputeItems, tt.wantDisputeItems)
			}
		})
	}
}

func TestClient_ShowDisputeDetails(t *testing.T) {
	type args struct {
		disputeID string
	}
	tests := []struct {
		name       string
		args       args
		wantDetail ShowDisputeDetailsResp
		wantErr    bool
	}{
		{
			name:    "fake dispute",
			args:    args{disputeID: "2311AD"},
			wantErr: true,
		},
		{
			name:    "true",
			args:    args{disputeID: "PP-010-233-987-341"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := testClient
			gotDetail, err := c.ShowDisputeDetails(context.Background(), tt.args.disputeID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ShowDisputeDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDetail, tt.wantDetail) {
				t.Logf("ShowDisputeDetails() gotDetail = %v, want %v", gotDetail, tt.wantDetail)
			}
		})
	}
}

func TestClient_GetTransaction(t *testing.T) {
	type args struct {
		transactionID string
		startTime     time.Time
		endTime       time.Time
	}
	tests := []struct {
		name            string
		args            args
		wantTransaction TransactionInfo
		wantErr         bool
	}{
		{
			name:    "fake transaction id",
			args:    args{transactionID: "awrffd", startTime: time.Now().Add(-30 * 24 * time.Hour), endTime: time.Now()},
			wantErr: true,
		},
		{
			name:    "true transaction id",
			args:    args{transactionID: "73C23547SY247022R", startTime: time.Now().Add(-30 * 24 * time.Hour), endTime: time.Now()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := testClient
			gotTransaction, err := c.GetTransaction(context.Background(), tt.args.transactionID, tt.args.startTime, tt.args.endTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotTransaction, tt.wantTransaction) {
				t.Logf("GetTransaction() gotTransaction = %v, want %v", gotTransaction, tt.wantTransaction)
			}
		})
	}
}

func TestClient_GetRefundTransaction(t *testing.T) {
	type args struct {
		startTime time.Time
		endTime   time.Time
	}
	tests := []struct {
		name             string
		args             args
		wantTransactions []TransactionInfo
		wantErr          bool
	}{
		{
			name: "test with success",
			args: args{
				startTime: time.Now().Add(-30 * 24 * time.Hour),
				endTime:   time.Now(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := testClient
			gotTransactions, err := c.GetRefundTransaction(context.Background(), tt.args.startTime, tt.args.endTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRefundTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotTransactions, tt.wantTransactions) {
				t.Logf("GetRefundTransaction() gotTransactions = %v, want %v", gotTransactions, tt.wantTransactions)
			}
		})
	}
}

func TestClient_ListOnePageDisputesV2(t *testing.T) {
	tests := []struct {
		name             string
		wantDisputeItems []DisputeItem
		wantErr          bool
	}{
		{
			name:    "test with success",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := testClient
			gotDisputeItems, next, err := c.ListOnePageDisputesV2(context.Background(), "", time.Now().Add(-20*24*time.Hour), time.Now().Add(-8*time.Hour))
			if (err != nil) != tt.wantErr {
				t.Errorf("ListDisputes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("next %s", next)
			if !reflect.DeepEqual(gotDisputeItems, tt.wantDisputeItems) {
				t.Logf("ListDisputes() gotDisputeItems = %v, want %v", gotDisputeItems, tt.wantDisputeItems)
			}
		})
	}
}

func TestClient_ListOnePageDisputesV3(t *testing.T) {
	err := syncPaypalRefund(context.Background(), time.Now().Add(time.Duration(-173)*time.Hour*24), 3)
	if err != nil {
		fmt.Println(err)
	}

}

func syncPaypalRefund(ctx context.Context, endTime time.Time, stepDay int) error {
	client := testClient
	pageToken := ""
	for {
		allDisputes, nextPageToken, err := client.ListOnePageDisputesV2(ctx, pageToken, endTime.Add(time.Duration(-stepDay)*24*time.Hour), endTime)
		if err != nil {
			log.S(ctx).Errorw("ListDisputes", "err",err)

			return err
		}

		for _, d := range allDisputes {
			fmt.Println(d)
		}

		userRefundDisputeDetailByOrderNumber, noOrderNoUserRefundDisputeDetails := client.GetUserRefundDisputedDetails(ctx, allDisputes)

		fmt.Println("userRefundDisputeDetailByOrderNumber", userRefundDisputeDetailByOrderNumber)
		fmt.Println("noOrderNoUserRefundDisputeDetails", noOrderNoUserRefundDisputeDetails)
		if nextPageToken == "" {
			break
		}

		pageToken = nextPageToken
	}

	return nil

}
