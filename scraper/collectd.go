package scraper

import (
	pcollectd "github.com/ccin2p3/collectd-prometheus-plugin/collectd"
)

func fieldToHashStringToCollectdFieldType(fieldName string) pcollectd.FieldType {
	switch fieldName {
	case "plugin_instance":
		return pcollectd.PluginInstanceFieldType
	case "type_instance":
		return pcollectd.TypeInstanceFieldType
	}
	return pcollectd.TypeInstanceFieldType
}
