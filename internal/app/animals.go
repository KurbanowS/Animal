package app

import (
	"errors"
	"log"
	"strings"

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

func AnimalUpdate(data models.AnimalRequest) (*models.AnimalResponse, error) {
	model := &models.Animal{
		Species:        data.Species,
		Characteristic: data.Characteristic,
	}
	data.ToModel(model)

	var err error
	model, err = store.Store().AnimalUpdate(model)
	if err != nil {
		return nil, err
	}
	res := &models.AnimalResponse{}
	res.FromModel(model)
	return res, nil
}

func AnimalDelete(ids []string) ([]*models.AnimalResponse, error) {
	animals, err := store.Store().AnimalFindByIds(ids)
	if err != nil {
		return nil, err
	}
	if len(animals) < 1 {
		return nil, errors.New("model not found: " + strings.Join(ids, ","))
	}
	animals, err = store.Store().AnimalDelete(animals)
	if err != nil {
		return nil, err
	}
	if len(animals) == 0 {
		return make([]*models.AnimalResponse, 0), nil
	}

	var animalsResponse []*models.AnimalResponse
	for _, animal := range animals {
		var a models.AnimalResponse
		a.FromModel(animal)
		animalsResponse = append(animalsResponse, &a)
	}
	return animalsResponse, nil
}
