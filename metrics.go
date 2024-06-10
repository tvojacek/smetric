package smetric

import (
	"fmt"
	"github.com/VictoriaMetrics/metrics"
	"reflect"
	"sync"
)

// InitMetricStruct register metrics defined in metricStruct.
// Names of metric are generated from tags `metric:"custom_suffix"` or names of struct fields converted to snake_case.
// Each nested struct add its name as prefix for nested metrics
// Implicit names for embed struct are ignored. Tags are never ignored.
//
// Supported metrics types:
//   - *metric.Counters,
//   - *metric.FloatCounters
//   - *metric.Summaries
//   - *metric.Histograms
//   - Gauges
//
// Function parameter metricStruct must be pointer to struct.
//
// set can be nil, default set used.
//
// Gauges helper allow to create gauges with same structured name like rest of the metrics.
//
// Func can panic if metric with duplicate name is already registered.
// It is user responsibility that metrics with same names do not exist in set prior calling function
// and tags does not have duplicate values.
func InitMetricStruct(metricStruct interface{}, nameBuilder NameBuilder, set *metrics.Set) error {
	if err := nameBuilder.Error(); err != nil {
		return err
	}
	if set == nil {
		set = metrics.GetDefaultSet()
	}
	val := reflect.ValueOf(metricStruct)

	if val.Kind() != reflect.Pointer {
		return fmt.Errorf("InitMetricStruct non-pointer (" + val.String() + ")")
	}

	if !val.IsValid() {
		return fmt.Errorf("init MetricStruct nil")
	}
	el := val.Elem()
	m := Metrics{
		set:           set,
		metricsStruct: metricStruct,
		nameBuilder:   nameBuilder,
	}
	err := m.initMetrics(el, nameBuilder)

	if err != nil {
		_ = m.unregisterMetricMetrics(el, nameBuilder)
		return err
	}

	return nil
}

type structUpdateFunc func(value reflect.Value, nameBuilder NameBuilder) error

func add[T any](value reflect.Value, nameBuilder NameBuilder, newMetric func(name string) *T) error {
	name, err := nameBuilder.String()
	if err != nil {
		return err
	}
	counter := newMetric(name)
	value.Set(reflect.ValueOf(counter))
	return nil
}

func (m *Metrics) initMetrics(value reflect.Value, nameBuilder NameBuilder) error {
	if !value.IsValid() {
		return fmt.Errorf("init MetricStruct nil or invalid :" + value.Type().String())
	}
	if !value.CanInterface() {
		return nil
	}
	switch v := value.Interface().(type) {
	case Metrics:
		gauges := Metrics{
			set:         m.set,
			gauges:      make(map[string]*metrics.Gauge),
			nameBuilder: nameBuilder,
			lock:        new(sync.Mutex),
		}
		value.Set(reflect.ValueOf(gauges))
	case metrics.Counter, metrics.FloatCounter, metrics.Summary, metrics.Histogram:
		return fmt.Errorf("non pointer type " + reflect.TypeOf(v).String())
	case *metrics.Counter:
		return add(value, nameBuilder, m.set.NewCounter)
	case *metrics.FloatCounter:
		return add(value, nameBuilder, m.set.NewFloatCounter)
	case *metrics.Summary:
		return add(value, nameBuilder, m.set.NewSummary)
	case *metrics.Histogram:
		return add(value, nameBuilder, m.set.NewHistogram)
	default:
		err := crawlStruct(value, nameBuilder, m.initMetrics)
		if err != nil {
			return err
		}
	}
	return nil
}

func crawlStruct(value reflect.Value, nameBuilder NameBuilder, structUpdate structUpdateFunc) error {
	if value.Kind() == reflect.Struct {
		for i := 0; i < value.NumField(); i++ {
			field := value.Field(i)
			t := value.Type().Field(i)
			metricName, ok := t.Tag.Lookup("metric")

			if !ok {
				if !t.Anonymous {
					var err error
					metricName, err = GetSnakeMetricName(t.Name)
					if err != nil {
						return err
					}
				} else {
					metricName = ""
				}
			}

			mb := nameBuilder
			if metricName != "" {
				mb = mb.WithSuffix(metricName)
			}
			err := structUpdate(field, mb)
			if err != nil {
				return fmt.Errorf("%w, metric: %s", err, mb.fullName)
			}
		}
	}
	return nil
}

func (m *Metrics) unregisterMetricMetrics(value reflect.Value, nameBuilder NameBuilder) error {
	if m == nil {
		return fmt.Errorf("failed: metrics can not be nil")
	}
	if !value.IsValid() {
		return fmt.Errorf("init MetricStruct nil or invalid :" + value.Type().String())
	}
	if !value.CanInterface() {
		return nil
	}

	switch v := value.Interface().(type) {
	case Metrics:
		err := v.DestroyAll()
		if err != nil {
			name, _ := nameBuilder.String()
			return fmt.Errorf("%w, metrics: %s", err, name)
		}
	case metrics.Counter, metrics.FloatCounter, metrics.Summary, metrics.Histogram:
		return fmt.Errorf("non pointer type " + reflect.TypeOf(v).String())
	case *metrics.Counter, *metrics.FloatCounter, *metrics.Summary, *metrics.Histogram:
		name, err := nameBuilder.String()
		return fmt.Errorf("%w, metrics: %s", err, name)
	default:
		err := crawlStruct(value, nameBuilder, m.unregisterMetricMetrics)
		if err != nil {
			return err
		}
	}
	return nil
}
