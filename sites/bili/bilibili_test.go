package bili

import (
	"live-dl-go/sites/plats"
	"reflect"
	"testing"
)

func TestGetBiliLiveUrl(t *testing.T) {
	type args struct {
		uid string
	}
	tests := []struct {
		name    string
		args    args
		want    *plats.AnchorInfo
		wantErr bool
	}{
		{
			name:    "DYS",
			args:    args{uid: "8739477"},
			wantErr: false,
			want: &plats.AnchorInfo{
				Name:  "老实憨厚的笑笑",
				Title: "德云色 5点解说比赛！！",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAnchorInfo(tt.args.uid)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAnchorInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// StreamUrl 是变动的，跳过比较该项
			tt.want.IsLive = got.IsLive
			tt.want.StreamUrl = got.StreamUrl

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAnchorInfo() got = %v, want %v", got, tt.want)
			}
		})
	}
}
