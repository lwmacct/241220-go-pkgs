package mtime

import "time"

type this struct {
	location *time.Location
}

func New() *this {
	return &this{
		location: time.FixedZone("CST", 8*3600), // CST (China Standard Time) UTC+8
	}
}

// 将时间向前取整到最接近的 5 的倍数
func (t *this) Round5m(d time.Time) time.Time {
	// 获取分钟数
	minutes := d.Minute()
	// 向前取整到最接近的 5 的倍数
	roundedMinutes := (minutes / 5) * 5
	// 返回新时间
	return time.Date(d.Year(), d.Month(), d.Day(), d.Hour(), roundedMinutes, 0, 0, t.location)
}

func (t *this) ToString(d time.Time, format string) string {
	return d.In(t.location).Format(format)
}

func (t *this) ToUnix(d time.Time) int64 {
	return d.In(t.location).Unix()
}

func (t *this) ToTime(d int64) time.Time {
	return time.Unix(d, 0).In(t.location)
}

func (t *this) ToCST(d time.Time) time.Time {
	return d.In(t.location)
}
