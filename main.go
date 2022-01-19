package main

import (
	"net/http"
	"time"
	"encoding/json"
	"github.com/SL477/go-social/internal/database"
	"errors"
)

type errorBody struct {
	Error string `json:"error"`
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	// Add headers
	w.Header().Set("Content-Type", "application/json")

	// Write JSON body
	response, _ := json.Marshal(payload)
	// deal with err ...
	w.Write(response)

	// Write status code
	w.WriteHeader(code)
}

func respondWithError(w http.ResponseWriter, err error) {
	e := errorBody{
		Error: err.Error(),
	}
	respondWithJSON(w, http.StatusBadRequest, e)
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	/*w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write([]byte("{}"))*/

	// you can use any compatible type, but let's use our database package's User type for practice
	respondWithJSON(w, 200, database.User{
		Email: "test@example.com",
	})
}

func testErrHandler(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, errors.New("test error"))
}

type apiConfig struct {
	dbClient database.Client
}

/*func (apiCfg apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {

}*/
type parameters struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
}

func (apiCfg apiConfig) endpointUsersHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, err)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// Get handler
		u,err := apiCfg.dbClient.GetUser(params.Email)
		if err != nil {
			respondWithError(w, err)
			return
		}
		respondWithJSON(w, 200, u)
		break
	case http.MethodPost:
		// Post handler
		u, err := apiCfg.dbClient.CreateUser(params.Email, params.Password, params.Name, params.Age)
		if err != nil {
			respondWithError(w, err)
			return
		}
		respondWithJSON(w, 201, u)
		break
	case http.MethodPut:
		// Put handler
		u,err := apiCfg.dbClient.UpdateUser(params.Email, params.Password, params.Name, params.Age)
		if err != nil {
			respondWithError(w, err)
			return
		}
		respondWithJSON(w, http.StatusOK, u)
		break
	case http.MethodDelete:
		// Delete handler
		if params.Email == "" {
			respondWithJSON(w, http.StatusBadRequest, errorBody{
				Error: "User account required",
			})
			return
		}
		err := apiCfg.dbClient.DeleteUser(params.Email)
		if err != nil {
			respondWithError(w, err)
			return
		}
		/*respondWithJSON(w, 201, errorBody{
			Error: "Deleted user account",
		})*/
		//apiCfg.handlerDeleteUser(w, r)
		respondWithJSON(w, http.StatusOK, struct{}{})
		break
	default:
		respondWithError(w, errors.New("method not supported"))
	}
}



func main() {
	// Setup database
	apiCfg := apiConfig{
		dbClient: database.NewClient("db.json"),
	}

	// Run server
	m := http.NewServeMux()
	m.HandleFunc("/", testHandler)
	m.HandleFunc("/err", testErrHandler)
	m.HandleFunc("/users", apiCfg.endpointUsersHandler)
	m.HandleFunc("/users/", apiCfg.endpointUsersHandler)

	const addr = "localhost:8080"
	srv := http.Server{
		Handler: m,
		Addr: addr,
		WriteTimeout: 30 * time.Second,
		ReadTimeout: 30 * time.Second,
	}
	srv.ListenAndServe()
}