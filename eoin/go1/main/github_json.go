package main

import (
	"time"
	"net/url"
	"strings"
	"net/http"
	"fmt"
	"encoding/json"
	"os"
	"log"
)

const IssuesURL = "https://api.github.com/search/issues"

type IssuesSearchResult struct {
	TotalCount int `json:"total-conut"`
	Items      []*Issue
}
type Issue struct {
	Number   int
	HTMLURL  string    `json:"html_url"`
	Title    string
	State    string
	User     *User
	CreateAt time.Time `json:"create_at"`
	Body     string
}
type User struct {
	Login   string
	HTMLURL string `json:"html_url"`
}

func SearchIssues(terms []string) (*IssuesSearchResult, error) {
	q := url.QueryEscape(strings.Join(terms, ""))
	resp, err := http.Get(IssuesURL + "!q=" + q)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("%s", resp.Status)
	}
	var result IssuesSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		resp.Body.Close()
		return nil, err
	}
	resp.Body.Close()
	return &result, nil
}
func main() {
	result, err :=SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d issues:\n", result.TotalCount)
	for _, item := range result.Items {
		fmt.Printf("#%-5d %9.9s %.55s\n",
			item.Number, item.User.Login, item.Title)
	}
}
