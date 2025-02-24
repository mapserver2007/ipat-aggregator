package repository

type PathOptimizer interface {
	GetProjectRoot() (string, error)
}
