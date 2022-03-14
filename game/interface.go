package game

type Game interface {
	ReviveCell(int, int) error
	KillCell(int, int) error
	Evolve()
	GetCell(int, int) (*bool, error)
	GetGeneration() *generation
}
