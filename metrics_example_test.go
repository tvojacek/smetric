package smetric

import (
	"fmt"
	"testing"

	"github.com/VictoriaMetrics/metrics"
)

func ExampleInitMetricStruct() {
	type OperationMetrics struct {
		Success  *metrics.Counter
		Failures *metrics.Counter
	}
	type CommonMetrics struct {
		Healthy       *metrics.Counter
		Reconnections *metrics.Counter
	}
	type exampleMetrics struct {
		CommonMetrics
		MessagesConsumed OperationMetrics
		WithCustomName   OperationMetrics `metric:"writes"`
	}

	var myMetrics exampleMetrics

	metricNamesBuilder := NewNameBuilder("database").WithParameter("url", "localhost")
	set := metrics.NewSet()
	err := InitMetricStruct(&myMetrics, metricNamesBuilder, set)
	if err != nil {
		panic(err)
	}

	for _, name := range set.ListMetricNames() {
		fmt.Println(name)
	}
	set.UnregisterAllMetrics()
	// Unordered output: database_healthy{url="localhost"}
	// database_reconnections{url="localhost"}
	// database_messages_consumed_failures{url="localhost"}
	// database_messages_consumed_success{url="localhost"}
	// database_writes_failures{url="localhost"}
	// database_writes_success{url="localhost"}
}

// ExampleInitMetricStruct_gauges cleanup can be omited
func ExampleInitMetricStruct_gauges() {

	type exampleMetrics struct {
		Input ChannelMetrics
	}

	var myMetrics exampleMetrics

	metricNamesBuilder := NewNameBuilder("data_processor").WithParameter("url", "localhost")
	set := metrics.NewSet()
	err := InitMetricStruct(&myMetrics, metricNamesBuilder, set)
	if err != nil {
		panic(err)
	}

	var myChan chan struct{}

	myMetrics.Input.AddLength(GetChanLength(myChan))

	if err != nil {
		panic(err)
	}
	fmt.Println(myMetrics.Input.Name())
	for _, name := range set.ListMetricNames() {
		fmt.Println(name)
	}
	// Unordered output: input
	//data_processor_input_length{url="localhost"}
	//data_processor_input_over_flow{url="localhost"}
	//data_processor_input_total{url="localhost"}
}

type deleter struct {
}

func (d deleter) Del() {}

func TestFoo(t *testing.T) {

}
