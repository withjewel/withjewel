package jedb

import (
	"container/list"
	"fmt"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

/*目前的实现比较简单。
 */
type dbPool struct {
	unused *list.List
	lock   sync.Mutex
}

var dbPoolGlobal dbPool

/*Init 初始化Jedb。
 */
func Init() {
	/*
		dbPoolGlobal.unused = list.New()
		for i := 0; i < 10; i++ {
			db, err := gorm.Open("mysql", "leslie:19941121sr@tcp(115.28.27.5:3306)/test")
			if err != nil {
				fmt.Printf("数据库连接出错：%s\n", err.Error())
			}
			db.SingularTable(true)
			dbPoolGlobal.unused.PushBack(db)
		}
	*/
}

/*Accquire 从Jedb获取一个已连接的Gorm.DB实例。
 */
func Accquire() *gorm.DB {
	db, err := gorm.Open("mysql", "leslie:19941121sr@tcp(115.28.27.5:3306)/test")
	if err != nil {
		fmt.Printf("数据库连接出错：%s\n", err.Error())
	}
	db.SingularTable(true)
	return db
	/*
		for {
			dbPoolGlobal.lock.Lock()
			if dbPoolGlobal.unused.Len() > 0 {
				elem := dbPoolGlobal.unused.Front()
				dbPoolGlobal.unused.Remove(elem)
				dbPoolGlobal.lock.Unlock()
				v := elem.Value
				db, ok := v.(*gorm.DB)
				if ok {
					return db
				}
			}
			dbPoolGlobal.lock.Unlock()
		}
	*/
}

/*Revert 将Gorm.DB实例归还给Jedb。
 */
func Revert(gormdb *gorm.DB) {
	gormdb.Close()
	//dbPoolGlobal.lock.Lock()
	//dbPoolGlobal.unused.PushBack(gormdb)
	//dbPoolGlobal.lock.Unlock()
}
