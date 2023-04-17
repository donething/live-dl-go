# sites

# 使用

通过 `plats.GenAnchor` 生成对应平台的主播实例

```go
anchorSite, err := plats.GenAnchor(&anchor)
if err != nil {
    return err
}
```

# 增加网站

新增网站需要：

1. 实现 `IAnchor` 接口的方法
2. 修改 `plats.Plats` 、 `plats.GenAnchor()`，增加新平台