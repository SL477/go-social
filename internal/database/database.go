package database

import (
	"encoding/json"
	"io/ioutil"
	"time"
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

// User -
type User struct {
	CreatedAt time.Time `json:"createdAt"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Name      string    `json:"name"`
	Age       int       `json:"age"`
}

// Post -
type Post struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UserEmail string    `json:"userEmail"`
	Text      string    `json:"text"`
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