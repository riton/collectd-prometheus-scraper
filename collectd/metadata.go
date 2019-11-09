package collectd

import "collectd.org/api"

func ExtendMetadataWithKeyValue(parent api.Metadata, key string, value interface{}) api.Metadata {
	newMeta := make(api.Metadata)
	newMeta.Set(key, value)
	return parent.CloneMerge(newMeta)
}
