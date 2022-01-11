//+build ignore

package paypal

import (
	"context"
	"reflect"
	"testing"
)

func TestClient_SetExpressCheckout(t *testing.T) {
	type args struct {
		req SetExpressCheckoutReq
	}
	tests := []struct {
		name    string
		args    args
		want    SetExpressCheckoutResp
		wantErr bool
	}{
		{
			name: "set express checkout",
			args: args{req: SetExpressCheckoutReq{
				ReturnURL:    "https://example.com",
				Amount:       "0.5",
				NoShipping:   "1",
				CurrencyCode: "USD",
				Desc:         "goodTitle0",
				Customer:     "myCustomer",
				Invoice:      "1242334323432442",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := testClient
			got, err := c.SetExpressCheckout(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetExpressCheckout() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Logf("SetExpressCheckout() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetExpressCheckoutDetails(t *testing.T) {
	type args struct {
		req GetExpressCheckoutDetailsReq
	}
	tests := []struct {
		name    string
		args    args
		want    GetExpressCheckoutDetailsResp
		wantErr bool
	}{
		{
			name: "test",
			args: args{req: GetExpressCheckoutDetailsReq{Token: "EC-8KN66271SK878335A"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := testClient
			got, err := c.GetExpressCheckoutDetails(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetExpressCheckoutDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Logf("GetExpressCheckoutDetails() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_DoExpressCheckoutPayment(t *testing.T) {
	type args struct {
		req GetExpressCheckoutDetailsReq
	}
	tests := []struct {
		name    string
		args    args
		want    GetExpressCheckoutDetailsResp
		wantErr bool
	}{
		{
			name: "test",
			args: args{req: GetExpressCheckoutDetailsReq{Token: "EC-3TX755711C115351G"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := testClient
			got, err := c.DoExpressCheckoutPayment(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoExpressCheckoutPayment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Logf("DoExpressCheckoutPayment() got = %v, want %v", got, tt.want)
			}
		})
	}
}