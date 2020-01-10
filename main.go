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
	"strconv"
	"time"
)

const (
	layout = "2006.01.02"
)

type Pattern struct {
	Name string `json:"name"`
	Age  string `json:"age"`
}

type Config struct {
	Patterns []Pattern `json:"patterns"`
	Domain   string    `json:"domain"`
	Username string    `json:"username"`
	Password string    `json:"password"`
}

type CatEntry struct {
	Index      string `json:"index"`
	StoreSize  string `json:"store.size"`
	DateString string
	IndexBase  string
	Age        int
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

func QueryCat(domain string, username string, password string) []CatEntry {
	url := "https://" + domain + "/_cat/indices/*,-%2E*?format=json&h=index,store.size"
	catClient := http.Client{
		Timeout: time.Second * 5, // Maximum of 2 secs
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	if username != "" {
		req.SetBasicAuth(username, password)
	}
	res, getErr := catClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}
	if res.StatusCode != 200 {
		log.Fatal(fmt.Sprintf("HTTP response code %d from Elasticsearch\n", res.StatusCode))
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

	now := time.Now()
	r := regexp.MustCompile(`(?P<Indexbase>.*)-(?P<Date>\d{4}.\d{2}.\d{2})$`)
	for idx, elem := range MyCats {
		res := r.FindStringSubmatch(elem.Index)
		names := r.SubexpNames()
		for i, _ := range res {
			if i != 0 {
				if names[i] == "Indexbase" {
					MyCats[idx].IndexBase = res[i]
				} else if names[i] == "Date" {
					MyCats[idx].DateString = res[i]
					t, err := time.Parse(layout, res[i])
					if err != nil {
						fmt.Println(err)
					}
					MyCats[idx].Age = int(now.Sub(t).Hours()) / 24
				}
			}
		}
	}
	return MyCats
}

func GetPurgeIndexes(c []CatEntry, config Config, del bool) []string {
	ret := []string{}
	if del == false {
		fmt.Println("This is a dry-run")
	}
	// Build map from config
	m := make(map[string]int)
	for _, v := range config.Patterns {
		m[v.Name], _ = strconv.Atoi(v.Age)
	}
	for _, v := range c {
		if _, ok := m[v.IndexBase]; ok {
			fmt.Printf("Found config entry for %s - ", v.Index)
			if m[v.IndexBase] < v.Age {
				fmt.Printf("Index %s is older than %d days and will be purged\n", v.Index, m[v.IndexBase])
				ret = append(ret, v.Index)
			} else {
				fmt.Printf("Index age is within limits\n")
			}
		}
	}
	return ret
}

// Purge takes a  list of indexes to DELETE from the ES domain and deletes them
func Purge(l []string, index string, username string, password string) {
	for _, val := range l {
		url := `https://` + index + "/" + val
		fmt.Printf("DELETE %s\n", url)
		delClient := http.Client{
			Timeout: time.Second * 5, // Maximum of 2 secs
		}
		req, err := http.NewRequest(http.MethodDelete, url, nil)
		if err != nil {
			log.Fatal(err)
		}
		if username != "" {
			req.SetBasicAuth(username, password)
		}
		res, getErr := delClient.Do(req)
		if getErr != nil {
			log.Fatal(getErr)
		}
		if res.StatusCode != 200 {
			log.Fatal(fmt.Sprintf("HTTP response code %d from Elasticsearch\n", res.StatusCode))
		}
		body, readErr := ioutil.ReadAll(res.Body)
		if readErr != nil {
			log.Fatal(readErr)
		}
		fmt.Println(string(body))
	}
}

func main() {
	c := flag.String("c", "./config.json", "Specify the configuration file.")
	del := flag.Bool("d", false, "Full delete run - not dry-run.")
	flag.Parse()

	MyConfig := LoadConfig(*c)
	Indexes := QueryCat(MyConfig.Domain, MyConfig.Username, MyConfig.Password)
	purges := GetPurgeIndexes(Indexes, MyConfig, *del)
	if *del == true {
		Purge(purges, MyConfig.Domain, MyConfig.Username, MyConfig.Password)
	}
}
