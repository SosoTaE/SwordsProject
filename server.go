package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


func findIndex(arr []string, val string) int {
	for i, v := range arr {
		if v == val {
			return i
		}
	}
	return -1
}

func search(arr []string, word string) []string {
	var results []string
	splited_word := strings.Split(word, "")
	for i := 0; i < len(arr); i++ {
		var score float32 = 0
		each := strings.Split(arr[i], "")
		for j := 0; j < len(splited_word); j++ {
			var index int = findIndex(each, splited_word[j])
			if index != -1 {
				score += 1
				each[index] = ""
			} else {
				score += 0
			}
		}

		// fmt.Println(score,len(splited_word),arr[i])

		if len(each) > len(splited_word) {
			score = score / float32(len(each))
		} else {
			score = score / float32(len(splited_word))

		}

		// fmt.Println(score)

		if score >= 0.8 {
			results = append(results, arr[i])

		} else {
			length := len(arr[i])
			Wlength := len(word)
			if length >= 3 * 4 && Wlength >=  3 * 4 {
				if arr[i][length - 3 * 4:length] == word[Wlength - 3 * 4:Wlength] {
					results = append(results, arr[i])
				} 
				// else if arr[i][length - 3 * 3:length] == word[Wlength - 3 * 3:Wlength] {
				// 	results = append(results, arr[i])
				// }
			}
		}

	}

	return results
}

func main() {
	fmt.Println("server is started")
	var arr []string

	uri := "mongodb+srv://sosotae:SosoTaENaoko;@sword.yah6gsy.mongodb.net/?retryWrites=true&w=majority"
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	db := client.Database("texts")
	coll := db.Collection("words")

	cursor, err := coll.Find(context.TODO(), bson.D{})
	if err != nil {
		panic(err)
	}

	// iterate code goes here
	for cursor.Next(context.TODO()) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			panic(err)
		}
		word := result["word"].(string)
		arr = append(arr, word)
	}
	if err := cursor.Err(); err != nil {
		panic(err)
	}

	// fmt.Println(search(arr,"გამარჯობა"))

	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		// fmt.Println("GET request")
		word := r.URL.Query().Get("word")
		data := map[string]interface{}{
			"words": search(arr, word),
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(jsonData)

	}

	http.HandleFunc("/api/search", handlerFunc)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	http.ListenAndServe(":8080", nil)

}
