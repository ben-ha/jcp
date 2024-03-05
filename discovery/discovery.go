package discovery

type Discoverer interface {
	Next() (FileInformation, error)
}
