# smetric structured metrics
go helper to simply initializing structs holding prometheus metrics based on "github.com/VictoriaMetrics/metrics"
InitMetricStruct initialize any variable holding pointer to metrics, register metrics with names generated from names of the fields of struct

## Usage
Basic usage
```go
    package main

    import (
    	"github.com/VictoriaMetrics/metrics"
		"github.com/tvojacek/smetric"
    )
	
    func main() {
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
		
	    myMetrics.MessagesConsumed.Success.Inc()
}
```
