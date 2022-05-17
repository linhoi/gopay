package httpx

import (
	"context"
	"io/ioutil"
	"strings"
	"testing"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient()
			res, err := client.Get("https://baidu.com/robots.txt")
			if err != nil {
				if strings.Contains(err.Error(), context.DeadlineExceeded.Error()) {
					t.Log("is context error")
				}
				t.Error(err)
			}
			res.Body.Close()
			t.Log(res)

			byte, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Log(err)
			} else {
				t.Log(string(byte))
			}
			t.Log(res.StatusCode)
		})
	}
}
