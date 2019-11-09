package scraper

import (
	"context"
	"fmt"
	"hash"
	"net/http"
	"sort"
	"time"

	"collectd.org/api"
	"golang.org/x/crypto/blake2b"

	pcollectd "gitlab.in2p3.fr/rferrand/collectd-prometheus-plugin/collectd"

	//"collectd.org/plugin"
	"github.com/pkg/errors"
	dto "github.com/prometheus/client_model/go"
	pexpfmt "github.com/prometheus/common/expfmt"
)

// See https://github.com/prometheus/client_model/blob/master/go/metrics.pb.go

const (
	// promMetaPrefix is the Metadata prefix used for prometheus specific
	// metadatas
	promMetaPrefix = "prom."
)

var (
	backgroundCtx = context.Background()
)

type PrometheusScraper struct {
	PluginName string
	MetaPrefix string

	// TargetURL is the destination of the scraper
	TargetURL string

	// HTTPTimeout is the timeout used when performing the GET
	// on TargetURL
	HTTPTimeout time.Duration

	// FieldToHash is the collectd field the unique metadata Hash
	// should be concatenated with. Supported values are
	// PluginInstanceFieldType or TypeInstanceFieldType
	// This is of no use if `TypeInstanceOnlyHashedMeta` is set to `true`
	FieldToHash pcollectd.FieldType

	// TypeInstanceOnlyHashedMeta defines how prometheus values
	// should be mapped to collectd namespace
	// If TypeInstanceOnlyHashedMeta is set to `true`, the `type_instance`
	// value will only contain a hashed version of the value metadatas
	// and no usable information
	TypeInstanceOnlyHashedMeta bool

	// HashLabelFunctionHashSize is the size of the hash used
	// to ensure unicity of values in the way collectd knows it
	// (metadata are not considered)
	// This size can be a value between 1 and 64 but it is highly
	// recommended to use values equal or greater than 32
	HashLabelFunctionHashSize int

	// AdditionalMetadata is a set of metadata that should be added
	// to every metric dispatched
	// Note: those metadata does not inherits the MetaPrefix value
	AdditionalMetadata api.Metadata

	labelHasher hash.Hash
	valueWriter api.Writer
	httpClient  httpDoer
}

func NewPrometheusScraper(pluginName string) *PrometheusScraper {

	additionalMetadata := make(api.Metadata)
	//additionalMetadata.Set("api-stable", false)

	targetURL := "http://coredns:9253/metrics"
	//targetURL := "http://traefik:8082/metrics"
	httpTimeout := 5 * time.Second

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

	valueWriter := pcollectd.NewFileWriter("/tmp/value_list.txt")

	return &PrometheusScraper{
		PluginName:                 pluginName,
		TargetURL:                  targetURL,
		HTTPTimeout:                httpTimeout,
		MetaPrefix:                 fmt.Sprintf("%s.", pluginName),
		FieldToHash:                fieldToHashCollectdField,
		TypeInstanceOnlyHashedMeta: typeInstanceOnlyForHashedMeta,
		HashLabelFunctionHashSize:  hasherHashSize,
		AdditionalMetadata:         additionalMetadata,
		labelHasher:                hasher,
		valueWriter:                valueWriter,
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
		Timeout: ps.HTTPTimeout,
	}

	req, err := http.NewRequest("GET", ps.TargetURL, nil)
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
		case dto.MetricType_GAUGE, dto.MetricType_UNTYPED, dto.MetricType_COUNTER:
			vls, err = ps.promSimpleValueToCollectdValueLists(mFamily)
		case dto.MetricType_SUMMARY, dto.MetricType_HISTOGRAM:
			vls, err = ps.promCompoundedValueToCollectdValueLists(mFamily)
		default:
			panic(fmt.Sprintf("unknown ptype %d", mType))
		}

		// for _, metric := range mFamily.Metric {
		// 	fmt.Printf("%+v\n", metric)
		// }

		//fmt.Printf("value-lists = %+v\n", vls)

		for _, vl := range vls {
			//putvalWriter.Write(context.Background(), vl)

			if len(ps.AdditionalMetadata) > 0 {
				nMeta := vl.Metadata.CloneMerge(ps.AdditionalMetadata)
				vl.Metadata = nMeta
			}

			if err := ps.valueWriter.Write(backgroundCtx, vl); err != nil {
				return err
			}
		}
	}

	return nil
}

func (ps PrometheusScraper) newCollectdValueList(name string, mTime time.Time, wTypeInstance string,
	value api.Value, meta api.Metadata) *api.ValueList {

	var pluginInstance, typeInstance string
	if ps.TypeInstanceOnlyHashedMeta {
		pluginInstance = name
		typeInstance = ps.hashMetadata(meta)
	} else {
		// the unique Hash can be concatenated to plugin_instance or
		// type_instance based on the configuration
		pluginInstance = ps.computePluginInstance(meta, name)
		typeInstance = ps.computeTypeInstance(meta, wTypeInstance)
	}

	identifier := api.Identifier{
		Plugin:         ps.PluginName,
		PluginInstance: pluginInstance,
		Type:           value.Type(),
		TypeInstance:   typeInstance,
	}

	return &api.ValueList{
		Identifier: identifier,
		Time:       mTime,
		Values:     []api.Value{value},
		DSNames:    []string{"value"},
		Metadata:   meta,
	}

}

type collectdValueListGeneratorFnc func(string, api.Value, api.Metadata) *api.ValueList

func (ps PrometheusScraper) newCollectdCompoundedValueListFnc(mName string, metric *dto.Metric) collectdValueListGeneratorFnc {
	metricTime := promTimestampToTime(metric.TimestampMs)

	return func(typeInstance string, value api.Value, meta api.Metadata) *api.ValueList {
		return ps.newCollectdValueList(mName, metricTime, typeInstance, value, meta)
	}
}

func (ps PrometheusScraper) promCompoundedValueToCollectdValueLists(mf *dto.MetricFamily) ([]*api.ValueList, error) {
	var vls []*api.ValueList

	for _, metric := range mf.GetMetric() {

		labelBasedMeta := ps.metadataFromMetricLabels(metric.GetLabel())
		newCollectdValueListFnc := ps.newCollectdCompoundedValueListFnc(mf.GetName(), metric)

		var typeVl []*api.ValueList
		switch mf.GetType() {
		case dto.MetricType_SUMMARY:
			typeVl = ps.promSummaryMetricToValueLists(newCollectdValueListFnc, metric.Summary, labelBasedMeta)
		case dto.MetricType_HISTOGRAM:
			typeVl = ps.promHistogramMetricToValueLists(newCollectdValueListFnc, metric.Histogram, labelBasedMeta)
		}

		vls = append(vls, typeVl...)
	}

	return vls, nil
}

func (ps PrometheusScraper) promHistogramMetricToValueLists(newValueListFnc collectdValueListGeneratorFnc,
	histogram *dto.Histogram, labelBasedMeta api.Metadata) []*api.ValueList {

	var vls []*api.ValueList

	vlCountMeta := pcollectd.ExtendMetadataWithKeyValue(labelBasedMeta,
		metadataKeyWithPromPrefix("bucket.sample_count"), true)
	vlCount := newValueListFnc("sample_count", api.Counter(*histogram.SampleCount),
		vlCountMeta)

	vlSumMeta := pcollectd.ExtendMetadataWithKeyValue(labelBasedMeta,
		metadataKeyWithPromPrefix("bucket.sample_sum"), true)
	vlSum := newValueListFnc("sample_sum", api.Counter(*histogram.SampleSum),
		vlSumMeta)

	vls = append(vls, vlCount, vlSum)

	for _, bucket := range histogram.Bucket {

		// set bucket metadata under the dedicated prometheus metadata namespace
		vlMeta := pcollectd.ExtendMetadataWithKeyValue(labelBasedMeta,
			metadataKeyWithPromPrefix("bucket.upper_bound"), *bucket.UpperBound /* float64 */)

		typeInstance := fmt.Sprintf("bucket_%f", *bucket.UpperBound) // used only if TypeInstanceOnlyHashedMeta is set to `false`
		vlHistogram := newValueListFnc(typeInstance, api.Counter(*bucket.CumulativeCount),
			vlMeta)
		vls = append(vls, vlHistogram)
	}

	return vls
}

func (ps PrometheusScraper) promSummaryMetricToValueLists(newValueListFnc collectdValueListGeneratorFnc,
	summary *dto.Summary, labelBasedMeta api.Metadata) []*api.ValueList {

	var vls []*api.ValueList

	vlCountMeta := pcollectd.ExtendMetadataWithKeyValue(labelBasedMeta,
		metadataKeyWithPromPrefix("quantile.sample_count"), true)
	vlCount := newValueListFnc("sample_count", api.Counter(*summary.SampleCount),
		vlCountMeta)

	vlSumMeta := pcollectd.ExtendMetadataWithKeyValue(labelBasedMeta,
		metadataKeyWithPromPrefix("quantile.sample_sum"), true)
	vlSum := newValueListFnc("sample_sum", api.Counter(*summary.SampleSum),
		vlSumMeta)

	vls = append(vls, vlCount, vlSum)

	for _, quantile := range summary.Quantile {

		// set quantile metadata under the dedicated prometheus metadata namespace
		vlMeta := pcollectd.ExtendMetadataWithKeyValue(labelBasedMeta,
			metadataKeyWithPromPrefix("quantile.quantile"), *quantile.Quantile /* float64 */)

		typeInstance := fmt.Sprintf("quantile_%f", *quantile.Quantile) // used only if TypeInstanceOnlyHashedMeta is set to `false`

		vlQuantile := newValueListFnc(typeInstance, api.Gauge(*quantile.Value),
			vlMeta)

		vls = append(vls, vlQuantile)
	}

	return vls
}

func extractSimpleValueFromMetric(mftype dto.MetricType, metric *dto.Metric) api.Value {
	var value api.Value

	switch mftype {
	case dto.MetricType_GAUGE:
		value = api.Gauge(*(metric.Gauge.Value))
	case dto.MetricType_COUNTER:
		value = api.Counter(*(metric.Counter.Value))
	case dto.MetricType_UNTYPED:
		value = api.Gauge(*(metric.Untyped.Value))
	}

	return value
}

func (ps PrometheusScraper) promSimpleValueToCollectdValueLists(mf *dto.MetricFamily) ([]*api.ValueList, error) {
	var vls []*api.ValueList

	for metricIdx, metric := range mf.GetMetric() {
		mTime := promTimestampToTime(metric.TimestampMs)
		mValue := extractSimpleValueFromMetric(mf.GetType(), metric)

		labelBasedMeta := ps.metadataFromMetricLabels(metric.GetLabel())

		typeInstance := fmt.Sprintf("value%d", metricIdx) // value used if TypeInstanceOnlyHashedMeta is `false`
		vl := ps.newCollectdValueList(mf.GetName(), mTime, typeInstance, mValue,
			labelBasedMeta)

		vls = append(vls, vl)
	}

	return vls, nil
}

func (ps PrometheusScraper) metadataFromMetricLabels(labels []*dto.LabelPair) api.Metadata {
	labelBasedMeta := make(api.Metadata)
	if len(labels) > 0 {
		for _, label := range labels {
			labelBasedMeta.Set(ps.getLabelName(label.GetName()), label.GetValue())
		}
	}
	return labelBasedMeta
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

func (ps PrometheusScraper) metaKeyWithPrefix(key string) string {
	return ps.MetaPrefix + key
}

func metadataKeyWithPromPrefix(key string) string {
	return promMetaPrefix + key
}
