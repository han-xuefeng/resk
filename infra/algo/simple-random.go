package algo

import (
	"math/rand"
	"time"
)

const min = int64(1)

//简单随机算法
//红包的数量，红包金额
//金额单位为分，1元钱=100分
func SimpleRand(count, amount int64) int64 {
	//当红包数量剩余一个的时候，就直接返回剩余金额
	if count == 1 {
		return amount
	}
	//计算最大可调度金额
	max := amount - min*count
	rand.Seed(time.Now().UnixNano())
	x := rand.Int63n(max) + min
	return x
}
