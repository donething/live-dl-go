package bili

import (
	"dl-live-go/hanlders"
	"github.com/donething/utils-go/dotg"
	"os"
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
		want    *sites.AnchorInfo
		wantErr bool
	}{
		{
			name:    "DYS",
			args:    args{uid: "8739477"},
			wantErr: false,
			want: &sites.AnchorInfo{
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

func TestStartCapture(t *testing.T) {
	tgHandler := &hanlders.TGHandler{
		TG:        dotg.NewTGBot(os.Getenv("MY_TG_TOKEN")),
		LocalPort: 0,
		ChatID:    os.Getenv("MY_TG_CHAT_LIVE"),
	}

	err := StartCapture("8739477", "D:/Tmp/live/7777", 15*1024*1024, tgHandler)
	if err != nil {
		t.Fatal(err)
	}
}
