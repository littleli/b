package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"net/http"
	"net/url"
)

var baseURI = "https://api.search.brave.com/res/v1/web/search"

type WebSearchApiResponse struct {
	Type   string        `json:"type"`
	Mixed  MixedResponse `json:"mixed,omitempty"`
	Query  Query         `json:"query,omitempty"`
	Videos Videos        `json:"videos,omitempty"`
	Web    Search        `json:"web,omitempty"`
}

type MixedResponse struct {
	Type string            `json:"mixed"`
	Main []ResultReference `json:"main,omitempty"`
	Top  []ResultReference `json:"top,omitempty"`
	Side []ResultReference `json:"side,omitempty"`
}

type ResultReference struct {
	Type  string `json:"type"`
	Index int    `json:"index,omitempty"`
	All   bool   `json:"all"`
}

type Query struct {
	Original          string `json:"original"`
	ShowStrictWarning bool   `json:"show_strict_warning,omitempty"`
	IsNavigational    bool   `json:"is_navigational,omitempty"`
	IsNewsBreaking    bool   `json:"is_news_breaking,omitempty"`
	IsSpellcheckOff   bool   `json:"spellcheck_off,omitempty"`
	Country           string `json:"country,omitempty"`
	IsBadResults      bool   `json:"bad_results,omitempty"`
	ShouldFallback    bool   `json:"should_fallback,omitempty"`
	PostalCode        string `json:"postal_code,omitempty"`
	City              string `json:"city,omitempty"`
	HeaderCountry     string `json:"header_country,omitempty"`
	IsThereMore       bool   `json:"more_results_available,omitempty"`
	State             string `json:"state,omitempty"`
}

type Videos struct {
	Type               string        `json:"type"`
	Results            []VideoResult `json:"results"`
	IsMutatedByGoggles bool          `json:"mutated_by_goggles,omitempty"`
}

type VideoResult struct {
	Type      string    `json:"type"`
	Video     VideoData `json:"video_data"`
	MetaURL   MetaURL   `json:"meta_url,omitempty"`
	Thumbnail Thumbnail `json:"thumbnail,omitempty"`
	Age       string    `json:"age,omitempty"`
}

type VideoData struct {
	Duration               string    `json:"duration,omitempty"`
	Views                  string    `json:"views,omitempty"`
	Creator                string    `json:"creator,omitempty"`
	Publisher              string    `json:"publisher,omitempty"`
	Thumbnail              Thumbnail `json:"thumbnail,omitempty"`
	Tags                   []string  `json:"tags,omitempty"`
	Author                 Profile   `json:"author,omitempty"`
	IsSubscriptionRequired bool      `json:"requires_subscription,omitempty"`
}

type Thumbnail struct {
	Source   string `json:"src"`
	Original string `json:"original,omitempty"`
}

type Profile struct {
	Name     string `json:"name"`
	LongName string `json:"long_name"`
	URL      string `json:"url,omitempty"`
	Image    string `json:"img,omitempty"`
}

type MetaURL struct {
	Scheme          string `json:"scheme"`
	NetworkLocation string `json:"netloc"`
	DomainName      string `json:"hostname,omitempty"`
	Favicon         string `json:"favicon"`
	Path            string `json:"path"`
}

type Search struct {
	Type             string          `json:"type"`
	Results          []SearchResults `json:"results"`
	IsFamilyFriendly bool            `json:"family_friendly"`
}

type Result struct {
	Title            string  `json:"title"`
	URL              string  `json:"url"`
	IsSourceLocal    bool    `json:"is_source_local"`
	IsSourceBoth     bool    `json:"is_source_both"`
	Description      string  `json:"description,omitempty"`
	PageAge          string  `json:"page_age,omitempty"`
	PageFetched      string  `json:"page_fetched,omitempty"`
	Profile          Profile `json:"profile,omitempty"`
	Language         string  `json:"language,omitempty"`
	IsFamilyFriendly bool    `json:"family_friendly"`
}

type SearchResults struct {
	Result
	Type    string `json:"type"`
	SubType string `json:"subtype"`
	IsLive  bool   `json:"is_live"`

	MetaURL   MetaURL   `json:"meta_url,omitempty"`
	Thumbnail Thumbnail `json:"thumbnail,omitempty"`
}

func main() {
	token := os.Getenv("BRAVE_TOKEN")
	if token == "" {
		fmt.Fprintln(os.Stderr, "BRAVE_TOKEN environment variable not set")
		os.Exit(1)
	}

	args := os.Args
	if len(args) < 2 {
		fmt.Println("Usage: b query string...")
		return
	}

	strs := strings.Join(args[1:], " ")
	fmt.Printf("Args at once: %s", url.QueryEscape(strs))

	req, err := http.NewRequest("GET", baseURI, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		os.Exit(20)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Subscription-Token", token)

	query := req.URL.Query()
	query.Add("q", url.QueryEscape(strs))
	req.URL.RawQuery = query.Encode()

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		fmt.Fprintf(os.Stderr, "An error during request")
		os.Exit(1)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error durng response read")
		os.Exit(1)
	}

	var o WebSearchApiResponse

	err = json.Unmarshal(data, &o)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(len(o.Web.Results))
}
