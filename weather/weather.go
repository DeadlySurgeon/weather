package weather

// Service exposes the different functions we can execute on the weather
// service.
//
//go:generate mockgen -source=weather.go -package=mock -destination=mock/weather.go
type Service interface {
	At(lat string, long string) (Report, error)
}
