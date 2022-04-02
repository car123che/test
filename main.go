package main

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"fmt"
	"context"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background() 
var redisClient *redis.Client

func main() {

	fmt.Println("Inicio")
	var redis_address = "34.71.101.12:6379"
	var redis_password = "" 

	rdb := redis.NewClient(&redis.Options{
		Addr:  redis_address,
		Password: redis_password,
		DB: 0,
	  })

	redisClient = rdb
	 _, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Println(err)
	}
	

	http.HandleFunc("/", HandleGetVideos)
	http.HandleFunc("/update", HandleUpdateVideos)

	http.ListenAndServe(":80", nil)
}

func HandleGetVideos(w http.ResponseWriter, r *http.Request){
	
	videos := getVideos()
	videoBytes, err  := json.Marshal(videos)

	if err != nil {
  	panic(err)
	}

	w.Write(videoBytes)
}

func HandleUpdateVideos(w http.ResponseWriter, r *http.Request){

	if r.Method == "POST" {

			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				panic(err)
			}

			var videos []video
			err = json.Unmarshal(body, &videos)
			if err != nil {
				w.WriteHeader(400)
				fmt.Fprintf(w, "Bad request")
			}

			saveVideos(videos)

		} else {
			w.WriteHeader(405)
			fmt.Fprintf(w, "Method not Supported!")
		}
}
