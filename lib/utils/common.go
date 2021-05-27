package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/lithammer/shortuuid/v3"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	tnet "github.com/toolkits/net"
)

// JSON alias type
type JSON = map[string]interface{}

//ServiceError ...
type ServiceError struct {
	Code int
	Err  error
}

//FieldTran ...
type FieldTran map[string]string

var (
	once     sync.Once
	clientIP = "127.0.0.1"
)

// XRequestID
var XRequestID = "X-Request-ID"

// GetLocalIP
func GetLocalIP() string {
	once.Do(func() {
		ips, _ := tnet.IntranetIP()
		if len(ips) > 0 {
			clientIP = ips[0]
		} else {
			clientIP = "127.0.0.1"
		}
	})
	return clientIP
}

// GenRequestID eg: 76d27e8c-a80e-48c8-ad20-e5562e0f67e4
func GenRequestID() string {
	return UUID()
}

//UUID ...
func UUID() string {
	return shortuuid.NewWithNamespace("")
}

//MapToJSON ...
func MapToJSON(m interface{}) ([]byte, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return b, nil
}

//JSONToMap ...
func JSONToMap(json string) (interface{}, error) {
	if !gjson.Valid(json) {
		return nil, errors.New("invalid json")
	}

	m := gjson.Parse(json).Value()

	return m, nil
}

//YearMonthToDay 获取xx月有多少天
func YearMonthToDay(date string) int {
	tmp := strings.Split(date, "-")
	year := cast.ToInt(tmp[0])
	month := cast.ToInt(tmp[1])
	var days int
	if month != 2 {
		if month == 4 || month == 6 || month == 9 || month == 11 {
			days = 30

		} else {
			days = 31
			//fmt.Fprintln(os.Stdout, "The month has 31 days")
		}
	} else {
		if ((year%4) == 0 && (year%100) != 0) || (year%400) == 0 {
			days = 29
		} else {
			days = 28
		}
	}

	return days
}

//TimeConversion 时间换成成xx秒、xx分钟、xx小时、xx天前
func TimeConversion(atime int64) string {
	var byTime = []int64{365 * 24 * 60 * 60, 24 * 60 * 60, 60 * 60, 60, 1}
	var unit = []string{"年前", "天前", "小时前", "分钟前", "秒前"}
	now := time.Now().Unix()
	ct := now - atime
	if ct < 0 {
		return "刚刚"
	}
	var res string
	for i := 0; i < len(byTime); i++ {
		if ct < byTime[i] {
			continue
		}
		var temp = math.Floor(float64(ct / byTime[i]))
		ct = ct % byTime[i]
		if temp > 0 {
			var tempStr string
			tempStr = strconv.FormatFloat(temp, 'f', -1, 64)
			res = mergeString(tempStr, unit[i])
		}
		break //我想要的形式是精确到最大单位，即："2天前"这种形式，如果想要"2天12小时36分钟48秒前"这种形式，把此处break去掉，然后把字符串拼接调整下即可（别问我怎么调整，这如果都不会我也是无语）
	}
	return res
}

//GetWeek 判断时间是当年的第几周
func GetWeek(t time.Time) (year, week int) {
	year, week = t.ISOWeek()
	return
}

func mergeString(args ...string) string {
	buffer := bytes.Buffer{}
	for i := 0; i < len(args); i++ {
		buffer.WriteString(args[i])
	}
	return buffer.String()
}
