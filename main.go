package main

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
)

type User struct {
	ID    int
	Token string
	Name  string
	Age   int
}

func create_users() map[string]User {
	var users map[string]User = make(map[string]User)
	users["1"] = User{1, "123", "Сильвестр В Столовой", 48}
	users["2"] = User{2, "234", "Иосиф В Прихожей", 45}
	users["3"] = User{3, "345", "Роберт В Сауне Младший", 12}
	users["4"] = User{4, "4qwerty", "Стивен Сигнал", 32}
	users["5"] = User{5, "567", "Джек Воробей", 89}
	users["6"] = User{6, "123", "Муртаза Рахимов", 78}
	return users
}

var userdb map[string]User = create_users()

func UserHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Printf("%+v\n", r)
	if r.Method == http.MethodGet {
		UserGetHandler(w, r)
	} else if r.Method == http.MethodPost {
		UserPostHandler(w, r)
	} else {
		// error
	}
}

func UserGetHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[6:]
	js, _ := json.Marshal(userdb[id])
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	//fmt.Fprintf(w, "Имя: %s, Возраст: %d, Ид: %d, Токен: %s", userdb[id].Name, userdb[id].Age, userdb[id].ID, userdb[id].Token)
}

func UserPostHandler(w http.ResponseWriter, r *http.Request) {

	if !IsContentJson(r) {
		http.Error(w, "Content-Type is not application/json", http.StatusUnsupportedMediaType)
		return
	}

	RequestUser, err := DecodeJson(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	EditUser(r.URL.Path[6:], RequestUser)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	jsonResp, _ := json.Marshal(userdb[r.URL.Path[6:]])
	w.Write(jsonResp)
}

func IsContentJson(r *http.Request) bool {
	contentType := r.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return false
	}
	if mediatype != "application/json" {
		return false
	}
	return true
}

func DecodeJson(r *http.Request) (User, error) {
	var err error = nil
	var ReqUser User
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&ReqUser)
	// fmt.Printf("json decoded: %+v\n", RequestUser)
	return ReqUser, err
}

func EditUser(srcId string, dstUser User) {
	fmt.Printf("%+v\n", dstUser)
	User_temp := userdb[srcId]
	if len(dstUser.Name) > 0 {
		User_temp.Name = dstUser.Name
	}
	if dstUser.Age > 0 {
		User_temp.Age = dstUser.Age
	}
	userdb[srcId] = User_temp
}

func main() {
	http.HandleFunc("/user/", UserHandler)
	fmt.Println("starting")
	http.ListenAndServe(":80", nil)
}
