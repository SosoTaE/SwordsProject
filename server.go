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


// func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
//     return func(w http.ResponseWriter, r *http.Request) {
//         // set the Access-Control-Allow-Origin header to allow requests from any origin
//         w.Header().Set("Access-Control-Allow-Origin", "*")

//         // set the Access-Control-Allow-Methods header to allow GET and POST requests
//         w.Header().Set("Access-Control-Allow-Methods", "GET, POST")

//         // set the Access-Control-Allow-Headers header to allow Content-Type and Authorization headers
//         w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

//         // if the request method is OPTIONS, return immediately with a 200 status code
//         if r.Method == "OPTIONS" {
//             w.WriteHeader(http.StatusOK)
//             return
//         }

//         // call the next handler function
//         next(w, r)
//     }
// }





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
	splited_word := strings.Split(word,"")
	for i := 0;i < len(arr);i++ {
		var score float32 = 0;
		each := strings.Split(arr[i],"")
		for j := 0;j < len(splited_word);j++ {
			var index int = findIndex(each,splited_word[j])
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
			
		}

	}

	return results
}

// func handleStaticFiles(w http.ResponseWriter, r *http.Request) {
//     http.ServeFile(w, r, r.URL.Path[1:])
// }

func main() {
	var arr []string
	
	uri := "mongodb://localhost:27017"
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
		arr = append(arr,word)
	}
	if err := cursor.Err(); err != nil {
		panic(err)
	}

	// fmt.Println(search(arr,"გამარჯობა"))


	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("GET request")
		word := r.URL.Query().Get("word")
		data := map[string]interface{}{
			"words": search(arr,word),
		}

		jsonData,err := json.Marshal(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(jsonData)

    }


    http.HandleFunc("/api/search", handlerFunc)

    fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

    http.ListenAndServe(":8000", nil)

  }