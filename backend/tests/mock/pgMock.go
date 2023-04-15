package mock

// PgMock implements db.PgCaller interface
type PgMock struct {
	Urls   []Url
	Users  []User
	Clicks []ClickInfo
}

// NewPgMock returns new PgMock
func NewPgMock() *PgMock {
	return &PgMock{
		Urls:   []Url{},
		Users:  []User{},
		Clicks: []ClickInfo{},
	}
}
