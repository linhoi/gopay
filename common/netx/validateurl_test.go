package netx

import "testing"

func TestValidateUrl(t *testing.T) {
	type args struct {
		urlStr string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "a invalid url",
			args:    args{urlStr: "xxxx53rfdadscsrcsaqw"},
			wantErr: true,
		},
		{
			name:    "a valid url unreachable",
			args:    args{urlStr: "https://githubabcd.com/efg"},
			wantErr: true,
		},
		{
			name:    "a valid url reachable",
			args:    args{urlStr: "https://www.baidu.com/s?tn=baidutop10&rsv_idx=2"},
			wantErr: false,
		},
		//{
		//	name:    "a valid address",
		//	args:    args{urlStr: "http://127.0.0.1:1087"},
		//	wantErr: false,
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUrl(tt.args.urlStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUrl() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Log(err)
		})
	}
}
