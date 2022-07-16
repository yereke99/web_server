package repository

import (
	"pro1/entity"

	"github.com/jinzhu/gorm"
)

type DataRepository interface {
	InsertData(d entity.Data) entity.Data
	UpdateData(d entity.Data) entity.Data
	DeleteData(d entity.Data)
	AllData() []entity.Data
	FindByDataID(dataID uint64) entity.Data
}

type dataConnection struct {
	connection *gorm.DB
}

func NewDataRepository(db *gorm.DB) DataRepository {
	return &dataConnection{
		connection: db,
	}
}

func (db *dataConnection) InsertData(d entity.Data) entity.Data {
	db.connection.Save(&d)
	db.connection.Preload("User").Find(&d)
	return d
}

func (db *dataConnection) UpdateData(d entity.Data) entity.Data {
	db.connection.Save(&d)
	db.connection.Preload("User").Find(&d)
	return d
}

func (db *dataConnection) DeleteData(d entity.Data) {
	db.connection.Delete(&d)
}


func (db *dataConnection) AllData() []entity.Data {
	var data []entity.Data
	db.connection.Preload("User").Find(&data)
	return data
}

func (db *dataConnection) FindByDataID(dataID uint64) entity.Data {
	var data entity.Data
	db.connection.Preload("User").Find(&data, dataID)
	return data
}
