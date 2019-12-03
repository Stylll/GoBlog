package models

import (
	"errors"
	"fmt"
	"html"
	"log"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int       `gorm:"primary_key;auto_increment" json:"id"`
	Firstname string    `gorm:"size:255;not null" json:"firstname"`
	Lastname  string    `gorm:"size:255;not null" json:"lastname"`
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	Password  string    `gorm:"size:100;not null" json:"password"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
}

func (u *User) BeforeSave() error {
	hashedPassword, err := Hash(u.Password)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)
	return nil
}

func (u *User) Prepare() {
	u.ID = 0
	u.Firstname = html.EscapeString(strings.TrimSpace(u.Firstname))
	u.Lastname = html.EscapeString(strings.TrimSpace(u.Lastname))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
}

func (u *User) Validate(operation string) error {
	switch strings.ToLower(operation) {
	case "update":
		if u.Firstname == "" {
			return errors.New("Firstname Required")
		}
		if u.Lastname == "" {
			return errors.New("Lastname Required")
		}
		if u.Email == "" {
			return errors.New("Email Required")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Email Invalid")
		}
		if u.Password == "" {
			return errors.New("Password Required")
		}

		return nil

	case "login":
		if u.Email == "" {
			return errors.New("Email Required")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			fmt.Println("Email invalid")
			return errors.New("Email Invalid")
		}
		if u.Password == "" {
			return errors.New("Password Required")
		}

		return nil

	default:
		if u.Firstname == "" {
			return errors.New("Firstname Required")
		}
		if u.Lastname == "" {
			return errors.New("Lastname Required")
		}
		if u.Email == "" {
			return errors.New("Email Required")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Email Invalid")
		}
		if u.Password == "" {
			return errors.New("Password Required")
		}

		return nil
	}
}

func (u *User) SaveUser(db *gorm.DB) (*User, error) {
	var err error
	err = db.Debug().Create(&u).Error
	if err != nil {
		return &User{}, err
	}

	return u, nil
}

func (u *User) FindAllUsers(db *gorm.DB) (*[]User, error) {
	var err error
	users := []User{}
	err = db.Debug().Model(&User{}).Limit(100).Find(&users).Error

	if err != nil {
		return &[]User{}, err
	}

	return &users, nil
}

func (u *User) FindUserByID(db *gorm.DB, uid uint64) (*User, error) {
	var err error
	err = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&u).Error

	if gorm.IsRecordNotFoundError(err) {
		return &User{}, errors.New("User Not Found")
	}

	if err != nil {
		return &User{}, err
	}

	return u, nil
}

func (u *User) UpdateAUser(db *gorm.DB, uid int64) (*User, error) {
	var err error

	// hash password
	err = u.BeforeSave()
	if err != nil {
		log.Fatal(err)
		return u, err
	}

	// update the record
	err = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&u).UpdateColumns(
		map[string]interface{}{
			"firstname": u.Firstname,
			"lastname":  u.Lastname,
			"email":     u.Email,
			"password":  u.Password,
			"updatedAt": time.Now(),
		},
	).Error

	if err != nil {
		return &User{}, err
	}

	// retrieve the updated record
	err = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&u).Error

	if err != nil {
		return &User{}, err
	}

	return u, nil
}

func (u *User) DeleteAUser(db *gorm.DB, uid int64) (int64, error) {

	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&u).Delete(&u)
	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil
}
