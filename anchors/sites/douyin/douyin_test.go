package douyin

import (
	"github.com/donething/live-dl-go/anchors/baseanchor"
	"strings"
	"testing"
)

func TestAnchorDouyin_GetAnchorInfo(t *testing.T) {
	type fields struct {
		Anchor *baseanchor.Anchor
	}
	tests := []struct {
		name    string
		fields  fields
		want    *baseanchor.AnchorInfo
		wantErr bool
	}{
		{
			name: "测试不在播 假树",
			fields: fields{Anchor: &baseanchor.Anchor{
				UID:  "165251594775",
				Plat: Platform,
			}},
			want: &baseanchor.AnchorInfo{
				Name: "假树",
			},
			wantErr: false,
		},
		{
			name: "测试不在播 ☁️云福晋☁️",
			fields: fields{Anchor: &baseanchor.Anchor{
				UID:  "483360313799",
				Plat: Platform,
			}},
			want: &baseanchor.AnchorInfo{
				Name: "☁️云福晋☁️",
			},
			wantErr: false,
		},
		{
			name: "测试在播 鹿哥电影",
			fields: fields{Anchor: &baseanchor.Anchor{
				UID:  "937098912324",
				Plat: Platform,
			}},
			want: &baseanchor.AnchorInfo{
				Name: "鹿哥电影",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AnchorDouyin{
				Anchor: tt.fields.Anchor,
			}
			got, err := a.GetAnchorInfo()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAnchorInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 只比较 Name 属性，其它很多值经常变动，不便比较
			if !strings.Contains(got.Name, tt.want.Name) {
				t.Errorf("GetAnchorInfo() got = %v, want %v", got, tt.want)
			}
		})
	}
}
