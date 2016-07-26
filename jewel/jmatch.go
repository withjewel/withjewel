package jewel

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type urlPattern struct {
	origURL     string
	cvtedURL    string // 暂时留作调试用，其实每次创建urlPattern都会调用MustCompiler强制转换为pat
	pat         *regexp.Regexp
	groupVarMap map[int]string
	handler     http.Handler
}

type urlPatternNode struct {
	pat    *urlPattern
	childs []*urlPatternNode
}

type urlPatternSlice []*urlPattern

func (s urlPatternSlice) Len() int {
	return len(s)
}
func (s urlPatternSlice) Less(left, right int) bool {
	return s[left].origURL > s[right].origURL
}
func (s urlPatternSlice) Swap(left, right int) {
	s[left], s[right] = s[right], s[left]
}

func (this *urlPattern) print() {
	fmt.Printf("原始URL：%s\n结果URL：%s\n分组索引：%+v\n", this.origURL, this.cvtedURL, this.groupVarMap)
}

/*JewelMatchSystem Jewel-URL模式匹配系统。
 */
type JewelMatchSystem struct {
	regvarTypeMap        map[string]string
	handlerURLPatternMap urlPatternSlice
	myPat                *regexp.Regexp
	urlPatternTree       *urlPatternNode
}

func (this *JewelMatchSystem) convert2Tree(pathsegs []string) *urlPatternNode {
	var retNode *urlPatternNode
	i := 1
	for ipathseg := len(pathsegs) - 1; ipathseg >= 0; ipathseg++ {
		pathseg := pathsegs[ipathseg]

		newURLPattern := &urlPattern{origURL: pathseg, groupVarMap: make(map[int]string), handler: nil}
		newURLPatternNode := &urlPatternNode{pat: newURLPattern}
		if retNode != nil {
			newURLPatternNode.childs = append(newURLPatternNode.childs, retNode)
		} else {

		}
		retNode = newURLPatternNode
		newURLPattern.cvtedURL = this.myPat.ReplaceAllStringFunc(pathseg, func(s string) string {
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
		newURLPattern.pat = regexp.MustCompile(newURLPattern.cvtedURL)
	}
	return retNode
}

func (this *JewelMatchSystem) AddPattern2(url string, h http.Handler) {
	pathsegs := strings.Split(url, "/")
	curNode := this.urlPatternTree // assert(this.urlPatternTree != nil)
	for ipathseg, pathseg := range pathsegs {
		fmt.Printf("[DEBUG] 匹配%s..\n", pathseg)
		found := false
		for _, childNode := range curNode.childs {
			submatchs := childNode.pat.pat.FindStringSubmatch(pathseg)
			if submatchs != nil && len(submatchs[0]) == len(pathseg) {
				found = true
				curNode = childNode
				break
			}
		}
		if !found {
			curNode.childs = append(curNode.childs, this.convert2Tree(pathsegs[ipathseg:len(pathsegs)]))
			return
		}
	}
}

func (this *JewelMatchSystem) Match2(url string) {
	pathsegs := strings.Split(url, "/")
	curNode := this.urlPatternTree // assert(this.urlPatternTree != nil)
	for _, pathseg := range pathsegs {
		fmt.Printf("[DEBUG] 匹配%s..\n", pathseg)
		found := false
		for _, childNode := range curNode.childs {
			submatchs := childNode.pat.pat.FindStringSubmatch(pathseg)
			if submatchs != nil && len(submatchs[0]) == len(pathseg) {
				found = true
				curNode = childNode
				break
			}
		}
		if !found {
			fmt.Printf("URL<%s>在匹配路径<%s>时失败——没有在URL模式树中找到相应的路径。\n", url, pathseg)
			return
		}
	}
}

/*NewJewelMatchSystem 创建JewelMatchSystem的一个实例。
 */
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

/*AddURL 向JewelMatchSystem 添加一条Jewel-URL模式。
 */
func (this *JewelMatchSystem) AddPattern(url string, h http.Handler) {
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
	newURLPattern.pat = regexp.MustCompile(newURLPattern.cvtedURL)

	this.handlerURLPatternMap = append(this.handlerURLPatternMap, newURLPattern)
	sort.Sort(this.handlerURLPatternMap)
}

/*Match JewelMatchSystem尝试寻找匹配字符串url的Jewel-URL模式，
若成功匹配，返回该模式注册的处理器；若无模式匹配，返回nil。
*/
func (this *JewelMatchSystem) Match(url string) (http.Handler, string) {
	for _, up := range this.handlerURLPatternMap {
		//fmt.Printf("尝试匹配模式%s...\n", up.cvtedURL)
		submatchs := up.pat.FindStringSubmatch(url)
		if submatchs != nil && len(submatchs[0]) == len(url) {
			if JewelHanler, ok := up.handler.(*JewelHandler); ok {
				params := make(map[string]string)
				for k, v := range up.groupVarMap {
					params[v] = submatchs[k]
				}
				JewelHanler.InitParams(params)
			}
			return up.handler, up.origURL
		}
	}
	return nil, ""
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

	case reflect.Int:
		i, err := strconv.ParseInt(fieldv, 10, strconv.IntSize)
		if err == nil {
			fieldValue.SetInt(i)
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

	case reflect.Uint:
		i, err := strconv.ParseUint(fieldv, 10, strconv.IntSize)
		if err == nil {
			fieldValue.SetUint(i)
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
