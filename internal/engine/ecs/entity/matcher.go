package entity

type Matcher interface {
	Matches(e Entity) bool
}
