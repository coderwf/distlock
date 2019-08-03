package lock

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)


//获取设置的锁超时自动失效时间
func getExpireTime(expireKey string) (d time.Duration, err error){

	//DIST_REDIS_EXPIRE
	//先从环境变量中获取超时时间,如果没有则从默认key中获取
	if expireKey == ""{
		expireKey = "DIST_REDIS_EXPIRE"
	}

	expireTimeStr := os.Getenv(expireKey)

	if expireTimeStr == ""{
		//默认使用30s

		d = time.Second * 30
		return
	}

	//解析时间
	expireTime, err := strconv.ParseInt(expireTimeStr, 10, 64)

	if err != nil{
		err = fmt.Errorf("%s from env is invalid number", expireKey)
		return
	}

	//ms
	d = time.Duration(expireTime) * time.Millisecond
	return
}


func randId(idLen int) string{
	buffer := make([]byte, idLen)
	var gen int64

    for i:= 0; i< idLen; i++{
    	if i % 8 == 0{
    		gen = rand.Int63()
		}

    	buffer[i] = byte(gen & 0x0FF)
    	gen >>= 8
	}//for

	return fmt.Sprintf("%x", buffer)
}


func getIp() string{
	ip := os.Getenv("LOCAL_IP")
	if ip == ""{
		ip = "none"
	}
	return ip
}


func lockId() string{
	return fmt.Sprintf("%s-%d-%s", randId(5), time.Now().Unix(), getIp())
}
