package scraper

import (
	pcollectd "gitlab.in2p3.fr/rferrand/collectd-prometheus-plugin/collectd"
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
