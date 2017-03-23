package entity

type Photo struct {
	Key   string `datastore:",noindex"`
	Proto []byte `datastore:",noindex"`
}

type Comparison struct {
	Proto []byte `datastore:",noindex"`
}
