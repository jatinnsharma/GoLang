package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/jatinnsharma/GoLang/models"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

type UserController struct {
	session *mgo.Session
}

func NewUserController(s *mgo.Session) *UserController {
	return &UserController{s}
}

func (uc UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	// check id is hex else give 404
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound)
	}

	//oid is used for finding data in mongodb
	oid := bson.ObjectIdHex(id)

	// u:= user model for finding user data
	u := models.User{}

	// found the data in db
	if err := uc.session.DB("mongo-golang").C("users").FindId(oid).One(&u); err != nil {
		w.WriteHeader(404)
		return
	}

	// Marshal : convert to json , json to simple data
	// uj store in uj
	uj, err := json.Marshal(u)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s\n", uj)
}

func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	u := models.User{}

	// decoder body from json to
	json.NewDecoder(r.Body).Decode(&u)

	// create new object id
	u.Id = bson.NewObjectId()
	uc.session.DB("mongo-golang").C("users").Insert(u)

	// Marshal data into JSON
	uj,err := json.Marshal(u)

	if err != nil {
		fmt.Println(err) 
	}

	// send back to frontend
	w.Header().Set("Content-Type" , "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Println(w,"%s\n" , uj) 

}

func (uc UserController) DeleteUser(w http.ResponseWriter , r *http.Request, p httprouter.Params){
	// get id from params
	id:= p.ByName("id")
	
	// check id is hex else 404
	if !bson.IsObjectIdHex(id){
		w.WriteHeader(404)
		return 
	}
	
	// oid is used for finding data in mongodb
	oid:= bson.ObjectIdHex(id)

	if err :=uc.session.DB("mongo-golang").C("users").RemoveId(oid); err!=nil{
		w.WriteHeader(404)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w,"Deleted user " ,oid , "\n")
}
