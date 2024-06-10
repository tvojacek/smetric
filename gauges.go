package smetric

import (
	"fmt"
	"sync"

	"github.com/VictoriaMetrics/metrics"
)

type Metrics struct {
	set           *metrics.Set
	metricsStruct interface{}
	nameBuilder   NameBuilder
	gauges        map[string]*metrics.Gauge
	lock          *sync.Mutex
}

//type Gauges struct {
//	set         *metrics.Set
//	gauges      map[string]*metrics.Gauge
//	nameBuilder NameBuilder
//	lock        *sync.Mutex
//}

// Name
func (g *Metrics) Name() string {
	err := g.isValid()
	if err != nil {
		return ""
	}
	return g.nameBuilder.LastSuffix()
}

// DestroyAll unregister and remove all gauges.
// Function is safe to call repeatably amd concurrently.
func (g *Metrics) DestroyAll() error {
	err := g.isValid()
	if err != nil {
		return err
	}
	g.lock.Lock()
	defer g.lock.Unlock()
	for name := range g.gauges {
		g.destroyGauge(name)
	}
	return nil
}

var ErrorDuplicateMetric = fmt.Errorf("metric already exist")
var ErrorNotExistMetric = fmt.Errorf("metric does not exist")

func (g *Metrics) isValid() error {
	if g == nil {
		return fmt.Errorf("nil Gauges")
	}
	if g.set == nil {
		return fmt.Errorf("uninicialized Gauges")
	}
	return nil
}

// Add register gauge. Function return error if gauge with same prefix already exist.
// Function is safe to call repeatably amd concurrently.
func (g *Metrics) AddGauge(suffix string, getValue func() float64) error {
	err := g.isValid()
	if err != nil {
		return err
	}
	g.lock.Lock()
	defer g.lock.Unlock()
	_, exist := g.gauges[suffix]
	if exist {
		return fmt.Errorf("%w: %s", ErrorDuplicateMetric, suffix)
	}
	g.add(suffix, getValue)
	return nil
}

func (g *Metrics) add(suffix string, getValue func() float64) error {
	name, err := g.nameBuilder.WithSuffix(suffix).String()
	if err != nil {
		return err
	}
	g.gauges[suffix] = g.set.NewGauge(name, getValue)
	return nil
}

// Destroy unregister gauge with suffix. Function return error if gauge does not exist.
// Function is safe to call repeatably amd concurrently.
func (g *Metrics) Destroy(suffix string) error {
	err := g.isValid()
	if err != nil {
		return err
	}
	g.lock.Lock()
	defer g.lock.Unlock()
	_, exist := g.gauges[suffix]
	if exist {
		g.destroyGauge(suffix)
	} else {
		return fmt.Errorf("%w: %s", ErrorNotExistMetric, suffix)
	}
	return nil
}

// AddOrReplace register gauge with suffix if Gauge with same prefix exist it is unregistered first.
// Function is safe to call repeatably amd concurrently.
func (g *Metrics) AddOrReplace(suffix string, getValue func() float64) error {
	err := g.isValid()
	if err != nil {
		return err
	}
	g.lock.Lock()
	defer g.lock.Unlock()
	_, exist := g.gauges[suffix]
	if exist {
		g.destroyGauge(suffix)
	}
	g.add(suffix, getValue)
	return nil
}

func (g *Metrics) destroyGauge(suffix string) error {
	name, err := g.nameBuilder.WithSuffix(suffix).String()
	if err != nil {
		return err
	}
	g.set.UnregisterMetric(name)
	delete(g.gauges, suffix)
	return nil
}
