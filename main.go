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
	var redis_address = "localhost:6379"
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
	
	videos := GetVideos()
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

			SaveVideos(videos)

		} else {
			w.WriteHeader(405)
			fmt.Fprintf(w, "Method not Supported!")
		}
}

type video struct{
	Id string
	Title string
	Description string
	Imageurl string
	Url string
}

func GetVideos()(videos []video){
	
	keys, err := redisClient.Keys(ctx,"*").Result()

	if err != nil {
		panic(err)
	}

	for _, key := range keys {
		video := GetVideo(key)
		videos = append(videos, video)
	}

	return videos
}


	

func GetVideo(id string)(video video) {
	
	value, err := redisClient.Get(ctx, id).Result()

	if err != nil {
		panic(err)
	}

	if err != redis.Nil {
		err = json.Unmarshal([]byte(value), &video)
	}
	
	return video
}


func SaveVideo(video video)(){

	videoBytes, err  := json.Marshal(video)
	if err != nil {
		  panic(err)
	  }
  
	err = redisClient.Set(ctx, video.Id, videoBytes, 0).Err()
	if err != nil {
		  panic(err)
	  }
  }


func SaveVideos(videos []video)(){
	for _, video := range videos {
		SaveVideo(video)
	}

}
