package scraper

import (
	"fmt"
	"hash"
	"net/http"
	"sort"
	"time"

	"collectd.org/api"
	"collectd.org/plugin"
	"golang.org/x/crypto/blake2b"

	pcollectd "gitlab.in2p3.fr/rferrand/collectd-prometheus-plugin/collectd"

	//"collectd.org/plugin"
	"github.com/pkg/errors"
	dto "github.com/prometheus/client_model/go"
	pexpfmt "github.com/prometheus/common/expfmt"
)

// See https://github.com/prometheus/client_model/blob/master/go/metrics.pb.go

type PrometheusScraper struct {
	PluginName                 string
	MetaPrefix                 string
	FieldToHash                pcollectd.FieldType
	TypeInstanceOnlyHashedMeta bool
	HashLabelFunctionHashSize  int
	labelHasher                hash.Hash
}

func NewPrometheusScraper(pluginName string) *PrometheusScraper {

	typeInstanceOnlyForHashedMeta := true
	hasherHashSize := 8
	fieldToHash := "plugin_instance"
	var fieldToHashCollectdField pcollectd.FieldType

	switch fieldToHash {
	case "plugin_instance":
		fieldToHashCollectdField = pcollectd.PluginInstanceFieldType
	case "type_instance":
		fieldToHashCollectdField = pcollectd.TypeInstanceFieldType
	}

	hasher, _ := blake2b.New(hasherHashSize, nil)

	return &PrometheusScraper{
		PluginName:                 pluginName,
		MetaPrefix:                 fmt.Sprintf("%s.", pluginName),
		FieldToHash:                fieldToHashCollectdField,
		TypeInstanceOnlyHashedMeta: typeInstanceOnlyForHashedMeta,
		HashLabelFunctionHashSize:  hasherHashSize,
		labelHasher:                hasher,
	}
}

func (ps PrometheusScraper) getLabelName(name string) string {
	return fmt.Sprintf("%s%s", ps.MetaPrefix, name)
}

func (ps *PrometheusScraper) Read() error {
	return ps.Parse()
}

func (ps *PrometheusScraper) Parse() error {

	// td, err := os.Open("./testdata/traefik_2_metrics.txt")
	// if err != nil {
	// 	panic(err)
	// }
	// defer td.Close()
	hClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", "http://localhost:8082/metrics", nil)
	if err != nil {
		return errors.Wrap(err, "building new HTTP request")
	}

	resp, err := hClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "performing HTTP request")
	}
	defer resp.Body.Close()

	tparser := pexpfmt.TextParser{}

	metrics, err := tparser.TextToMetricFamilies(resp.Body)
	if err != nil {
		return errors.Wrap(err, "parsing prometheus metrics")
	}

	//putvalWriter := cformat.NewPutvalWithMeta(os.Stdout)

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
			//putvalWriter.Write(context.Background(), vl)
			if err := plugin.Write(vl); err != nil {
				return err
			}
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
		mMeta.Add(ps.getLabelName(label.GetName()), label.GetValue())
	}

	histogram := metric.Histogram

	newValueList := func(typeInstance string, value api.Value) *api.ValueList {

		pluginInstance := ps.computePluginInstance(mMeta, mf.GetName())
		cTypeInstance := ps.computeTypeInstance(mMeta, typeInstance)

		identifier := api.Identifier{
			Plugin:         ps.PluginName,
			PluginInstance: pluginInstance,
			Type:           value.Type(),
			TypeInstance:   cTypeInstance,
		}

		return &api.ValueList{
			Identifier: identifier,
			Time:       mTime,
			Values:     []api.Value{value},
			DSNames:    []string{"value"},
			Metadata:   extendMetadataWithIdentifier(mMeta, identifier),
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

	// TODO: FIXME, should iterate over Metrics
	metric := mf.Metric[0]

	mTime := promTimestampToTime(metric.TimestampMs)
	mMeta := make(api.Metadata)
	for _, label := range metric.GetLabel() {
		mMeta.Add(ps.getLabelName(label.GetName()), label.GetValue())
	}

	summary := metric.Summary

	newValueList := func(typeInstance string, value api.Value) *api.ValueList {

		pluginInstance := ps.computePluginInstance(mMeta, mf.GetName())
		cTypeInstance := ps.computeTypeInstance(mMeta, typeInstance)

		identifier := api.Identifier{
			Plugin:         ps.PluginName,
			PluginInstance: pluginInstance,
			Type:           value.Type(),
			TypeInstance:   cTypeInstance,
		}

		return &api.ValueList{
			Identifier: identifier,
			Time:       mTime,
			Values:     []api.Value{value},
			DSNames:    []string{"value"},
			Metadata:   extendMetadataWithIdentifier(mMeta, identifier),
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
			mMeta.Add(ps.getLabelName(label.GetName()), label.GetValue())
		}

		pluginInstance := ps.computePluginInstance(mMeta, mf.GetName())
		typeInstance := ps.computeTypeInstance(mMeta,
			fmt.Sprintf("%d", metricIdx))

		identifier := api.Identifier{
			Plugin:         ps.PluginName,
			PluginInstance: pluginInstance,
			Type:           "gauge",
			TypeInstance:   typeInstance,
		}

		vl := &api.ValueList{
			Identifier: identifier,
			Time:       mTime,
			Values:     []api.Value{api.Gauge(mValue)},
			DSNames:    []string{"value"},
			Metadata:   extendMetadataWithIdentifier(mMeta, identifier),
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
			mMeta.Add(ps.getLabelName(label.GetName()), label.GetValue())
		}

		pluginInstance := ps.computePluginInstance(mMeta, mf.GetName())
		typeInstance := ps.computeTypeInstance(mMeta,
			fmt.Sprintf("value%d", metricIdx))

		identifier := api.Identifier{
			Plugin:         ps.PluginName,
			PluginInstance: pluginInstance,
			Type:           "gauge",
			TypeInstance:   typeInstance,
		}

		vl := &api.ValueList{
			Identifier: identifier,
			Time:       mTime,
			Values:     []api.Value{api.Gauge(mValue)},
			DSNames:    []string{"value"},
			Metadata:   extendMetadataWithIdentifier(mMeta, identifier),
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
			mMeta.Add(ps.getLabelName(label.GetName()), label.GetValue())
		}

		pluginInstance := ps.computePluginInstance(mMeta, mf.GetName())
		typeInstance := ps.computeTypeInstance(mMeta, fmt.Sprintf("value%d", metricIdx))

		identifier := api.Identifier{
			Plugin:         ps.PluginName,
			PluginInstance: pluginInstance,
			Type:           "counter",
			TypeInstance:   typeInstance,
		}

		vl := &api.ValueList{
			Identifier: identifier,
			Time:       mTime,
			Values:     []api.Value{api.Counter(mValue)},
			DSNames:    []string{"value"},
			Metadata:   extendMetadataWithIdentifier(mMeta, identifier),
		}

		vls = append(vls, vl)
	}

	return vls, nil
}

func extractValueFromMetric(mftype dto.MetricType, metric *dto.Metric) api.Value {
	var value api.Value

	switch mftype {
	case dto.MetricType_GAUGE:
		value = api.Gauge(*(metric.Gauge.Value))
	case dto.MetricType_COUNTER:
		value = api.Counter(*(metric.Counter.Value))
	case dto.MetricType_UNTYPED:
		value = api.Gauge(*(metric.Untyped.Value))
		// case dto.MetricType_SUMMARY:
		//
		// case dto.MetricType_HISTOGRAM:
	}

	return value
}

func (ps PrometheusScraper) promSimpleValueToCollectdValueLists(mf *dto.MetricFamily) ([]*api.ValueList, error) {
	var vls []*api.ValueList

	for metricIdx, metric := range mf.GetMetric() {
		mTime := promTimestampToTime(metric.TimestampMs)
		mValue := extractValueFromMetric(mf.GetType(), metric)

		labelBasedMeta := make(api.Metadata)
		for _, label := range metric.GetLabel() {
			labelBasedMeta.Add(ps.getLabelName(label.GetName()), label.GetValue())
		}

		var pluginInstance, typeInstance string
		if ps.TypeInstanceOnlyHashedMeta {
			pluginInstance = mf.GetName()
			typeInstance = ps.hashMetadata(labelBasedMeta)
		} else {
			// the unique Hash can be concatenated to plugin_instance or
			// type_instance based on the configuration
			pluginInstance = ps.computePluginInstance(labelBasedMeta, mf.GetName())
			typeInstance = ps.computeTypeInstance(labelBasedMeta, fmt.Sprintf("value%d", metricIdx))
		}

		identifier := api.Identifier{
			Plugin:         ps.PluginName,
			PluginInstance: pluginInstance,
			Type:           mValue.Type(),
			TypeInstance:   typeInstance,
		}

		vl := &api.ValueList{
			Identifier: identifier,
			Time:       mTime,
			Values:     []api.Value{mValue},
			DSNames:    []string{"value"},
			Metadata:   extendMetadataWithIdentifier(labelBasedMeta, identifier),
		}

		vls = append(vls, vl)
	}

	return vls, nil
}

func (ps PrometheusScraper) computeTypeInstance(meta api.Metadata, wantedInstance string) string {
	return ps.computeInstance(pcollectd.TypeInstanceFieldType, meta, wantedInstance)
}

func (ps PrometheusScraper) computePluginInstance(meta api.Metadata, wantedInstance string) string {
	return ps.computeInstance(pcollectd.PluginInstanceFieldType, meta, wantedInstance)
}

// computePluginInstance takes a `wantedTypeInstance` as input
// and eventually concat a unique metadata hash
func (ps PrometheusScraper) computeInstance(fieldType pcollectd.FieldType, meta api.Metadata, wantedInstance string) string {
	if fieldType != ps.FieldToHash ||
		len(meta) == 0 {
		return wantedInstance
	}

	return fmt.Sprintf("%s_%s", wantedInstance, ps.hashMetadata(meta))
}

func (ps *PrometheusScraper) hashMetadata(meta api.Metadata) string {
	if len(meta) == 0 {
		return ""
	}

	ps.labelHasher.Reset()

	keys := sort.StringSlice(meta.Toc())
	keys.Sort()

	for _, mKey := range keys {
		s := mKey + meta.GetAsString(mKey)
		ps.labelHasher.Write([]byte(s))
	}

	return fmt.Sprintf("%x", ps.labelHasher.Sum(nil))
}

func extendMetadataWithIdentifier(meta api.Metadata,
	id api.Identifier) api.Metadata {
	var newMeta api.Metadata

	for _, key := range meta.Toc() {
		newMeta.Add(key, meta.Get(key))
	}

	newMeta.Add("prom.metric_name", id.PluginInstance)
	newMeta.Add("prom.type_instance", id.TypeInstance)

	return newMeta
}
