package quality

import (
	"math/rand"
	"time"
)

type DelayQualityAgent struct {
	delaymin int
	delaymax int
}

func NewDelayQualityAgent(min int, max int) (a *DelayQualityAgent) {
	return &DelayQualityAgent{min, max}
}

func (a *DelayQualityAgent) Pass(input interface{}) interface{} {
	n := rand.Intn(a.delaymax-a.delaymin) + a.delaymin
	time.Sleep(time.Duration(n) * time.Second)
	return input
}
