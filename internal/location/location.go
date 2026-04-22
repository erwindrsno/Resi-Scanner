package location

type Location struct {
	Id   int
	Code string
}

type LocationRepository interface {
	Save(l *Location) error
}
