package smetric

import (
	"github.com/VictoriaMetrics/metrics"
)

type ChannelMetrics struct {
	Metrics
	Total    *metrics.Counter
	OverFlow *metrics.Counter
}

func (channelMetrics *ChannelMetrics) Reset() {
	channelMetrics.Total.Set(0)
	channelMetrics.OverFlow.Set(0)
}

func GetChanLength[T any](chanel chan T) func() float64 {
	return func() float64 { return float64(len(chanel)) }
}

func (channelMetrics *ChannelMetrics) AddLength(getLen func() float64) {
	channelMetrics.AddGauge("length", getLen)
}
