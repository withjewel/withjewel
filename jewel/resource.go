package jewel

import (
    "io/ioutil"
    "fmt"
    "os"
    "path/filepath"
    "strings"
)

/*Get 获取name指定的资源
*/
func Get(name string)([]byte, error) {
    name = strings.Replace(name, "/", "\\", -1)
    programAbsDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
    bytes, err := ioutil.ReadFile(programAbsDir + name)
    if err != nil {
        fmt.Printf("在读取文件%s的时候失败。\n", programAbsDir + name)
        return nil, err
    }
    return bytes, err
}

/**
 * 判断文件是否存在  存在返回 true 不存在返回false
 */
func checkFileIsExist(filename string) (bool) {
 var exist = true;
 if _, err := os.Stat(filename); os.IsNotExist(err) {
  exist = false;
 }
 return exist;
}