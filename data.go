package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/jhunt/go-firehose"
	"github.com/jhunt/go-log"
)

var (
	/* deployment/job -> index -> name -> value */
	Metrics map[string]map[string]map[string]Metric
	Lock    sync.Mutex
)

func init() {
	Metrics = make(map[string]map[string]map[string]Metric)
}

func Prune(expiry int) {
	Lock.Lock()
	defer Lock.Unlock()

	log.Debugf("pruning metrics that haven't been seen in %d seconds", expiry)
	deadline := time.Now().Unix() - int64(expiry)
	for job := range Metrics {
		for idx := range Metrics[job] {
			for name, m := range Metrics[job][idx] {
				if m.LastSeen < deadline {
					log.Infof("removing [%s] metric from [%s][%s]", name, job, idx)
					delete(Metrics[job][idx], name)
				}
			}
			if len(Metrics[job][idx]) == 0 {
				log.Infof("removing [%s] idx from [%s]", idx, job)
				delete(Metrics[job], idx)
			}
		}
		if len(Metrics[job]) == 0 {
			log.Infof("removing [%s]", job)
			delete(Metrics, job)
		}
	}
}

type Metric struct {
	Type     firehose.EventType
	Name     string
	LastSeen int64

	Value float64
	Unit  string

	Tally uint64
}

func (m Metric) String() string {
	switch m.Type {
	case firehose.ValueMetric:
		return fmt.Sprintf("value %f (%s)", m.Value, m.Unit)
	case firehose.CounterEvent:
		return fmt.Sprintf("counter %d", m.Tally)
	default:
		return "((unknown type))"
	}
}

func (m *Metric) UnmarshalJSON(b []byte) error {
	return fmt.Errorf("unmarshalling of metrics is not yet implemented!")
}

func (m Metric) MarshalJSON() ([]byte, error) {
	switch m.Type {
	case firehose.ValueMetric:
		out := struct {
			Type  string  `json:"type"`
			Name  string  `json:"name"`
			Value float64 `json:"value"`
			Unit  string  `json:"unit"`
		}{
			Type:  "value",
			Name:  m.Name,
			Value: m.Value,
			Unit:  m.Unit,
		}
		return json.Marshal(out)

	case firehose.CounterEvent:
		out := struct {
			Type  string `json:"type"`
			Name  string `json:"name"`
			Value uint64 `json:"value"`
		}{
			Type:  "counter",
			Name:  m.Name,
			Value: m.Tally,
		}
		return json.Marshal(out)
	}

	return nil, fmt.Errorf("cannot marshal this metric type (not a value or a counter)!")
}
