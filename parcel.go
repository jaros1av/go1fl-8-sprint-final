package main

import (
	"database/sql"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

const (
	InsertParcelTpl                = "INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :createdAt)"
	GetParselTpl                   = "SELECT number, client, status, address, created_at FROM parcel WHERE number = :number"
	GetParselListTpl               = "SELECT number, client, status, address, created_at FROM parcel WHERE client = :client"
	UpdateParselStatusByNumberTpl  = "UPDATE parcel SET status = :status WHERE number = :number"
	UpdateParselAddressByNumberTpl = "UPDATE parcel SET address = :address WHERE number = :number AND status = :status"
	DeleteParselByNumberTpl        = "DELETE FROM parcel WHERE number = :number AND status = :status"
)

func (s ParcelStore) Add(p Parcel) (int, error) {
	res, err := s.db.Exec(
		InsertParcelTpl,
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("createdAt", p.CreatedAt),
	)

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	p := Parcel{}

	row := s.db.QueryRow(GetParselTpl, sql.Named("number", number))
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)

	if err == sql.ErrNoRows {
		return Parcel{}, fmt.Errorf("parcel with number %d not found", number)
	}

	if err != nil {
		return Parcel{}, err
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	var res []Parcel

	rows, err := s.db.Query(GetParselListTpl, sql.Named("client", client))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		p := Parcel{}

		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return nil, err
		}

		res = append(res, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	res, err := s.db.Exec(
		UpdateParselStatusByNumberTpl,
		sql.Named("status", status),
		sql.Named("number", number))

	if err != nil {
		return err
	}

	affectedRow, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if affectedRow == 0 {
		return fmt.Errorf("parcel with number %d not found", number)
	}

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	_, err := s.db.Exec(
		UpdateParselAddressByNumberTpl,
		sql.Named("address", address),
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered),
	)

	return err
}

func (s ParcelStore) Delete(number int) error {
	_, err := s.db.Exec(
		DeleteParselByNumberTpl,
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered),
	)

	return err
}
