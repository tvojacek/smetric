package smetric

import (
	"testing"
)

func TestAddChanLength(t *testing.T) {
	type args[T any] struct {
		channelMetrics ChannelMetrics
		chanel         chan T
	}
	type testCase[T any] struct {
		name string
		args args[T]
	}
	tests := []testCase[int]{
		{
			name: "",
			args: args[int]{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.args.channelMetrics.AddLength(GetChanLength(tt.args.chanel))
		})
	}
}

func BenchmarkAddChanLength(b *testing.B) {
	input := make(chan int, 10)

	var metrics struct {
		inputChan ChannelMetrics
	}
	err := InitMetricStruct(&metrics, NewNameBuilder("foo"), nil)
	if err != nil {
		b.Fatal(err)
	}

	b.Run("AddLength", func(b *testing.B) {
		b.ReportAllocs()
		for _ = range b.N {
			metrics.inputChan.AddLength(GetChanLength(input))
			metrics.inputChan.Destroy("length")
		}
	})
}
