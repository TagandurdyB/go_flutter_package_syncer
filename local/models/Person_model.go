package models

import (
	// "encoding/json"
	"encoding/json"
	"strconv"

	// "gorm.io/driver/mysql"
	// "gorm.io/gorm"

	helpers "flutter_package_syncer/helpers"

)

type Person struct {
	ID          int     `json:"id"`
	LastName    string  `json:"last_name"`
	FirstName   string  `json:"first_name"`
	MidleName   string  `json:"midle_name"`
	Address     *string `json:"address"`
	WorkAddress *string `json:"work_address"`
	Phone       string  `json:"phone"`
	IDNumber    *string `json:"id_number"`
	EPCNumber   *string `json:"epc_number"`
	Photo       *string `json:"photo"`
	Birthday    *string `json:"birth_day"`
	// Birthday    *time.Time `json:"birth_day"`
}

// func _GetPersonDb() *gorm.DB {
// 	db, err := gorm.Open(mysql.Open(Dns), &gorm.Config{})
// 	if err != nil {
// 		panic(err)
// 	}
// 	return db
// }

const jsonFile = "./.base/persons.json"

func (person Person) Migrate() {
	// helpers.MkDir("./.base")
	// helpers.MkDir("./.base/uploads")
	// if !helpers.IsExist(jsonFile) {
	// 	helpers.CreateFile(jsonFile)
	// 	people := []Person{}
	// 	data, err := json.Marshal(people)
	// 	helpers.ErrH("Error in Person Migrate:", err)
	// 	helpers.WriteJson(jsonFile, data)
	// 	helpers.WriteFile("./.base/last_id.lock", "0")
	// }
}

func nilCheck(val *string) (result string) {
	if val != nil {
		result = *val
	} else {
		result = ""
	}
	return
}

func (person Person) ToArr() (arr []interface{}) {
	arr = append(arr, person.LastName)
	arr = append(arr, person.FirstName)
	arr = append(arr, person.MidleName)
	arr = append(arr, nilCheck(person.Birthday))
	arr = append(arr, nilCheck(person.Address))
	arr = append(arr, nilCheck(person.WorkAddress))
	arr = append(arr, person.Phone)
	arr = append(arr, nilCheck(person.IDNumber))
	arr = append(arr, nilCheck(person.EPCNumber))
	// arr = append(arr, person.Photo)
	return
}

func (person Person) Add() {
	lastId := helpers.ReadFile("./.base/last_id.lock")
	id, err := strconv.Atoi(lastId[0])
	id++
	helpers.ErrH("Error in Person Add get last ID:", err)
	people := Person{}.ReadAll()
	person.ID = id
	people = append(people, person)
	jsonArray, err := json.Marshal(people)
	helpers.ErrH("Error Add person:", err)
	helpers.WriteJson(jsonFile, jsonArray)
	strID := strconv.Itoa(id)
	helpers.WriteFile("./.base/last_id.lock", strID)

}

func (person Person) ReadAll() (people []Person) {
	jsonData := helpers.ReadAllJson(jsonFile)
	err := json.Unmarshal(jsonData, &people)
	helpers.ErrH("Error ReadAll person:", err)
	return
}

func (person Person) Read(id int) (result Person) {
	people := []Person{}
	jsonData := helpers.ReadAllJson(jsonFile)
	err := json.Unmarshal(jsonData, &people)
	helpers.ErrH("Error Read(", id, ") person:", err)
	for _, v := range people {
		if v.ID == id {
			result = v
			break
		}
	}
	return
}

func (person Person) Delete() {
	var newList []Person
	people := Person{}.ReadAll()
	for _, v := range people {
		if v.ID != person.ID {
			newList = append(newList, v)
		}
	}
	jsonArray, err := json.Marshal(newList)
	helpers.ErrH("Error Add person:", err)
	helpers.WriteJson(jsonFile, jsonArray)
}

// func (Person Person) Get(where ...interface{}) Person {
// 	db := _GetPersonDb()
// 	db.First(&Person, where...)
// 	return Person
// }

// func (Person Person) GetAll(where ...interface{}) []Person {
// 	db := _GetPersonDb()
// 	var Persons []Person
// 	db.Find(&Persons, where...)
// 	return Persons
// }

// func (Person Person) Update(column string, value interface{}) {
// 	db := _GetPersonDb()
// 	db.Model(&Person).Update(column, value)
// }

// func (Person Person) Updates(data Person) {
// 	db := _GetPersonDb()
// 	db.Model(&Person).Updates(data)
// }

// func (Person Person) Delete() {
// 	db := _GetPersonDb()
// 	db.Delete(&Person, Person.ID)
// }
