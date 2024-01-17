package pgx

import (
	"context"
	"strconv"
	"strings"

	"github.com/KurbanowS/Animal/internal/models"
	"github.com/KurbanowS/Animal/internal/utils"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const sqlAnimalFields = `a.id, a.species, a.characteristic`
const sqlAnimalSelect = `select ` + sqlAnimalFields + ` from animal a where id = ANY($1::int[])`
const sqlAnimalSelectMany = `select ` + sqlAnimalFields + `, count(*) over() as total from animal a where a.id = a.id limit $1 offset $2`
const sqlAnimalInsert = `insert into animal`
const sqlAnimalUpdate = `update animal a set id = id`
const sqlAnimalDelete = `delete from animal where id = ANY($1::int[])`

func scanAnimals(rows pgx.Row, m *models.Animal, addColumns ...interface{}) (err error) {
	err = rows.Scan(parseColumnsForScan(m, addColumns...)...)
	return
}

func (d *PgxStore) AnimalFindById(ID string) (*models.Animal, error) {
	row, err := d.AnimalFindByIds([]string{ID})
	if err != nil {
		return nil, err
	}
	if len(row) < 1 {
		return nil, pgx.ErrNoRows
	}
	return row[0], nil
}

func (d *PgxStore) AnimalFindByIds(Ids []string) ([]*models.Animal, error) {
	animals := []*models.Animal{}
	err := d.runQuery(context.Background(), func(tx *pgxpool.Conn) (err error) {
		rows, err := tx.Query(context.Background(), sqlAnimalSelect, (Ids))
		for rows.Next() {
			m := models.Animal{}
			err := scanAnimals(rows, &m)
			if err != nil {
				utils.LoggerDesc("Scan error").Error(err)
				return err
			}
			animals = append(animals, &m)
		}
		return
	})
	if err != nil {
		utils.LoggerDesc("Query error").Error(err)
		return nil, err
	}
	return animals, nil
}

func (d *PgxStore) AnimalFindBy(f models.AnimalFilterRequest) (animals []*models.Animal, total int, err error) {
	args := []interface{}{f.Limit, f.Offset}
	qs, args := AnimalListBuildQuery(f, args)
	err = d.runQuery(context.Background(), func(tx *pgxpool.Conn) (err error) {
		rows, err := tx.Query(context.Background(), qs, args...)
		for rows.Next() {
			animal := models.Animal{}
			err = scanAnimals(rows, &animal, &total)
			if err != nil {
				utils.LoggerDesc("Scan error").Error(err)
				return err
			}
			animals = append(animals, &animal)
		}
		return
	})
	if err != nil {
		utils.LoggerDesc("Query error").Error(err)
		return nil, 0, err
	}
	return animals, total, nil
}

func (d *PgxStore) AnimalsCreate(model *models.Animal) (*models.Animal, error) {
	qs, args := AnimalsCreateQuery(model)
	qs += " RETURNING id"
	err := d.runQuery(context.Background(), func(tx *pgxpool.Conn) (err error) {
		err = tx.QueryRow(context.Background(), qs, args...).Scan(&model.ID)
		return
	})
	if err != nil {
		utils.LoggerDesc("Query error").Error(err)
		return nil, err
	}
	editModel, err := d.AnimalFindById(strconv.Itoa(int(model.ID)))
	if err != nil {
		return nil, err
	}
	return editModel, nil
}

func (d *PgxStore) AnimalUpdate(model *models.Animal) (*models.Animal, error) {
	qs, args := AnimalUpdateQuery(model)
	err := d.runQuery(context.Background(), func(tx *pgxpool.Conn) (err error) {
		_, err = tx.Query(context.Background(), qs, args...)
		return
	})
	if err != nil {
		utils.LoggerDesc("Query error").Error(err)
		return nil, err
	}
	editModel, err := d.AnimalFindById(strconv.Itoa(int(model.ID)))
	if err != nil {
		return nil, err
	}
	return editModel, nil
}

func (d *PgxStore) AnimalDelete(items []*models.Animal) ([]*models.Animal, error) {
	ids := []uint{}
	for _, i := range items {
		ids = append(ids, i.ID)
	}
	err := d.runQuery(context.Background(), func(tx *pgxpool.Conn) (err error) {
		_, err = tx.Query(context.Background(), sqlAnimalDelete, (ids))
		return
	})
	if err != nil {
		utils.LoggerDesc("Query error").Error(err)
		return nil, err
	}
	return items, nil
}

// ______________________________QUERIES______________________________________

func AnimalsCreateQuery(m *models.Animal) (string, []interface{}) {
	args := []interface{}{}
	cols := ""
	vals := ""
	q := AnimalsAtomicQuery(m)
	for k, v := range q {
		args = append(args, v)
		cols += ", " + k
		vals += ", $" + strconv.Itoa(len(args))
	}
	qs := sqlAnimalInsert + " (" + strings.Trim(cols, ", ") + ") VALUES (" + strings.Trim(vals, ", ") + ")"
	return qs, args
}

func AnimalsAtomicQuery(m *models.Animal) map[string]interface{} {
	q := map[string]interface{}{}
	q["species"] = m.Species
	q["characteristic"] = m.Characteristic
	return q
}

func AnimalListBuildQuery(f models.AnimalFilterRequest, args []interface{}) (string, []interface{}) {
	var wheres string = ""

	if f.ID != nil && *f.ID != 0 {
		args = append(args, *f.ID)
		wheres += "and a.id=$" + strconv.Itoa(len(args))
	}
	wheres += "order by a.id desc"
	qs := sqlAnimalSelectMany
	qs = strings.ReplaceAll(qs, "a.id=a.id", "a.id=a.id "+wheres+" ")
	return qs, args
}

func AnimalUpdateQuery(m *models.Animal) (string, []interface{}) {
	args := []interface{}{}
	sets := ""
	q := AnimalsAtomicQuery(m)
	for k, v := range q {
		args = append(args, v)
		sets += ", " + k + "=$" + strconv.Itoa(len(args))
	}
	args = append(args, m.ID)
	qs := strings.ReplaceAll(sqlAnimalUpdate, "set id=id", "set id=id "+sets+"") + "where id=" + strconv.Itoa(len(args))
	return qs, args
}
