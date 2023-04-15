# do-live-go

录制有限几个平台的直播

* 抖音
* 哔哩哔哩
* 足迹

# 使用

可查看`stream`包里的测试函数，查看使用方法

推荐使用的例子

```go
// 哔哩哔哩的用户 ID
func TestStartFlv(t *testing.T) {
	anchor := plats.Anchor{
		ID:   "8739477",
		Plat: plats.PlatBili,
	}

	err := StartFlvAnchor(&anchor, "D:/Tmp/live/bili_8739477.flv", 20*1024*1024, &tgHandler)
	if err != nil {
		t.Fatal(err)
	}
}

// 抖音的直播间号
func TestStartFlv2(t *testing.T) {
	anchor := plats.Anchor{
		ID:   "249406961231",
		Plat: plats.PlatDouyin,
	}

	err := StartFlvAnchor(&anchor, "D:/Tmp/live/douyin_249406961231.flv", 10*1024*1024, &tgHandler)
	if err != nil {
		t.Fatal(err)
	}
}

// 足迹的用户ID
func TestStartM3u8(t *testing.T) {
	anchor := plats.Anchor{
		ID:   "61667788",
		Plat: plats.PlatZuji,
	}

	err := StartM3u8Anchor(&anchor, "D:/Tmp/live/zuji_61667788.flv", 10*1024*1024, &tgHandler)
	if err != nil {
		t.Fatal(err)
	}
}
```