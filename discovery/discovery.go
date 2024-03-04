package discovery

type Discoverer interface {
	Next() (string, error)
}
