package app

import (
	"log"

	"github.com/KurbanowS/Animal/internal/models"
	"github.com/KurbanowS/Animal/internal/store"
)

func AnimalList(f models.AnimalFilterRequest) ([]*models.AnimalResponse, int, error) {
	animals, total, err := store.Store().AnimalFindBy(f)
	log.Println(animals, "+++++++", total)
	if err != nil {
		return nil, 0, err
	}
	animalResponse := []*models.AnimalResponse{}
	for _, animal := range animals {
		a := models.AnimalResponse{}
		a.FromModel(animal)
		animalResponse = append(animalResponse, &a)
	}
	return animalResponse, total, nil
}

func AnimalCreate(data models.AnimalRequest) (*models.AnimalResponse, error) {
	model := &models.Animal{}
	data.ToModel(model)
	var err error
	model, err = store.Store().AnimalsCreate(model)
	if err != nil {
		return nil, err
	}
	res := &models.AnimalResponse{}
	res.FromModel(model)
	return res, nil
}
