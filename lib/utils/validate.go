package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	zhongwen "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	validator "gopkg.in/go-playground/validator.v9"
	zh_translations "gopkg.in/go-playground/validator.v9/translations/zh"
)

const (
	lower = false
	upper = true
)

var (
	trans ut.Translator
	//Validate ...
	Validate *validator.Validate
	//ValidateServerError ...
	ValidateServerError = "服务内部错误"
	//ValidateNotFound ...
	ValidateNotFound = "接口不存在"
	//ValidateUnauthorized ...
	ValidateUnauthorized = "未授权"
	//ValidateError ...
	ValidateError = "请求参数不正确"
	//ValidateRequestJSONError ...
	ValidateRequestJSONError = "请求体格式不正确"

	//驼峰格式转换
	commonInitialisms              = []string{"ACL", "API", "ASCII", "CPU", "CSS", "DNS", "EOF", "GUID", "HTML", "HTTP", "HTTPS", "ID", "IP", "JSON", "LHS", "QPS", "RAM", "RHS", "RPC", "SLA", "SMTP", "SQL", "SSH", "TCP", "TLS", "TTL", "UDP", "UI", "UID", "UUID", "URI", "URL", "UTF8", "VM", "XML", "XMPP", "XSRF", "XSS"}
	commonInitialismsReplacer      *strings.Replacer
	uncommonInitialismsReplacer    *strings.Replacer
	commonInitialismsForReplacer   []string
	uncommonInitialismsForReplacer []string
)

//CustomError 错误格式
type CustomError map[string]interface{}

func init() {
	//实例化需要转换的语言
	zh := zhongwen.New()
	uni := ut.New(zh, zh)
	trans, _ = uni.GetTranslator("zh")
	Validate = validator.New()

	//注册转换的语言为默认语言
	zh_translations.RegisterDefaultTranslations(Validate, trans)

	//驼峰转换
	for _, initialism := range commonInitialisms {
		commonInitialismsForReplacer = append(commonInitialismsForReplacer, initialism, strings.Title(strings.ToLower(initialism)))
		uncommonInitialismsForReplacer = append(uncommonInitialismsForReplacer, strings.Title(strings.ToLower(initialism)), initialism)
	}
	commonInitialismsReplacer = strings.NewReplacer(commonInitialismsForReplacer...)
	uncommonInitialismsReplacer = strings.NewReplacer(uncommonInitialismsForReplacer...)
}

//NewValidatorError ...
func NewValidatorError(err error, m map[string]string) CustomError {
	res := CustomError{}
	errs := err.(validator.ValidationErrors)
	for _, e := range errs {
		transtr := e.Translate(trans)
		//将结构体字段转换map中的key为小写
		f := underscoreName(e.Field())

		//Key: 'Deployment.Team' Error:Field validation for 'Team' failed on the '(team & env) 或者 ns 只有有一个必须指定' tag
		var valid = regexp.MustCompile(".*failed on the '(.*)' tag")
		transtrs := valid.FindStringSubmatch(transtr)
		if len(transtrs) > 1 {
			transtr = transtrs[1]
		}

		//判断错误字段是否在命名中，如果在，则替换错误信息中的字段
		if rp, ok := m[e.Field()]; ok {
			res[f] = strings.Replace(transtr, e.Field(), rp, 1)
		} else {
			res[f] = transtr
		}
	}
	//返回错误信息
	return res
}

//NewValidatorJSONError ...
func NewValidatorJSONError(err error, m map[string]string) CustomError {
	res := CustomError{}
	switch err.(type) {
	case *json.SyntaxError:
		res["expection"] = "非法的JSON格式"
	case *json.UnmarshalTypeError:
		errs := err.(*json.UnmarshalTypeError)
		if errs.Field != "" {
			value := strings.Split(errs.Field, ".")
			key := marshal(value[0])
			if _, ok := m[key]; ok {
				res[errs.Field] = fmt.Sprintf("%s非法参数格式[%s]", m[key], errs.Value)
			} else {
				res[errs.Field] = fmt.Sprintf("%s非法参数格式[%s]", key, errs.Value)
			}
		}
	default:
		res["expection"] = err
	}
	return res
}

func underscoreNameEasy(name string) string {
	buffer := NewBuffer()
	for i, r := range name {
		if unicode.IsUpper(r) {
			if i != 0 {
				buffer.Append('_')
			}
			buffer.Append(unicode.ToLower(r))
		} else {
			buffer.Append(r)
		}
	}

	fmt.Println(name, "--------->", buffer.String())
	return buffer.String()
}

func underscoreName(name string) string {

	if name == "" {
		return ""
	}

	var (
		value                                    = commonInitialismsReplacer.Replace(name)
		buf                                      = bytes.NewBufferString("")
		lastCase, currCase, nextCase, nextNumber bool
	)

	for i, v := range value[:len(value)-1] {
		nextCase = bool(value[i+1] >= 'A' && value[i+1] <= 'Z')
		nextNumber = bool(value[i+1] >= '0' && value[i+1] <= '9')

		if i > 0 {
			if currCase == upper {
				if lastCase == upper && (nextCase == upper || nextNumber == upper) {
					buf.WriteRune(v)
				} else {
					if value[i-1] != '_' && value[i+1] != '_' {
						buf.WriteRune('_')
					}
					buf.WriteRune(v)
				}
			} else {
				buf.WriteRune(v)
				if i == len(value)-2 && (nextCase == upper && nextNumber == lower) {
					buf.WriteRune('_')
				}
			}
		} else {
			currCase = upper
			buf.WriteRune(v)
		}
		lastCase = currCase
		currCase = nextCase
	}

	buf.WriteByte(value[len(value)-1])

	s := strings.ToLower(buf.String())
	return s
}

func marshal(name string) string {
	if name == "" {
		return ""
	}

	temp := strings.Split(name, "_")
	var s string
	for _, v := range temp {
		vv := []rune(v)
		if len(vv) > 0 {
			if bool(vv[0] >= 'a' && vv[0] <= 'z') { //首字母大写
				vv[0] -= 32
			}
			s += string(vv)
		}
	}

	s = uncommonInitialismsReplacer.Replace(s)

	return s
}
