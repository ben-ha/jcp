package discovery

type bfsDiscoverer struct {
	BasePath string
}

func MakeBfsDiscoverer(basePath string) Discoverer {
	return bfsDiscoverer{BasePath: basePath}
}
