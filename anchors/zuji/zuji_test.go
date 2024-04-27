package zuji

import (
	"github.com/donething/live-dl-go/anchors/base"
	"strings"
	"testing"
)

func TestAnchorZuji_GetAnchorInfo(t *testing.T) {
	type fields struct {
		Anchor *base.Anchor
	}

	tests := []struct {
		name    string
		fields  fields
		want    *base.AnchorInfo
		wantErr bool
	}{
		{
			name: "测试 蝴蝶曼",
			fields: fields{Anchor: &base.Anchor{
				UID:  "15050303",
				Plat: Platform,
			}},
			want: &base.AnchorInfo{
				Name: "蝴蝶曼",
			},
			wantErr: false,
		},
		{
			name: "测试 虎妮",
			fields: fields{Anchor: &base.Anchor{
				UID:  "29608771",
				Plat: Platform,
			}},
			want: &base.AnchorInfo{
				Name: "虎妮",
			},
			wantErr: false,
		},
		{
			name: "测试 铁锤",
			fields: fields{Anchor: &base.Anchor{
				UID:  "20233311",
				Plat: Platform,
			}},
			want: &base.AnchorInfo{
				Name: "铁锤",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AnchorZuji{
				Anchor: tt.fields.Anchor,
			}
			got, err := a.GetAnchorInfo()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAnchorInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 只比较 Name 属性，其它很多值经常变动，不便比较
			if !strings.Contains(got.Name, tt.want.Name) || (got.IsLive && got.StreamUrl == "") {
				t.Errorf("GetAnchorInfo() got = %v, want %v", got, tt.want)
			}
		})
	}
}
