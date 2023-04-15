package m3u8decoder

import (
	"testing"
)

func Test_getPrefix(t *testing.T) {
	type args struct {
		m3u8Url    string
		segmentUrl string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "完整的http",
			args: args{
				m3u8Url:    "https://host.com/path/video.m3u8?key=abc",
				segmentUrl: "https://host.com/path/1.ts?key=abc",
			},
			want: "",
		},
		{
			name: "以'/'开头",
			args: args{
				m3u8Url:    "https://host.com/path/video.m3u8?key=abc",
				segmentUrl: "/path/1.ts?key=abc",
			},
			want: "https://host.com/path",
		},
		{
			name: "以文件名开头",
			args: args{
				m3u8Url:    "https://host.com/path/video.m3u8?key=abc",
				segmentUrl: "1.ts?key=abc",
			},
			want: "https://host.com/path/",
		},
		{
			name: "以目录名开头",
			args: args{
				m3u8Url:    "https://host.com/path/video.m3u8?key=abc",
				segmentUrl: "live/1.ts?key=abc",
			},
			want: "https://host.com/path/live/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getPrefix(tt.args.m3u8Url, tt.args.segmentUrl); got != tt.want {
				t.Errorf("getPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestM3u8Decoder_Decode(t *testing.T) {
	m3u8 := M3u8Decoder{}
	err := m3u8.Decode("http://bjlive.szsbtech.com/record/b9Jbf4epv93cP9l.m3u8?"+
		"auth_key=1681409489-0-0-d74c14681a199e6dacac732253e0c14a", nil)
	if err != nil {
		t.Fatal(err)
	}

	for _, segment := range m3u8.Segments {
		t.Log(segment)
	}
}
