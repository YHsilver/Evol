package evol

type Entity interface {
	// EntityIdentity for domain codec routing
	EntityIdentity() string
}
