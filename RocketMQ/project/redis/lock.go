package redis

import (
	"github.com/garyburd/redigo/redis"
	uuid "github.com/satori/go.uuid"
	"log"
	"sync"
)

var mutex sync.Mutex

func Lock(lockKey string) {
	mutex.Lock()
	defer mutex.Unlock()

	clientID := uuid.NewV4().String()
	ok := Redis.SetWitLock(lockKey, clientID, 10)
	if !ok {
		return
	}
	defer func() {
		// 删除自己所占用的锁， 看值是否一致，一致则删除，lua脚本实现
		script := "if redis.call('get',KEYS[1]) == ARGV[1] then return redis.call('del',KEYS[1]) else return 0  end"
		s := redis.NewScript(1, script)
		_, err := s.Do(Redis.RS, lockKey, clientID)
		if err != nil {
			log.Fatal(err)
		}
	}()

	//stack, _ := strconv.Atoi(Redis.Get("num"))
	//if stack > 0 {
	//	newStack := stack - 1
	//	res := Redis.Set("num", newStack)
	//	if res {
	//		fmt.Println("库存修改完毕, 剩余库存：" + strconv.Itoa(newStack) + "-8050")
	//		return
	//	}
	//	fmt.Println("库存修改失败-8050")
	//	return
	//}
	// 没有库存
	//return
}
