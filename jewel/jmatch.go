package jmatch

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type urlPattern struct {
	origURL     string
	cvtedURL    string
	groupVarMap map[int]string
	handler     interface{}
}

func (this *urlPattern) print() {
	fmt.Printf("原始URL：%s\n结果URL：%s\n分组索引：%+v\n", this.origURL, this.cvtedURL, this.groupVarMap)
}

/*JewelMatchSystem Jewel-URL模式匹配系统。
 */
type JewelMatchSystem struct {
	regvarTypeMap        map[string]string
	handlerURLPatternMap []*urlPattern
	myPat                *regexp.Regexp
}

func NewJewelMatchSystem() *JewelMatchSystem {
	this := &JewelMatchSystem{regvarTypeMap: make(map[string]string)}
	this.regvarTypeMap[`id`] = `(\w+)`
	this.regvarTypeMap[`int`] = `(\d+)`
	this.myPat = regexp.MustCompile(`\\\\|\\<|<\w+:\w+>|.`)
	return this
}

func (this *JewelMatchSystem) Print() {
	//fmt.Println(this.regvarTypeMap)
	for _, up := range this.handlerURLPatternMap {
		up.print()
	}
}

/*AddURL 向JewelMatch 添加一条Jewel-URL模式以待匹配。
 */
func (this *JewelMatchSystem) AddURL(url string, h interface{}) {
	i := 1
	newURLPattern := &urlPattern{origURL: url, groupVarMap: make(map[int]string), handler: h}
	newURLPattern.cvtedURL = this.myPat.ReplaceAllStringFunc(url, func(s string) string {
		if len([]rune(s)) == 1 {
			return regexp.QuoteMeta(s)
		}
		switch s {
		case `\\`:
			return regexp.QuoteMeta(`\`)
		case `\<`:
			return regexp.QuoteMeta(`<`)
		default:
			r := []rune(s)
			nameAndType := string(r[1 : len(r)-1])
			nameTypePair := strings.Split(nameAndType, ":")
			name := nameTypePair[0]
			ty := nameTypePair[1]

			regpat, exist := this.regvarTypeMap[ty]
			if !exist {
				fmt.Println("FATAL, the regvar's type doesn't exist")
				return s
			}

			newURLPattern.groupVarMap[i] = name
			i += countGroupBegin(regpat)

			//fmt.Printf("Name: %s, Type: %s RegPat: %s\n", name, ty, regpat)
			return regpat
		}
	})

	this.handlerURLPatternMap = append(this.handlerURLPatternMap, newURLPattern)
}

func (this *JewelMatchSystem) Match(url string) interface{} {
	for _, up := range this.handlerURLPatternMap {
		fmt.Printf("尝试匹配模式%s...\n", up.cvtedURL)
		pat := regexp.MustCompile(up.cvtedURL)
		submatchs := pat.FindStringSubmatch(url)
		if submatchs != nil {
			for k, v := range up.groupVarMap {
				fmt.Printf("%s的匹配结果是%s\n", v, submatchs[k])
				setField(up.handler, v, submatchs[k])
			}
			return up.handler
		}
	}
	return nil
}

func countGroupBegin(pat string) int {
	mypat := regexp.MustCompile(`\\*\(`)
	n := 0
	ss := mypat.FindAllString(pat, -1 /*no limit*/)
	for _, s := range ss {
		if len([]rune(s))%2 != 0 {
			n++
		}
	}
	return n
}

func setField(v interface{}, fieldName string, fieldv string) {
	structVal := reflect.ValueOf(v)
	if structVal.Kind() == reflect.Ptr {
		structVal = structVal.Elem()
	}
	fieldValue := structVal.FieldByName(fieldName)
	if !fieldValue.IsValid() {
		fmt.Printf("%+v没有字段%s\n", v, fieldName)
		return
	}
	if !fieldValue.CanSet() {
		fmt.Printf("%+v的字段%s不可设置\n", v, fieldName)
		return
	}

	switch fieldValue.Kind() {
	case reflect.String:
		fieldValue.SetString(fieldv)

	case reflect.Bool:
		i, err := strconv.ParseBool(fieldv)
		if err == nil {
			fieldValue.SetBool(i)
		}

	case reflect.Int8:
		i, err := strconv.ParseInt(fieldv, 10, 8)
		if err == nil {
			fieldValue.SetInt(i)
		}
	case reflect.Int16:
		i, err := strconv.ParseInt(fieldv, 10, 16)
		if err == nil {
			fieldValue.SetInt(i)
		}
	case reflect.Int32:
		i, err := strconv.ParseInt(fieldv, 10, 32)
		if err == nil {
			fieldValue.SetInt(i)
		}
	case reflect.Int64:
		i, err := strconv.ParseInt(fieldv, 10, 64)
		if err == nil {
			fieldValue.SetInt(i)
		}

	case reflect.Uint8:
		i, err := strconv.ParseUint(fieldv, 10, 8)
		if err == nil {
			fieldValue.SetUint(i)
		}
	case reflect.Uint16:
		i, err := strconv.ParseUint(fieldv, 10, 16)
		if err == nil {
			fieldValue.SetUint(i)
		}
	case reflect.Uint32:
		i, err := strconv.ParseUint(fieldv, 10, 32)
		if err == nil {
			fieldValue.SetUint(i)
		}
	case reflect.Uint64:
		i, err := strconv.ParseUint(fieldv, 10, 64)
		if err == nil {
			fieldValue.SetUint(i)
		}

	case reflect.Float32:
		i, err := strconv.ParseFloat(fieldv, 32)
		if err == nil {
			fieldValue.SetFloat(i)
		}
	case reflect.Float64:
		i, err := strconv.ParseFloat(fieldv, 64)
		if err == nil {
			fieldValue.SetFloat(i)
		}
	}
}

type MyStruct struct {
	Page  int
	Child string
}

func Demo() {
	v := &MyStruct{}
	pat := `/\123\<user:id>>ffda/\\<Page:int>/123</<1/<Child:id>`
	jm := NewJewelMatchSystem()
	jm.AddURL(pat, v)
	str := `/\123<user:id>>ffda/\999/123</<1/USER`
	fmt.Printf("用户希望截获的URL：%s\n", pat)
	fmt.Printf("请求URL：%s\n", str)
	jm.Match(str)
	fmt.Printf("ok\n")
}
