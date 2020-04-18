package dao

import "github.com/garyburd/redigo/redis"

var rdb redis.Conn

func InitRedis() (err error) {
	rdb,err = redis.Dial("tcp","localhost:6379")
	return
}

func Rclose()  {
	rdb.Close()
}

func AddCaptcha(num int) error {
	_,err := rdb.Do("SET","captcha",num,"EX","300")
	return err
}

func GetCaptcha() (string, error) {
	num,err := redis.String(rdb.Do("GET","captcha"))
	return num,err
}