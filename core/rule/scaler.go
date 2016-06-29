package rule

import "fmt"

// Scaler control the service
type Scaler interface {
	Describe() string
	Up() error
	Down() error
}

// MockService is a wrapper for the MockScaler
func MockService(name string) Service {
	return Service{
		Name:  name,
		Scale: new(MockScaler),
	}
}

// MockScaler write "scale up" or "scale down" to stdout
type MockScaler struct{}

// Describe scaler
func (s *MockScaler) Describe() string {
	return "A mock scaler writing to stdout"
}

// Up mock
func (s *MockScaler) Up() error {
	fmt.Println("SCALE UP")
	return nil
}

// Down mock
func (s *MockScaler) Down() error {
	fmt.Println("SCALE DOWN")
	return nil
}
