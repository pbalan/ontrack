package repositories

import (
	"github.com/pbalan/ontrack/src/graph/model"
	"gorm.io/gorm"
)

//create a user
func CreateUser(db *gorm.DB, NewUser *model.User) (User *model.User, err error) {
	err = db.Create(&NewUser).Error

	if err != nil {
		return NewUser, err
	}
	return NewUser, nil
}

//get users
func GetUsers(db *gorm.DB, User *[]model.User) (err error) {
	err = db.Find(User).Error
	if err != nil {
		return err
	}
	return nil
}

//get user by email
func UserGetByEmail(db *gorm.DB, email string) (*model.User, error) {
	user := model.User{}
	err := db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

//get user by username
func UserGetByUsername(db *gorm.DB, username string) (*model.User, error) {
	user := model.User{}
	err := db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

//get user by id
func GetUser(db *gorm.DB, User *model.User, id string) (err error) {
	err = db.Where("id = ?", id).First(User).Error
	if err != nil {
		return err
	}
	return nil
}

//update user
func UpdateUser(db *gorm.DB, User *model.User) (err error) {
	db.Save(User)
	return nil
}

//delete user
func DeleteUser(db *gorm.DB, User *model.User, id string) (err error) {
	db.Where("id = ?", id).Delete(User)
	return nil
}
