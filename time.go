package utils

import (
	"time"
)


/*
返回1970到今天多少天
*/
func UnixDay(times ...interface{}) int {
	if len(times) == 0 {
		return int((time.Now().Unix() + 8*3600) / 86400)
	}

	ttime := times[0]
	if t, ok := ttime.(time.Time); ok {
		return int((t.Unix() + 8*3600) / 86400)
	} else if t, ok := ttime.(int); ok {
		return int((time.Unix(int64(t), 0).Unix() + 8*3600) / 86400)
	} else if t, ok := ttime.(int64); ok {
		return int((time.Unix(t, 0).Local().Unix() + 8*3600) / 86400)
	}

	panic("Gday's param is not time.Time")
}

/*
将字符串类型的时间转换为time类型
@format string，格式定义可选，默认"2016-01-02 15:04:05"
*/
func Strttime(strtime string, format ...string) (time.Time, error) {
	f := "2006-01-02 15:04:05"
	if len(format) > 0 {
		f = format[0]
	}
	return time.ParseInLocation(f, strtime, time.Local)
}

/*
将time类型、int类型（Unix时间戳）转换为 字符串类型， 

@format string，格式定义可选，默认"2016-01-02 15:04:05"
*/
func Strftime(times ...interface{}) string {
	f := "2006-01-02 15:04:05"
	if len(times) == 0 {
		return time.Now().Local().Format(f)
	}

	ttime := times[0]

	if len(times) > 1 {
		f, _ = times[1].(string)
	}

	if t, ok := ttime.(time.Time); ok {
		return t.Format(f)
	} else if t, ok := ttime.(int); ok {
		return time.Unix(int64(t), 0).Format(f)
	} else if t, ok := ttime.(int64); ok {
		return time.Unix(t, 0).Format(f)
	}

	panic("Strftime's param is not time.Time or int、int64")
}

// return xxxx-xx-xx 00:00:00 的unix时间戳int
func TodayUnix(times ...interface{}) int {
	var tt time.Time
	if len(times) == 0 {
		tt = time.Now()
	} else {
		ttime := times[0]
		if t, ok := ttime.(time.Time); ok {
			tt = t
		} else if t, ok := ttime.(int); ok {
			tt = time.Unix(int64(t), 0)
		} else if t, ok := ttime.(int64); ok {
			tt = time.Unix(t, 0)
		}
	}

	tt = tt.Local()
	tt = time.Date(tt.Year(), tt.Month(), tt.Day(), 0, 0, 0, 0, time.Local)

	return int(tt.Unix())
}

// return yyyymmdd int格式
func Today(times ...interface{}) int {
	var tt time.Time
	if len(times) == 0 {
		tt = time.Now()
	} else {
		ttime := times[0]
		if t, ok := ttime.(time.Time); ok {
			tt = t
		} else if t, ok := ttime.(int); ok {
			tt = time.Unix(int64(t), 0)
		} else if t, ok := ttime.(int64); ok {
			tt = time.Unix(t, 0)
		}
	}

	tt = tt.Local()
	return tt.Year()*10000 + int(tt.Month())*100 + tt.Day()
}


// return yyyymmdd int格式
func Yesterday(times ...interface{}) int {
	var tt time.Time
	if len(times) == 0 {
		tt = time.Now()
	} else {
		ttime := times[0]
		if t, ok := ttime.(time.Time); ok {
			tt = t
		} else if t, ok := ttime.(int); ok {
			tt = time.Unix(int64(t), 0)
		} else if t, ok := ttime.(int64); ok {
			tt = time.Unix(t, 0)
		}
	}

	tt = tt.AddDate(0, 0, -1)
	tt = tt.Local()
	mm := int(tt.Month())
	return tt.Year()*10000 + mm*100 + tt.Day()

}

// 计算两个时间的天数差值
func Diffday(time1 int, time2 int) int {

	tt1 := time.Unix(int64(time1), 0)
	tt2 := time.Unix(int64(time2), 0)

	return tt1.YearDay() - tt2.YearDay()
}

