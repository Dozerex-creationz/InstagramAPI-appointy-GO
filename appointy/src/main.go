package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MyMux struct {
}
type user struct {
	UserID   int    `json:"UserID"`
	UserName string `json:"UserName"`
	Password string `json:"Password"`
	Email    string `json:"Email"`
}
type post struct {
	PostID    int    `json:"PostID"`
	Caption   string `json:"Caption"`
	ImageURL  string `json:"ImageURL"`
	Timestamp string `json:"Timestamp"`
}

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func (p *MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/userRegister" {
		if r.Method == "GET" {
			userDataSend(w, r)
			return
		}
	} else if r.URL.Path == "/users" {
		if r.Method == "POST" {
			userCall(w, r)
			return
		}
	} else if r.URL.Path == "/postRegister" {
		if r.Method == "GET" {
			postDataSend(w, r)
			return
		}
	} else if r.URL.Path == "/posts" {
		if r.Method == "POST" {
			postCall(w, r)
			return
		}
	} else if r.URL.Path == "/users/" {
		// UserId := r.URL.Path[7:]
		// clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

		// client, err := mongo.Connect(context.TODO(), clientOptions)
		// if err != nil {
		// 	fmt.Println(err)
		// 	os.Exit(1)
		// }

		// col := client.Database("InstaLitre").Collection("UserDataBase")
		// result, insertErr := col.FindOne(context.Background(), bson.M{"objectID": UserId})
		// fmt.Println(result.Decode(user))
	} else if r.URL.Path == "/posts/" {
		fmt.Println("Yet To Complete")
	} else if r.URL.Path == "/posts/users/" {
		fmt.Println("Yet To Complete")
	}

	http.NotFound(w, r)
}

func userCall(writer http.ResponseWriter, reader *http.Request) {
	reader.ParseForm()
	fmt.Println(reader.Form)
	fmt.Println(reader.Form["UserID"])
	var bUser interface{}
	TempUser := new(user)
	TempUser.UserID = 10
	TempUser.UserName = reader.FormValue("username")
	TempUser.Email = reader.FormValue("email")
	TempUser.Password = reader.FormValue("password")

	jsonUser, err := json.Marshal(TempUser)
	if err != nil {
		http.Error(writer, err.Error(), 500)
		return
	}
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

	col := client.Database("InstaLitre").Collection("UserDataBase")
	bson.UnmarshalExtJSON(jsonUser, true, &bUser)
	result, insertErr := col.InsertOne(ctx, &bUser)
	if insertErr != nil {
		fmt.Println(insertErr)
		os.Exit(1)
	}
	fmt.Println(jsonUser)
	fmt.Println("ID :", result.InsertedID)
	fmt.Fprintf(writer, "Hello to the client!")
}

func userDataSend(writer http.ResponseWriter, reader *http.Request) {
	fmt.Println("method: ", reader.Method)
	tpl.ExecuteTemplate(writer, "userDataSend.gohtml", nil)
}

func postCall(writer http.ResponseWriter, reader *http.Request) {
	reader.ParseForm()
	fmt.Println(reader.Form)
	TempPost := new(post)
	TempPost.PostID = 20
	TempPost.Caption = reader.FormValue("caption")
	TempPost.ImageURL = reader.FormValue("imageURL")
	TempPost.Timestamp = reader.FormValue("timestamp")
	var bPost interface{}

	jsonPost, err := json.Marshal(TempPost)
	if err != nil {
		http.Error(writer, err.Error(), 500)
		return
	}
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

	col := client.Database("InstaLitre").Collection("PostDataBase")
	bson.UnmarshalExtJSON(jsonPost, true, &bPost)
	result, insertErr := col.InsertOne(ctx, &bPost)
	if insertErr != nil {
		fmt.Println(insertErr)
		os.Exit(1)
	}
	fmt.Println(jsonPost)
	fmt.Println("ID :", result.InsertedID)
	fmt.Fprintf(writer, "Hello to the client!")
}

func postDataSend(writer http.ResponseWriter, reader *http.Request) {
	fmt.Println("method: ", reader.Method)
	tpl.ExecuteTemplate(writer, "postDataSend.gohtml", nil)
}

func main() {
	fmt.Println("started...in PORT no:9090....")
	mux := &MyMux{}
	http.ListenAndServe(":9090", mux)

}
