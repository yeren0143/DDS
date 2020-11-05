package common

type Property struct {
	Name      string
	Value     string
	Propagate bool
}

type PropertySeq = []Property
