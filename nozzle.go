package main

import (
	"github.com/jhunt/go-firehose"
	"github.com/jhunt/go-log"
)

type Nozzle struct{}

func (nozzle *Nozzle) set(job, idx string, m Metric) {
	Lock.Lock()
	defer Lock.Unlock()

	if _, ok := Metrics[job]; !ok {
		log.Infof("extending metrics to include [%s]...", job)
		Metrics[job] = make(map[string] map[string] Metric)
	}

	if _, ok := Metrics[job][idx]; !ok {
		log.Infof("extending metrics[%s] to include index [%s]...", job, idx)
		Metrics[job][idx] = make(map[string] Metric)
	}

	log.Debugf("ingesting [%s][%s][%s] = %s", job, idx, m.Name, m)
	Metrics[job][idx][m.Name] = m
}

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

		nozzle.set(e.GetJob(), e.GetIndex(), m)

	case firehose.CounterEvent:
		c := e.GetCounterEvent()
		m.Name = e.GetOrigin()+"."+c.GetName()
		m.Tally = c.GetTotal()

		nozzle.set(e.GetJob(), e.GetIndex(), m)
	}
}

func (nozzle *Nozzle) Flush() error {
	return nil
}

func (nozzle *Nozzle) SlowConsumer() {
}
