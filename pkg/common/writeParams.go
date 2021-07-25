package common

// WriteParams contains additional information of a CacheChange.
type WriteParamsT struct {
	SampleIdentity         SampleIdentityT
	ReleatedSampleIdentity SampleIdentityT
}

var KWriteParamDefault WriteParamsT
