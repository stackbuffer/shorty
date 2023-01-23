package main

import (
	"fmt"
	"net/http"
	"crypto/sha1"
	"encoding/json"
	"github.com/go-redis/redis"
)

type Response struct {
	Url string `json:"url"`
}

func connectToRedis() *redis.Client {

	//notice the "s" after redis, I don't know why this fixed the connection issue :|
	opt, _ := redis.ParseURL("rediss://default:514168316552439f89740797dfc30d33@eu2-firm-cow-31460.upstash.io:31460")
	client := redis.NewClient(opt)

	return client
}


func shorty(w http.ResponseWriter, r *http.Request){
	//CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)

	mp := r.URL.Query()

	urlToShorten := mp["url"][0]

	hash := sha1.New()
	hash.Write([]byte(urlToShorten))
	bs := hash.Sum(nil)

	finalHash := fmt.Sprintf("%x\n",bs)

	finalHash_6 := finalHash[:6]

	//map the hash to the url in the redis
	client := connectToRedis()

	client.Set(finalHash_6, urlToShorten, 0)

	resp := Response{Url : fmt.Sprintf("localhost:8080/%v", finalHash_6)}

	json.NewEncoder(w).Encode(resp)
	
	//fmt.Fprintf(w, resp)
}


func goToShortyURL(w http.ResponseWriter, r *http.Request){

	key := (r.URL.Path)[1:]

	client := connectToRedis()

	url, _ := client.Get(key).Result()


	if url != "" {
		http.Redirect(w, r, url, http.StatusSeeOther)
	} else {
		url = "Not found"
			
		//CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprintf(w, url)
	}

}



func main(){
	http.HandleFunc("/", goToShortyURL)
	http.HandleFunc("/shorty", shorty)

	http.ListenAndServe(":8080", nil)
	fmt.Println("Listening on port 8080")
}