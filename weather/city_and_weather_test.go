package weather

import (
	"reflect"
	"testing"
)

func Test_getCityAndWeatherImfor(t *testing.T) {
	type args struct {
		adcode string
	}
	tests := []struct {
		name    string
		args    args
		want    *CityAndWeather
		wantErr bool
	}{
		//{"status":"1","count":"1","info":"OK",
		//"infocode":"10000","lives":[{"province":"河北","city":"衡水市",
		//"adcode":"131100","weather":"阴","temperature":"24",
		//"winddirection":"东北","windpower":"≤3","humidity":"90",
		//"reporttime":"2020-08-16 19:55:23"}]}
		{
			name: "test1",
			args: args{adcode: "131100"},
			want: &CityAndWeather{
				Status: "1",
				Info:   "OK",
				Count:  1,
				Lives: []Live{{
					City:        "衡水市",
					Weather:     "阴",
					Temperature: 24,
					ReportTime:  "2020-08-16 19:55:23",
				}},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getCityAndWeatherImforFromGaode(tt.args.adcode)
			if (err != nil) != tt.wantErr {
				t.Errorf("getCityAndWeatherImforFromGaode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getCityAndWeatherImforFromGaode() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_addWeatherImforToRedis(t *testing.T) {
	type args struct {
		adcode       string
		weatherImfor *WeatherImformationResponse
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test add weather",
			args: args{
				adcode: "131100",
				weatherImfor: &WeatherImformationResponse{
					Status:      0,
					City:        "衡水市",
					Weather:     "多云",
					Temperature: 24,
					ReportTime:  "2020-08-17 08:25:24",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func Test_getWeatherImforFromRedis(t *testing.T) {
	type args struct {
		adcode string
	}
	tests := []struct {
		name string
		args args
		want *WeatherImformationResponse
	}{
		{
			name: "Test_getWeatherImforFromRedis",
			args: args{adcode: "131100"},
			want: &WeatherImformationResponse{
				Status:      0,
				City:        "衡水市",
				Weather:     "多云",
				Temperature: 24,
				ReportTime:  "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getWeatherImforFromRedis(tt.args.adcode); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getWeatherImforFromRedis() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getWeatherImformationResponseFromAdcode(t *testing.T) {
	type args struct {
		adcode string
	}
	tests := []struct {
		name string
		args args
		want *WeatherImformationResponse
	}{
		{
			name: "Test_getWeatherImformationResponseFromAdcode",
			args: args{adcode: "131100"},
			want: &WeatherImformationResponse{
				Status:      0,
				City:        "衡水市",
				Weather:     "多云",
				Temperature: 24,
				ReportTime:  "",
			},
		},
		{
			name: "Test_getWeatherImformationResponseFromAdcode2",
			args: args{adcode: "110101"},
			want: &WeatherImformationResponse{
				Status:      0,
				City:        "东城区",
				Weather:     "阴",
				Temperature: 25,
				ReportTime:  "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getWeatherImformationResponseFromAdcode(tt.args.adcode); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getWeatherImformationResponseFromAdcode() = %v, want %v", got, tt.want)
			}
		})
	}
}
