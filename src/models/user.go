package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Email     *string   `gorm:"size:255;unique;not null"`
	Username  string    `gorm:"size:255;unique;not null"`
	Password  string    `gorm:"size:255;not null"`
	FirstName string    `gorm:"size:255;not null"`
	LastName  string    `gorm:"size:255;not null"`
	NickName  string    `gorm:"size:512;default:NULL"`
	CreatedAt time.Time `sql:"default:ON UPDATE CURRENT_TIMESTAMP" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `sql:"default:NULL" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

//create a user
func CreateUser(db *gorm.DB, User *User) (err error) {
	err = db.Create(User).Error
	if err != nil {
		return err
	}
	return nil
}

//get users
func GetUsers(db *gorm.DB, User *[]User) (err error) {
	err = db.Find(User).Error
	if err != nil {
		return err
	}
	return nil
}

//get user by id
func GetUser(db *gorm.DB, User *User, id string) (err error) {
	err = db.Where("id = ?", id).First(User).Error
	if err != nil {
		return err
	}
	return nil
}

//update user
func UpdateUser(db *gorm.DB, User *User) (err error) {
	db.Save(User)
	return nil
}

//delete user
func DeleteUser(db *gorm.DB, User *User, id string) (err error) {
	db.Where("id = ?", id).Delete(User)
	return nil
}
