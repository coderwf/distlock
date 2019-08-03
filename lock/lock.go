package lock

import "time"

//lock interface


type Lock interface {

	//Lock,尝试获取锁并阻塞一直等待获取
	Lock(key string)

	//尝试获取锁,如果未获取到则立刻返回
	TryLock(key string) bool

	//尝试获取锁,如果未获取到则阻塞一段时间继续获取,
	TryLockTimeout(key string, timeout time.Duration) bool

	//释放锁
	UnLock(key string)

	//设置超时时间key,从环境变量中获取超时时间
	SetExpireKey(expireKey string)
}
