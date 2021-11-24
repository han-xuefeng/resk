package algo

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestSimpleRand(t *testing.T) {

	ForTest("简单随机算法", t, SimpleRand)
}
func TestBeforeShuffle(t *testing.T) {
	ForTest("先洗牌算法", t, BeforeShuffle)
}
//func TestBeforeShuffle(t *testing.T) {
//	ForTest("后洗牌算法", t, AfterShuffle)
//}

func TestDoubleRandom(t *testing.T) {
	ForTest("二次随机算法", t, DoubleRandom)
}
func TestDoubleAverage(t *testing.T) {
	ForTest("二倍均值算法", t, DoubleAverage)
}
func ForTest(message string, t *testing.T, fn func(count, amount int64) int64) {
	count, amount := int64(10), int64(10000)
	remain := amount
	sum := int64(0)
	for i := int64(0); i < count; i++ {
		x := fn(count-i, remain)
		fmt.Println(x)
		remain -= x
		sum += x
	}

	Convey(message, t, func() {
		Convey(message, func() {
			So(sum, ShouldEqual, amount)
		})
	})
}
