package database

import (
	"encoding/json"
	"io/ioutil"
	"time"
	"errors"
)

type Client struct {
	filepath string
}

func NewClient(path string) Client {
	c := Client{path}
	return c
}

type databaseSchema struct {
	Users map[string]User `json:"users"`
	Posts map[string]Post `json:"posts"`
}

// Create the database
func (c Client) createDB() error {
	data, err := json.Marshal(databaseSchema{
		Users: make(map[string]User),
		Posts: make(map[string]Post),
	})
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(c.filepath, data, 0600)
	return err
}

func (c Client) EnsureDB() error {
	_, err := ioutil.ReadFile(c.filepath)
	/*if errors.Is(err, os.ErrNotExist) {
		return c.createDB()
	}*/
	if err != nil {
		return c.createDB()
	}
	return err
}

// Database helper functions
func (c Client) updateDB(db databaseSchema) error {
	data, err := json.Marshal(db)
	if err == nil {
		return ioutil.WriteFile(c.filepath, data, 0600)
	}
	return err
}

func (c Client) readDB() (databaseSchema, error) {
	data, err := ioutil.ReadFile(c.filepath)
	// turn into database
	db := databaseSchema{
		Users: make(map[string]User),
		Posts: make(map[string]Post),
	}
	if err == nil {
		err = json.Unmarshal(data, &db)
		return db, err
	}
	return db, err
}

// Create user
func (c Client) CreateUser(email, password, name string, age int) (User, error) {
	u := User{
		CreatedAt: time.Now().UTC(),
		Email: email,
		Password: password,
		Name: name,
		Age: age,
	}

	// Get the current database
	db, err := c.readDB()
	if err != nil {
		return u, err
	}
	db.Users[email] = u
	
	err = c.updateDB(db)

	return u, err
}

// Update user
func (c Client) UpdateUser(email, password, name string, age int) (User, error) {
	// Get the current database
	db, err := c.readDB()

	u, exists := db.Users[email]

	if !exists {
		return User{}, errors.New("user doesn't exist")
	}
	u.Email = email
	u.Password = password
	u.Name = name
	u.Age = age

	err = c.updateDB(db)

	return u, err
}

// Get user
func (c Client) GetUser(email string) (User, error) {
	db, err := c.readDB()
	u, exists := db.Users[email]
	if !exists {
		return User{}, errors.New("user doesn't exist")
	}
	return u, err
}

// Delete user
func (c Client) DeleteUser(email string) error {
	db, err := c.readDB()
	_, exists := db.Users[email]
	if exists {
		delete(db.Users, email)
		err = c.updateDB(db)
	}
	return err
}