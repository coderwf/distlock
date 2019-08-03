package lock

import (
	"fmt"
	"github.com/go-redis/redis"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var RedisClient *redis.Client

//初始化
func init(){

	//rand seed
	rand.Seed(time.Now().Unix())


	// DIST_REDIS_ADDR
	// DIST_REDIS_PASS
	// DIST_REDIS_CLIENT


	//从环境变量中获取redis的地址
	//如果获取失败则默认使用本地地址
	distRedisAddr := os.Getenv("DIST_REDIS_ADDR")
	if distRedisAddr == ""{
		distRedisAddr = "127.0.0.1:6379"
	}

	//从环境变量中获取redis密码
	distRedisPassword := os.Getenv("DIST_REDIS_PASS")


	//从环境变量中获取redis使用的db,没有则默认为0
	var distRedisClient int64
	var err error
	distRedisClientStr := os.Getenv("DIST_REDIS_CLIENT")
	if distRedisClientStr != ""{
		distRedisClient, err = strconv.ParseInt(distRedisClientStr, 10, 32)
		if err != nil{
			panic(fmt.Errorf("'DIST_REDIS_CLIENT' from env is invalid number"))
		}//if
	}//if

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     distRedisAddr,
		Password: distRedisPassword,
		DB:       int(distRedisClient),
	})
}





//如果不存在则设置key-value,并设置过期时间,返回1
//如果存在则直接返回0
func SetIfNotExistAndExpire(key string, value string, expire time.Duration) (int, error){

	//如果key存在则直接返回1
	//如果key不存在则设置key并设置超时时间(ms为单位)然后返回
	script := "if (redis.call('exists', KEYS[1]) == 1) then return 0 end;" +
		      "if redis.call('psetex', KEYS[1], ARGV[2], ARGV[1]) then return 1 else return 0 end;"

	//以ms为单位则需要将time.Duration / 1000,000
	res, err := RedisClient.Eval(script, []string{key}, value, int64(expire)/ 1000000).Result()
	return int(res.(int64)), err
}



//如果key-value不存在则设置并且设置过期时间,返回1
//如果key-value存在,旧的value和此value相等,则更新过期时间并返回1
//如果key-value存在,且旧的value != value则返回0

func SetINEOrUpdate(key string, value string, expire time.Duration) (int, error){

	script := "local value = redis.call('get', KEYS[1]);" +
		"if value ~= false and value ~= ARGV[1] then return 0 end;" +
		"if redis.call('psetex', KEYS[1], ARGV[2], ARGV[1]) then return 1 else return 0 end;"

	//以ms为单位则需要将time.Duration / 1000,000
	res, err := RedisClient.Eval(script, []string{key}, value, int64(expire)/ 1000000).Result()
	return int(res.(int64)), err
}


//如果key在redis中对应的value是此value则删除这个key
//否则不能删除这个key
func DelIfNotSelfValue(key string, value string) (int, error){

	script := "local value = redis.call('get', KEYS[1]);" +
		"if value == false or value ~=ARGV[1] then return 0 end;" +
		"return redis.call('del', KEYS[1]);"

	res, err := RedisClient.Eval(script, []string{key}, value).Result()

	return int(res.(int64)), err
}