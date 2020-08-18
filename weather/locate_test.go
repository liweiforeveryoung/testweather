package weather

import "testing"

func Test_getLocateImformationFromGaode(t *testing.T) {
	type args struct {
		ip string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "test1",
			args:    args{ip: "110.247.50.2"},
			want:    "131100",
			wantErr: false,
		},
		{
			name:    "test2",
			args:    args{ip: "127.0.0.1"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getAdcodeFromIP(tt.args.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("getAdcodeFromIP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getAdcodeFromIP() got = %v, want %v", got, tt.want)
			}
		})
	}
}
