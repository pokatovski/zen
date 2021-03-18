package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

const zenApi = "https://zen.yandex.ru/api/v3/launcher/"

type zenItem struct {
	Title        string `json:"title"`
	Image        string `json:"image"`
	Link         string `json:"link"`
	CreationTime string `json:"creation_time"`
}

type zenChannel struct {
	Items []zenItem `json:"items"`
}

type channelData struct {
	Items []zenItem
}

var templates = template.Must(template.ParseGlob("templates/*"))

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/channel", channel)
	http.HandleFunc("/detail", detail)

	http.ListenAndServe(":8080", nil)
}

func channel(w http.ResponseWriter, r *http.Request) {
	var isNamed bool
	var ch, url string
	channelPath := r.URL.Query().Get("path")
	channelPath = html.EscapeString(channelPath)
	if channelPath == "" {
		err := errors.New("bad request: path is required")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//zen channels has two types: named and unnamed(https://zen.yandex.ru/channel_name or https://zen.yandex.ru/id/1)
	splittedPath := strings.Split(channelPath, "/")

	if len(splittedPath) == 4 {
		//https://zen.yandex.ru/crazydoge
		isNamed = true
		ch = splittedPath[3]
	} else if len(splittedPath) == 5 {
		//https://zen.yandex.ru/id/6022b7183646e21c6322408b
		ch = splittedPath[4]
	} else {
		err := errors.New("bad path")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("ch: ", ch)

	if isNamed {
		url = fmt.Sprintf("%smore?channel_name=%s", zenApi, ch)
	} else {
		url = fmt.Sprintf("%smore?channel_id=%s", zenApi, ch)
	}

	log.Println("url: ", url)
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
		}
		var zenCh zenChannel
		err = json.Unmarshal(bodyBytes, &zenCh)
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		//log.Println("len", len(zenCh.Items))
		//chData, err := json.Marshal(zenCh)
		//if err != nil {
		//	http.Error(w, "", http.StatusInternalServerError)
		//	return
		//}

		//w.Header().Set("Content-Type", "application/json")
		//w.Header().Set("Access-Control-Allow-Origin", "*")
		//_, err = w.Write(chData)
		//if err != nil {
		//	http.Error(w, "", http.StatusInternalServerError)
		//	return
		//}

		//fp := path.Join("templates", "channel.html")
		chData := channelData{Items: zenCh.Items}
		err = templates.ExecuteTemplate(w, "channel.html", chData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		badStatusErr := "status code from zen:" + strconv.Itoa(resp.StatusCode)
		http.Error(w, badStatusErr, http.StatusInternalServerError)
		return
	}
}

func detail(w http.ResponseWriter, r *http.Request) {

	//page := r.URL.Query().Get("page")
	//page = html.EscapeString(page)
	//boston
	page := "udivitelnye-bostonterery-5fcf95870a45a91cf4d3ed10"
	//gifs
	//page := "5-gifok-s-chihuhua-kotorye-dokazyvaiut-chto-oni-samye-smeshnye-sobaki-5ad60b0e00b3ddffc586cc31"
	//10 photos
	//page := "10-foto-dokazyvaiuscih-chto-ovcharki-luchshie-sobaki-na-zemle-5ab76aab00b3ddf3c334a7a5"

	if page != "" {
		log.Println("page: ", page)

		url := fmt.Sprintf("https://zen.yandex.ru/media/crazydoge/%s", page)
		c := colly.NewCollector()
		// On every a element which has href attribute call callback
		c.OnHTML(".article-render__block", func(e *colly.HTMLElement) {
			// Print link
			fmt.Printf("some found name: %q \n", e.Name)
			fmt.Printf("some found text: %q \n", e.Text)

		})
		c.OnHTML(".article-image__image", func(e *colly.HTMLElement) {
			// Print link
			fmt.Printf("some found name: %q \n", e.Name)
			fmt.Printf("found img src: %q \n", e.Attr("src"))
			fmt.Printf("found img data - src: %q \n", e.Attr("data-src"))
		})

		// Before making a request print "Visiting ..."
		c.OnRequest(func(r *colly.Request) {
			fmt.Println("Visiting", url)
		})

		// Start scraping on https://hackerspaces.org
		c.Visit(url)
		//resp, err := http.Get(url)
		//if err != nil {
		//	log.Println("bad request")
		//}
		//defer resp.Body.Close()
		//if resp.StatusCode == http.StatusOK {
		//	bodyBytes, err := ioutil.ReadAll(resp.Body)
		//	if err != nil {
		//		log.Fatal(err)
		//	}
		//	var zenCh zenChannel
		//	err = json.Unmarshal(bodyBytes, &zenCh)
		//	if err != nil {
		//		panic(err)
		//	}
		//	for _, i := range zenCh.Items {
		//		fmt.Println("Title", i.Title)
		//		fmt.Println("Link", i.Link)
		//		fmt.Println("CreationTime", i.CreationTime)
		//	}
		//	log.Println("len", len(zenCh.Items))
		//	chData, err := json.Marshal(zenCh)
		//	if err != nil {
		//		http.Error(w, err.Error(), http.StatusInternalServerError)
		//		return
		//	}
		//
		//	w.Header().Set("Content-Type", "application/json")
		//	w.Header().Set("Access-Control-Allow-Origin", "*")
		//	w.Write(chData)

	} else {
		log.Println("empty ch: ")
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	//templates := []string{"templates/index.html", "templates/header.html", "templates/footer.html"}
	err := templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
