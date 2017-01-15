package model

type Challenge struct {
	ID            int64
	Name          string
	InternalName  string
	Image         string
	ImageTestFile string
	ImageTestCmd  string
	Description   string
	Points        string
	FromDB        bool
}
