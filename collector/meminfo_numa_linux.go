// Copyright 2015 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build !nomeminfo_numa

package collector

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	memInfoNumaSubsystem = "memory_numa"
)

var nodeRE = regexp.MustCompile(`.*devices/system/node/node([0-9]*)`)

type meminfoKey struct {
	metricName, numaNode string
}

type meminfoNumaCollector struct {
	metricDescs map[string]*prometheus.Desc
}

func init() {
	Factories["meminfo_numa"] = NewMeminfoNumaCollector
}

// Takes a prometheus registry and returns a new Collector exposing
// memory stats.
func NewMeminfoNumaCollector() (Collector, error) {
	return &meminfoNumaCollector{
		metricDescs: map[string]*prometheus.Desc{},
	}, nil
}

func (c *meminfoNumaCollector) Update(ch chan<- prometheus.Metric) (err error) {
	gauges, counters, err := getMemInfoNuma()
	if err != nil {
		return fmt.Errorf("couldn't get NUMA meminfo: %s", err)
	}
	for k, v := range gauges {
		desc, ok := c.metricDescs[k.metricName]
		if !ok {
			desc = prometheus.NewDesc(
				prometheus.BuildFQName(Namespace, memInfoNumaSubsystem, k.metricName),
				fmt.Sprintf("Memory information field %s.", k.metricName),
				[]string{"node"}, nil)
			c.metricDescs[k.metricName] = desc
		}
		ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, v, k.numaNode)
	}
	for k, v := range counters {
		desc, ok := c.metricDescs[k.metricName]
		if !ok {
			desc = prometheus.NewDesc(
				prometheus.BuildFQName(Namespace, memInfoNumaSubsystem, k.metricName),
				fmt.Sprintf("Memory information field %s.", k.metricName),
				[]string{"node"}, nil)
			c.metricDescs[k.metricName] = desc
		}
		ch <- prometheus.MustNewConstMetric(desc, prometheus.CounterValue, v, k.numaNode)
	}
	return nil
}

func getMemInfoNuma() (map[meminfoKey]float64, map[meminfoKey]float64, error) {
	gaugeInfo := make(map[meminfoKey]float64)
	counterInfo := make(map[meminfoKey]float64)

	nodes, err := filepath.Glob(sysFilePath("devices/system/node/node[0-9]*"))
	if err != nil {
		return nil, nil, err
	}
	for _, node := range nodes {
		file, err := os.Open(path.Join(node, "meminfo"))
		if err != nil {
			return nil, nil, err
		}
		defer file.Close()

		numaInfo, err := parseMemInfoNuma(file)
		if err != nil {
			return nil, nil, err
		}
		for k, v := range numaInfo {
			gaugeInfo[k] = v
		}

		file, err = os.Open(path.Join(node, "numastat"))
		if err != nil {
			return nil, nil, err
		}

		nodeNumber := nodeRE.FindStringSubmatch(node)[1]
		numaStat, err := parseMemInfoNumaStat(file)
		if err != nil {
			return nil, nil, err
		}
		for k, v := range numaStat {
			counterInfo[meminfoKey{k, nodeNumber}] = v
		}
	}

	return gaugeInfo, counterInfo, nil
}

func parseMemInfoNuma(r io.Reader) (map[meminfoKey]float64, error) {
	var (
		memInfo = map[meminfoKey]float64{}
		scanner = bufio.NewScanner(r)
		re      = regexp.MustCompile("\\((.*)\\)")
	)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(string(line))

		fv, err := strconv.ParseFloat(parts[3], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value in meminfo: %s", err)
		}
		switch l := len(parts); {
		case l == 4: // no unit
		case l == 5 && parts[4] == "kB": // has unit
			fv *= 1024
		default:
			return nil, fmt.Errorf("invalid line in meminfo: %s", line)
		}
		metric := strings.TrimRight(parts[2], ":")

		// Active(anon) -> Active_anon
		metric = re.ReplaceAllString(metric, "_${1}")
		memInfo[meminfoKey{metric, parts[1]}] = fv
	}

	return memInfo, nil
}

func parseMemInfoNumaStat(r io.Reader) (map[string]float64, error) {
	var (
		numaStat = map[string]float64{}
		scanner  = bufio.NewScanner(r)
	)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(string(line))

		fv, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value in meminfo: %s", err)
		}

		numaStat[parts[0]] = fv
	}
	return numaStat, nil
}
