package bili

import (
	"github.com/donething/live-dl-go/sites/entity"
	"testing"
)

func TestAnchorBili_GetAnchorInfo(t *testing.T) {
	type fields struct {
		Anchor *entity.Anchor
	}
	tests := []struct {
		name    string
		fields  fields
		want    *entity.AnchorInfo
		wantErr bool
	}{
		{
			name: "测试 DYS",
			fields: fields{Anchor: &entity.Anchor{
				ID:   "8739477",
				Plat: Plat,
			}},
			want: &entity.AnchorInfo{
				Name: "老实憨厚的笑笑",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AnchorBili{
				Anchor: tt.fields.Anchor,
			}
			got, err := a.GetAnchorInfo()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAnchorInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 只比较 Name 属性，其它很多值经常变动，不便比较
			if got.Name != tt.want.Name {
				t.Errorf("GetAnchorInfo() got = %v, want %v", got, tt.want)
			}
		})
	}
}
