/*
PowerFlex 5.x Gen2 REST API

API version: 5.1.0
*/

package clientgen

// RestoreVolumePayload payload for restoring a volume from snapshot
type RestoreVolumePayload struct {
	// Source snapshot ID to restore from
	SourceSnapshotId string `json:"sourceSnapshotId"`
}

func (o RestoreVolumePayload) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["sourceSnapshotId"] = o.SourceSnapshotId
	return toSerialize, nil
}
