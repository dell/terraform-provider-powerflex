/*
PowerFlex 5.x Gen2 REST API

API version: 5.1.0
*/

package clientgen

// MapVolumeToHostPayload payload for mapping a volume to a host
type MapVolumeToHostPayload struct {
	// Host ID to map to
	HostId string `json:"hostId"`
	// Access mode: READ_ONLY or READ_WRITE
	AccessMode string `json:"accessMode,omitempty"`
}

func (o MapVolumeToHostPayload) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["hostId"] = o.HostId
	if o.AccessMode != "" {
		toSerialize["accessMode"] = o.AccessMode
	}
	return toSerialize, nil
}

// UnmapVolumeFromHostPayload payload for unmapping a volume from a host
type UnmapVolumeFromHostPayload struct {
	// Host ID to unmap from
	HostId string `json:"hostId"`
}

func (o UnmapVolumeFromHostPayload) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["hostId"] = o.HostId
	return toSerialize, nil
}
