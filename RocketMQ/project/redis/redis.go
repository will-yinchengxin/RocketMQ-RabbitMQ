package redis

import (
	"github.com/garyburd/redigo/redis"
	"log"
	_const "rock/const"
)

var Redis RedisResoce

type RedisResoce struct {
	RS redis.Conn
}

func InitRedis() *RedisResoce {
	host := "172.16.252.99:6379"

	rs, err := redis.Dial("tcp", host)
	//defer rs.Close()

	if err != nil {
		log.Println(err)
		return &RedisResoce{}
	}
	Redis = RedisResoce{
		rs,
	}
	return &Redis
}

func (rs *RedisResoce) Get(key string) string {
	intVal, err := redis.String(rs.RS.Do("get", key))
	if err != nil {
		log.Fatal("get err, ", err)
		return ""
	}
	return intVal
}

func (rs *RedisResoce) Set(key string, val int) bool {
	_, err := redis.String(rs.RS.Do("set", key, val))
	if err != nil {
		log.Fatal("set err, ", err)
		return false
	}
	return true
}

func (rs *RedisResoce) Expire(key string) {
	_, err := rs.RS.Do("expire", key, 10)
	if err != nil {
		log.Fatal("expire err, ", err)
		return
	}
}

func (rs *RedisResoce) SetWitLock(key string, val string, time int) bool { //SET test 1 EX 10 NX
	intVal, err := rs.RS.Do("set", key, val, "EX", time, "NX")
	if err != nil {
		log.Fatal(err)
		return false
	}
	if intVal != nil {
		return true
	}
	return false
}

func (rs *RedisResoce) SETNX(key string, val int) bool {
	intVal, err := redis.Int(rs.RS.Do("setnx", key, val))
	if err != nil {
		log.Fatal("setnx err, ", err)
		return false
	}
	return intVal > 0
}

func (rs *RedisResoce) Del(key string) {
	_, err := rs.RS.Do("del", key)
	if err != nil {
		log.Fatal("expire err, ", err)
		return
	}
}

func (rs *RedisResoce) HSETWithKey(hkey string, key string, val string) {
	_, err := rs.RS.Do("hset", hkey, key, val)
	if err != nil {
		log.Fatal("expire err, ", err)
		return
	}
}

func (rs *RedisResoce) HSET(key string, val string) {
	_, err := rs.RS.Do("hset", _const.RedisKey, key, val)
	if err != nil {
		log.Fatal("expire err, ", err)
		return
	}
}
