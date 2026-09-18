package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	// randSource источник псевдо случайных чисел.
	// Для повышения уникальности в качестве seed
	// используется текущее время в unix формате (в виде числа)
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange использует randSource для генерации случайных чисел
	randRange = rand.New(randSource)
)

// getTestParcel возвращает тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	id, err := store.Add(parcel)
	require.NoError(t, err, "error add parsel")
	require.NotZero(t, id, "parcel id must be set")

	// get
	parselFromDb, err := store.Get(id)
	require.NoError(t, err, "error get parsel")

	assert.Equal(t, parcel.Client, parselFromDb.Client, "expected: %d, got: %d", parcel.Client, parselFromDb.Client)
	assert.Equal(t, parcel.Status, parselFromDb.Status, "expected: %s, got: %s", parcel.Status, parselFromDb.Status)
	assert.Equal(t, parcel.Address, parselFromDb.Address, "expected: %s, got: %s", parcel.Address, parselFromDb.Address)
	assert.Equal(t, parcel.CreatedAt, parselFromDb.CreatedAt, "expected: %s, got: %s", parcel.CreatedAt, parselFromDb.CreatedAt)

	// delete
	err = store.Delete(id)
	require.NoError(t, err, "dell error")

	_, err = store.Get(id)
	assert.Error(t, err, "parcel with number %d must be deleted", id)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	id, err := store.Add(parcel)
	require.NoError(t, err, "error add parsel")
	require.NotZero(t, id, "parcel id must be set")

	// set address
	newAddress := "new test address"

	err = store.SetAddress(id, newAddress)
	require.NoError(t, err, "set address error")

	// check
	parselFromDb, err := store.Get(id)
	require.NoError(t, err, "error get parsel")
	assert.Equal(t, newAddress, parselFromDb.Address, "expected: %d, got: %d", newAddress, parselFromDb.Address)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	id, err := store.Add(parcel)
	require.NoError(t, err, "error add parsel")
	require.NotZero(t, id, "parcel id must be set")

	// set status
	err = store.SetStatus(id, ParcelStatusDelivered)
	require.NoError(t, err, "set status error")

	// check
	parselFromDb, err := store.Get(id)
	require.NoError(t, err, "error get parsel")
	assert.Equal(t, ParcelStatusDelivered, parselFromDb.Status, "expected: %d, got: %d", ParcelStatusDelivered, parselFromDb.Status)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	// задаём всем посылкам один и тот же идентификатор клиента
	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	// add
	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i])
		require.NoError(t, err, "error add parsel")
		require.NotZero(t, id, "parcel id must be set")
		// обновляем идентификатор добавленной у посылки
		parcels[i].Number = id

		parcelMap[id] = parcels[i]
	}

	// get by client
	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err, "get pasrel list error")

	// check
	for _, parcel := range storedParcels {
		expectedParsel, ok := parcelMap[parcel.Number]
		require.True(t, ok, "parcel with number %d not faund in parcelMap", parcel.Number)

		assert.Equal(t, parcel.Number, expectedParsel.Number, "expected: %d, got: %d", parcel.Number, expectedParsel.Number)
		assert.Equal(t, parcel.Client, expectedParsel.Client, "expected: %d, got: %d", parcel.Client, expectedParsel.Client)
		assert.Equal(t, parcel.Status, expectedParsel.Status, "expected: %s, got: %s", parcel.Status, expectedParsel.Status)
		assert.Equal(t, parcel.Address, expectedParsel.Address, "expected: %s, got: %s", parcel.Address, expectedParsel.Address)
		assert.Equal(t, parcel.CreatedAt, expectedParsel.CreatedAt, "expected: %s, got: %s", parcel.CreatedAt, expectedParsel.CreatedAt)
	}
}
