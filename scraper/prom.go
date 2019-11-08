package scraper

import (
	"context"
	"fmt"
	"os"

	"collectd.org/api"
	cformat "collectd.org/format"

	//"collectd.org/plugin"
	dto "github.com/prometheus/client_model/go"
	pexpfmt "github.com/prometheus/common/expfmt"
)

type PrometheusScraper struct {
	PluginName string
}

func NewPrometheusScraper() *PrometheusScraper {
	return &PrometheusScraper{}
}

func (ps *PrometheusScraper) Read() error {
	return ps.Parse()
}

func (ps *PrometheusScraper) Parse() error {

	td, err := os.Open("./testdata/traefik_2_metrics.txt")
	if err != nil {
		panic(err)
	}
	defer td.Close()

	tparser := pexpfmt.TextParser{}

	metrics, err := tparser.TextToMetricFamilies(td)
	if err != nil {
		panic(err)
	}

	putvalWriter := cformat.NewPutvalWithMeta(os.Stdout)

	var vls []*api.ValueList
	for _, mFamily := range metrics {
		// if mName != "http_request_duration_seconds" {
		// 	//if mName != "go_gc_duration_seconds" {
		// 	continue
		// }
		//fmt.Printf("[%s] %+v\n", mName, mFamily)

		switch mType := mFamily.GetType(); mType {
		case dto.MetricType_GAUGE:
			vls, err = ps.promGaugeToValueLists(mFamily)
		case dto.MetricType_COUNTER:
			vls, err = ps.promCounterToValueLists(mFamily)
		case dto.MetricType_SUMMARY:
			vls, err = ps.promSummaryToValueLists(mFamily)
		case dto.MetricType_HISTOGRAM:
			vls, err = ps.promHistogramToValueLists(mFamily)
		case dto.MetricType_UNTYPED:
			vls, err = ps.promUntypedToValueLists(mFamily)
		default:
			panic(fmt.Sprintf("unknown ptype %d", mType))

		}

		// for _, metric := range mFamily.Metric {
		// 	fmt.Printf("%+v\n", metric)
		// }

		//fmt.Printf("value-lists = %+v\n", vls)

		for _, vl := range vls {
			putvalWriter.Write(context.Background(), vl)
		}
	}

	return nil
}

func (ps PrometheusScraper) promHistogramToValueLists(mf *dto.MetricFamily) ([]*api.ValueList, error) {
	var vls []*api.ValueList

	metric := mf.Metric[0]
	mTime := promTimestampToTime(metric.TimestampMs)
	mMeta := make(api.Metadata)
	for _, label := range metric.GetLabel() {
		mMeta.Add(label.GetName(), label.GetValue())
	}

	histogram := metric.Histogram

	newValueList := func(typeInstance string, value api.Value) *api.ValueList {
		identifier := api.Identifier{
			Plugin:         ps.PluginName,
			PluginInstance: mf.GetName(),
			Type:           value.Type(),
			TypeInstance:   typeInstance,
		}

		return &api.ValueList{
			Identifier: identifier,
			Time:       mTime,
			Values:     []api.Value{value},
			DSNames:    []string{"value"},
			Metadata:   mMeta,
		}
	}

	vlCount := newValueList("sample_count", api.Counter(*histogram.SampleCount))
	vlSum := newValueList("sample_sum", api.Counter(*histogram.SampleSum))

	vls = append(vls, vlCount, vlSum)

	for _, bucket := range histogram.Bucket {
		typeInstance := fmt.Sprintf("bucket_%f", *bucket.UpperBound)
		vlHistogram := newValueList(typeInstance, api.Counter(*bucket.CumulativeCount))
		vls = append(vls, vlHistogram)
	}

	return vls, nil
}

func (ps PrometheusScraper) promSummaryToValueLists(mf *dto.MetricFamily) ([]*api.ValueList, error) {
	var vls []*api.ValueList

	metric := mf.Metric[0]

	mTime := promTimestampToTime(metric.TimestampMs)
	mMeta := make(api.Metadata)
	for _, label := range metric.GetLabel() {
		mMeta.Add(label.GetName(), label.GetValue())
	}

	summary := metric.Summary

	newValueList := func(typeInstance string, value api.Value) *api.ValueList {
		identifier := api.Identifier{
			Plugin:         ps.PluginName,
			PluginInstance: mf.GetName(),
			Type:           value.Type(),
			TypeInstance:   typeInstance,
		}

		return &api.ValueList{
			Identifier: identifier,
			Time:       mTime,
			Values:     []api.Value{value},
			DSNames:    []string{"value"},
			Metadata:   mMeta,
		}
	}

	vlCount := newValueList("sample_count", api.Counter(*summary.SampleCount))
	vlSum := newValueList("sample_sum", api.Counter(*summary.SampleSum))

	vls = append(vls, vlCount, vlSum)

	for _, quantile := range summary.Quantile {
		typeInstance := fmt.Sprintf("quantile_%f", *quantile.Quantile)
		vlQuantile := newValueList(typeInstance, api.Gauge(*quantile.Value))
		vls = append(vls, vlQuantile)
	}

	return vls, nil
}

func (ps PrometheusScraper) promUntypedToValueLists(mf *dto.MetricFamily) ([]*api.ValueList, error) {
	var vls []*api.ValueList

	for metricIdx, metric := range mf.GetMetric() {
		mTime := promTimestampToTime(metric.TimestampMs)
		mValue := *(metric.Untyped.Value)

		mMeta := make(api.Metadata)
		for _, label := range metric.GetLabel() {
			mMeta.Add(label.GetName(), label.GetValue())
		}

		identifier := api.Identifier{
			Plugin:         ps.PluginName,
			PluginInstance: mf.GetName(),
			Type:           "gauge",
			TypeInstance:   fmt.Sprintf("%d", metricIdx),
		}

		vl := &api.ValueList{
			Identifier: identifier,
			Time:       mTime,
			Values:     []api.Value{api.Gauge(mValue)},
			DSNames:    []string{"value"},
			Metadata:   mMeta,
		}

		vls = append(vls, vl)
	}

	return vls, nil
}

func (ps PrometheusScraper) promGaugeToValueLists(mf *dto.MetricFamily) ([]*api.ValueList, error) {
	var vls []*api.ValueList

	for metricIdx, metric := range mf.GetMetric() {
		mTime := promTimestampToTime(metric.TimestampMs)
		mValue := *(metric.Gauge.Value)

		mMeta := make(api.Metadata)
		for _, label := range metric.GetLabel() {
			mMeta.Add(label.GetName(), label.GetValue())
		}

		identifier := api.Identifier{
			Plugin:         ps.PluginName,
			PluginInstance: mf.GetName(),
			Type:           "gauge",
			TypeInstance:   fmt.Sprintf("value%d", metricIdx),
		}

		vl := &api.ValueList{
			Identifier: identifier,
			Time:       mTime,
			Values:     []api.Value{api.Gauge(mValue)},
			DSNames:    []string{"value"},
			Metadata:   mMeta,
		}

		vls = append(vls, vl)
	}

	return vls, nil
}

func (ps PrometheusScraper) promCounterToValueLists(mf *dto.MetricFamily) ([]*api.ValueList, error) {
	var vls []*api.ValueList

	for metricIdx, metric := range mf.GetMetric() {
		mTime := promTimestampToTime(metric.TimestampMs)
		mValue := *(metric.Counter.Value)

		mMeta := make(api.Metadata)
		for _, label := range metric.GetLabel() {
			mMeta.Add(label.GetName(), label.GetValue())
		}

		identifier := api.Identifier{
			Plugin:         ps.PluginName,
			PluginInstance: mf.GetName(),
			Type:           "counter",
			TypeInstance:   fmt.Sprintf("value%d", metricIdx),
		}

		vl := &api.ValueList{
			Identifier: identifier,
			Time:       mTime,
			Values:     []api.Value{api.Counter(mValue)},
			DSNames:    []string{"value"},
			Metadata:   mMeta,
		}

		vls = append(vls, vl)
	}

	return vls, nil
}
