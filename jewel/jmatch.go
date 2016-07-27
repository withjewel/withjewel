package jewel

import (
	"bytes"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
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
	origURL     string
	cvtedURL    string // 暂时留作调试用，其实每次创建urlPattern都会调用MustCompiler强制转换为pat
	pat         *regexp.Regexp
	groupVarMap map[int]string
	handler     http.Handler
	childs      []*urlPatternNode
}

/*JewelMatchSystem Jewel-URL模式匹配系统。
 */
type JewelMatchSystem struct {
	regvarTypeMap  map[string]string
	myPat          *regexp.Regexp
	urlPatternTree *urlPatternNode
}

/*convert2Tree 获取路径段pathsegs的树的表示方式。
 */
func (this *JewelMatchSystem) convert2Tree(pathsegs []string, handler http.Handler) *urlPatternNode {
	var retNode *urlPatternNode
	i := 1
	// 逆序
	for ipathseg := len(pathsegs) - 1; ipathseg >= 0; ipathseg-- {
		pathseg := pathsegs[ipathseg]

		newURLPatternNode := &urlPatternNode{origURL: pathseg, groupVarMap: make(map[int]string), handler: nil}
		if retNode != nil {
			newURLPatternNode.childs = append(newURLPatternNode.childs, retNode)
		} else {
			newURLPatternNode.handler = handler
		}
		retNode = newURLPatternNode

		newURLPatternNode.cvtedURL = this.myPat.ReplaceAllStringFunc(pathseg, func(s string) string {
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

				newURLPatternNode.groupVarMap[i] = name
				i += countGroupBegin(regpat)

				return regpat
			}
		})
		newURLPatternNode.pat = regexp.MustCompile(newURLPatternNode.cvtedURL)
	}
	return retNode
}

/*AddPattern 向JewelMatchSystem 添加一条Jewel-URL模式。
 */
func (this *JewelMatchSystem) AddPattern(url string, h http.Handler) {
	pathsegs := strings.Split(url, "/")
	if pathsegs[0] == "" { // 处理'/'开头的情况
		pathsegs = pathsegs[1:]
	}
	curNode := this.urlPatternTree // assert(this.urlPatternTree != nil)
	for ipathseg, pathseg := range pathsegs {
		//fmt.Printf("[DEBUG] 匹配%s..\n", pathseg)
		found := false
		for _, childNode := range curNode.childs {
			submatchs := childNode.pat.FindStringSubmatch(pathseg)
			if submatchs != nil && len(submatchs[0]) == len(pathseg) {
				found = true
				curNode = childNode
				break
			}
		}
		if !found {
			newSubTree := this.convert2Tree(pathsegs[ipathseg:len(pathsegs)], h)
			curNode.childs = append(curNode.childs, newSubTree)
			return
		}
	}
	fmt.Printf("在添加URL模式%s时发现重复的URL模式。\n", url)
}

/*Match JewelMatchSystem尝试寻找匹配字符串url的Jewel-URL模式（以下简称模式）。
Jewel-URL是以/分隔的路径，开头和结尾的/会被忽略；连续出现的//会被合并为一个/。
路径以模式“根”开始，所有URL都至少能匹配根；根的处理器为nil并且无法注册处理器。
模式“/”代表匹配根下的一个空模式。模式“/static”代表匹配根下的名为static的模式。
Match会返回在匹配到的最后一个模式处注册的处理器。如：
	/
		返回为模式“/”（根模式下的空模式）注册的处理器，若模式“/”未注册处理器，返回根的处理器（nil）；
	/static/images/img.jpg
		Match会先在根模式下寻找匹配static的模式，若没有此模式，返回根的处理器（nil）；
		若存在static模式，则在static模式下寻找匹配images的模式，若没有此模式，返回static模式注册的处理器；
		以此类推……
*/
func (this *JewelMatchSystem) Match(url string) (http.Handler, string) {
	pathsegs := strings.Split(url, "/")
	if pathsegs[0] == "" { // 处理'/'开头的情况
		pathsegs = pathsegs[1:]
	}
	curNode := this.urlPatternTree // assert(this.urlPatternTree != nil)
	var matchedPat bytes.Buffer
	for _, pathseg := range pathsegs {
		fmt.Printf("[DEBUG] 尝试匹配路径段%s..\n", pathseg)
		found := false
		for _, childNode := range curNode.childs {
			submatchs := childNode.pat.FindStringSubmatch(pathseg)
			if submatchs != nil && len(submatchs[0]) == len(pathseg) {
				if JewelHanler, ok := childNode.handler.(*JewelHandler); ok {
					params := make(map[string]string)
					for k, v := range childNode.groupVarMap {
						params[v] = submatchs[k]
					}
					JewelHanler.InitParams(params)
				}
				matchedPat.WriteString("/")
				matchedPat.WriteString(childNode.origURL)
				found = true
				curNode = childNode
				break
			}
		}
		if !found {
			//fmt.Printf("URL<%s>在匹配路径<%s>时失败——没有在URL模式树中找到相应的路径。\n", url, pathseg)
			//return nil, ""
			break
		}
	}
	return curNode.handler, matchedPat.String()
}

/*NewJewelMatchSystem 创建JewelMatchSystem的一个实例。
 */
func NewJewelMatchSystem() *JewelMatchSystem {
	this := &JewelMatchSystem{regvarTypeMap: make(map[string]string)}
	this.regvarTypeMap[`id`] = `(\w+)`
	this.regvarTypeMap[`int`] = `(\d+)`
	this.myPat = regexp.MustCompile(`\\\\|\\<|<\w+:\w+>|.`)
	this.urlPatternTree = &urlPatternNode{}
	return this
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
