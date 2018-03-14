package main

import (
	"github.com/jhunt/go-firehose"
)

type Nozzle struct{}

func (nozzle *Nozzle) Reset() {
}

func (nozzle *Nozzle) Configure(c firehose.Config) {
}

func (nozzle *Nozzle) Track(e firehose.Event) {
	if e.GetOrigin() == "MetronAgent" {
		return
	}

	m := Metric{Type: e.Type()}
	switch m.Type {
	case firehose.ValueMetric:
		v := e.GetValueMetric()
		m.Name = e.GetOrigin()+"."+v.GetName()
		m.Value = v.GetValue()
		m.Unit = v.GetUnit()

		Lock.Lock()
		defer Lock.Unlock()
		Metrics[m.Name] = m

	case firehose.CounterEvent:
		c := e.GetCounterEvent()
		m.Name = e.GetOrigin()+"."+c.GetName()
		m.Tally = c.GetTotal()

		Lock.Lock()
		defer Lock.Unlock()
		Metrics[m.Name] = m
	}
}

func (nozzle *Nozzle) Flush() error {
	return nil
}

func (nozzle *Nozzle) SlowConsumer() {
}
