package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"
)

type Pattern struct {
	Name string `json:"name"`
	Age  string `json:"age"`
}

type Config struct {
	Patterns []Pattern `json:"patterns"`
	Index    string    `json:"index"`
}

type CatEntry struct {
	Index      string `json:"index"`
	StoreSize  string `json:"store.size"`
	DateString string `json:"date"`
	IndexBase  string `json:"indexbase"`
}

func LoadConfig(f string) Config {
	file, err := os.Open(f)
	if err != nil {
		log.Fatal("can't open config file: ", err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	MyConfig := Config{}
	err = decoder.Decode(&MyConfig)
	if err != nil {
		log.Fatal("can't decode config JSON: ", err)
	}
	return MyConfig
}

func QueryCat(index string) []CatEntry {
	url := "https://" + index + "/_cat/indices?format=json&h=index,store.size"
	catClient := http.Client{
		Timeout: time.Second * 5, // Maximum of 2 secs
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	res, getErr := catClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	MyCats := []CatEntry{}
	jsonErr := json.Unmarshal(body, &MyCats)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	fmt.Println(MyCats)
	r := regexp.MustCompile(`(?P<Indexbase>.*)-(?P<Date>\d{4}.\d{2}.\d{2})$`)
	for _, elem := range MyCats {
		res := r.FindStringSubmatch(elem.Index)
		names := r.SubexpNames()
		for i, _ := range res {
			if i != 0 {
				fmt.Println(names[i], res[i])
				if names[i] == "Indexbase" {
					elem.IndexBase = res[i]
				} else if names[i] == "Date" {
					elem.DateString = res[i]
				}
			}
		}
		//fmt.Println(elem)
		fmt.Printf("%+v\n", elem)
	}
	return []CatEntry{}
}

func main() {
	c := flag.String("c", "./config.json", "Specify the configuration file.")
	flag.Parse()
	MyConfig := LoadConfig(*c)
	QueryCat(MyConfig.Index)
	fmt.Printf("%+v\n", MyConfig)
}
