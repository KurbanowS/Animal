package models

type Animal struct {
	ID             uint    `json:"id"`
	Species        *string `json:"species"`
	Characteristic *string `json:"characteristic"`
}

func (Animal) RelationFields() []string {
	return []string{}
}

type AnimalRequest struct {
	ID             *uint   `json:"id"`
	Species        *string `json:"species"`
	Characteristic *string `json:"characteristic"`
}

type AnimalResponse struct {
	ID             uint    `json:"id"`
	Species        *string `json:"species"`
	Characteristic *string `json:"characteristic"`
}

func (b *AnimalRequest) ToModel(m *Animal) {
	if b.ID != nil {
		m.ID = *b.ID
	}
	m.Species = b.Species
	m.Characteristic = b.Characteristic
}

func (r *AnimalResponse) FromModel(m *Animal) {
	r.ID = m.ID
	r.Species = m.Species
	r.Characteristic = m.Characteristic
}

type AnimalFilterRequest struct {
	ID *uint `form:"id"`
	PaginationRequest
}
