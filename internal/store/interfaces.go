package store

import "github.com/KurbanowS/Animal/internal/models"

type IStore interface {
	AnimalFindById(ID string) (*models.Animal, error)
	AnimalFindByIds(Ids []string) ([]*models.Animal, error)
	AnimalFindBy(f models.AnimalFilterRequest) (animal []*models.Animal, total int, err error)
	AnimalsCreate(model *models.Animal) (*models.Animal, error)
}
