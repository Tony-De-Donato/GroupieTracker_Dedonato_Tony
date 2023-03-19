package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Booklist struct {
	Docs []struct {
		Key string `json:"key"`
	}
}

type Onebook struct {
	Title   string `json:"title"`
	Covers  []int  `json:"covers"`
	Authors []struct {
		Author struct {
			Key string `json:"key"`
		} `json:"author"`
	} `json:"authors"`
	Key         string `json:"key"`
	Description string `json:"description"`
}
type Onebook2 struct {
	Title   string `json:"title"`
	Covers  []int  `json:"covers"`
	Authors []struct {
		Author struct {
			Key string `json:"key"`
		} `json:"author"`
	} `json:"authors"`
	Key         string `json:"key"`
	Description struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"description"`
}

func main() {
	get_booklist_info(search_book("Le petit prince"))

}

func search_book(text string) Booklist {

	text = make_request_usable(text)
	url := "https://openlibrary.org/search.json?title=" + text

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Erreur de demande : ", err)

	}

	defer resp.Body.Close()

	var books Booklist
	err = json.NewDecoder(resp.Body).Decode(&books)
	if err != nil {
		fmt.Println("Erreur de décodage JSON 1: ", err)
	}
	return books
}

func get_booklist_info(book Booklist) []Onebook {

	list_to_display := []Onebook{}

	for _, element := range book.Docs {
		url := "https://openlibrary.org/" + element.Key + ".json"

		resp, errror := http.Get(url)
		resp_copy, errror2 := http.Get(url)
		if errror != nil || errror2 != nil {
			fmt.Println("Erreur de demande")
		}

		defer resp.Body.Close()
		defer resp_copy.Body.Close()

		var onebook Onebook
		err2 := json.NewDecoder(resp.Body).Decode(&onebook)
		if err2 != nil {
			//fmt.Println("ERREUR :", err2, "  ", onebook.Key)
			var onebook2 Onebook2
			err3 := json.NewDecoder(resp_copy.Body).Decode(&onebook2)
			if err3 == nil {

				var onebook3 Onebook
				onebook3.Description = onebook2.Description.Value
				onebook3.Title = onebook2.Title
				onebook3.Covers = onebook2.Covers
				onebook3.Authors = onebook2.Authors
				onebook3.Key = onebook2.Key

				if len(onebook3.Description) > 0 && len(onebook3.Covers) >= 1 && len(onebook3.Authors) >= 1 && onebook3.Title != "" && len(onebook3.Authors[0].Author.Key) > 0 {
					list_to_display = append(list_to_display, onebook3)
					fmt.Println(onebook3.Title)
				}
			} else {
				fmt.Println("Erreur de décodage JSON 3: ", err3)
			}
		}

		if err2 == nil && len(onebook.Description) > 0 && len(onebook.Covers) >= 1 && len(onebook.Authors) >= 1 && onebook.Title != "" && len(onebook.Authors[0].Author.Key) > 0 {

			list_to_display = append(list_to_display, onebook)
			fmt.Println(onebook.Title)
		}

	}
	return list_to_display
}

func make_request_usable(recherche string) string {
	resultat := ""

	for i := 0; i < len(recherche); i++ {
		if recherche[i] == ' ' {
			resultat += string('+')
		} else {
			resultat += string(recherche[i])
		}
	}
	return resultat

}
