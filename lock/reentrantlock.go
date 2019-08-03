package lock

import (
    "fmt"
    "sync"
    "time"
)

//可重入分布式锁
//可重入是针对同一个ReentrantLock调用两次lock可以重入
//此重入和携程无关,不同的携程用同一个ReentrantLock去lock仍然可以重入

//想要实现并发锁住的地方必须要使用一个新的ReentrantLock即可

type ReentrantLock struct {
    //锁,更改counter的时候使用
    mu sync.Mutex

	//从环境变量中获取超时时间,以ms为单位
    expire time.Duration

    //锁的唯一标记
    id string

    //重入次数,为0则彻底释放锁
    counter int64
}


func (rl *ReentrantLock) lock(key string) bool{
    rl.mu.Lock()
    defer rl.mu.Unlock()

    success, _ := SetINEOrUpdate(key, rl.id, rl.expire)
    if success == 1{
        //counter + 1
        rl.counter += 1
        return true
    }
    return false
}


func (rl *ReentrantLock) Lock(key string){
    //
    for {
        if rl.lock(key){
            //加锁成功则直接返回
            return
        }//if

        //sleep一段时间(100ms)
        time.Sleep(100 * time.Millisecond)

    }//for

}

func (rl *ReentrantLock) TryLock(key string) bool{
    return rl.lock(key)
}

func (rl *ReentrantLock) TryLockTimeout(key string, timeout time.Duration) bool{
    //尝试加锁,超时则直接返回true/false

    //每10ms获取一次
    start := time.Now()

    for {
        if rl.lock(key){
            return true
        }

        if time.Now().Sub(start) > timeout{
            return false
        }

        time.Sleep(10 * time.Millisecond)
    }//for
}

func (rl *ReentrantLock) UnLock(key string) bool{
    rl.mu.Lock()
    defer rl.mu.Unlock()

    //没有加过锁则直接返回
    if rl.counter == 0{
        return true
    }

    rl.counter -= 1

    //如果counter为0了则可以彻底释放锁
    if rl.counter == 0{
        success, _ := DelIfNotSelfValue(key, rl.id)
        if success == 1{
            return true
        }else{
            return false
        }
    }

    return true
}


func (rl *ReentrantLock) SetExpireKey(expireKey string) (err error){

    rl.mu.Lock()
    defer rl.mu.Unlock()

    if rl.counter > 0{
        err = fmt.Errorf("Can't set expire when holding mutex ")
        return
    }
	//设置超时时间
    rl.expire, err = getExpireTime(expireKey)
    if err != nil{
        rl.expire = 30 * time.Second
        return
    }
    return
}


func NewReentrantLock(expireKey string) (dl *ReentrantLock, err error){
    var d time.Duration

    d, err = getExpireTime(expireKey)
    if err != nil{
        return
    }//if

    dl = &ReentrantLock{
        expire: d,
        id: lockId(),
    }
    return
}