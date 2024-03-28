//go:build !solution

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
)

var (
	sdata    string
	athletes []Athlete
)

type Athlete struct {
	Athlete string `json:"athlete"`
	Age     int    `json:"age"`
	Country string `json:"country"`
	Year    int    `json:"year"`
	Date    string `json:"date"`
	Sport   string `json:"sport"`
	Gold    int    `json:"gold"`
	Silver  int    `json:"silver"`
	Bronze  int    `json:"bronze"`
	Total   int    `json:"total"`
}

type AthleteInfo struct {
	Athlete       string         `json:"athlete"`
	Country       string         `json:"country"`
	Medals        Medals         `json:"medals"`
	MedalsByYears map[int]Medals `json:"medals_by_year"`
}

type Medals struct {
	Gold   int `json:"gold"`
	Silver int `json:"silver"`
	Bronze int `json:"bronze"`
	Total  int `json:"total"`
}

type TopCountries struct {
	Country string `json:"country"`
	Gold    int    `json:"gold"`
	Silver  int    `json:"silver"`
	Bronze  int    `json:"bronze"`
	Total   int    `json:"total"`
}

func ParseJSON(w *http.ResponseWriter) {
	contents, err := os.ReadFile(sdata)

	if err != nil {
		http.Error(*w, "error while parsing JSON", 400)
	}

	_ = json.Unmarshal(contents, &athletes)
}

func SearchByName(w *http.ResponseWriter, name string, sport string) AthleteInfo {
	var (
		a            []Athlete
		info         AthleteInfo
		medalsByYear = make(map[int]Medals)
	)

	currentCountry := ""
	target := name
	i := 0

	for _, athlete := range athletes {
		if athlete.Athlete == target {
			if sport == "" || athlete.Sport == sport {
				if i == 0 {
					currentCountry = athlete.Country
					a = append(a, athlete)
				} else {
					if athlete.Country == currentCountry {
						a = append(a, athlete)
					}
				}

				i++
			}
		}
	}

	if len(a) == 0 {
		http.Error(*w, "invalid athlete name", 404)
		(*w).WriteHeader(404)

		return AthleteInfo{}
	}

	info.Athlete = target
	info.Country = currentCountry

	for _, v := range a {
		info.Medals.Gold += v.Gold
		info.Medals.Silver += v.Silver
		info.Medals.Bronze += v.Bronze
		info.Medals.Total += v.Total

		medalsByYear[v.Year] = Medals{
			Gold:   medalsByYear[v.Year].Gold + v.Gold,
			Silver: medalsByYear[v.Year].Silver + v.Silver,
			Bronze: medalsByYear[v.Year].Bronze + v.Bronze,
			Total:  medalsByYear[v.Year].Total + v.Total,
		}
	}

	info.MedalsByYears = medalsByYear

	return info
}

func topInSports(w *http.ResponseWriter, sport string, limit int) []AthleteInfo {
	used := make(map[string]bool)
	a := make([]AthleteInfo, 0)

	for _, athlete := range athletes {
		if athlete.Sport == sport && !used[athlete.Athlete] {
			a = append(a, SearchByName(w, athlete.Athlete, athlete.Sport))
			used[athlete.Athlete] = true
		}
	}

	if len(a) == 0 {
		http.Error(*w, "invalid sport name", 404)
		(*w).WriteHeader(404)

		return make([]AthleteInfo, 0)
	}

	sort.Slice(a, func(i, j int) bool {
		if a[i].Medals.Gold == a[j].Medals.Gold {
			if a[i].Medals.Silver == a[j].Medals.Silver {
				if a[i].Medals.Bronze == a[j].Medals.Bronze {
					return a[i].Athlete < a[j].Athlete
				}

				return a[i].Medals.Bronze > a[j].Medals.Bronze
			}

			return a[i].Medals.Silver > a[j].Medals.Silver
		}

		return a[i].Medals.Gold > a[j].Medals.Gold
	})

	if len(a) < limit {
		return a
	}

	return a[:limit]
}

func TopCountriesPerYear(w *http.ResponseWriter, year int, limit int) []TopCountries {
	countries := make(map[string]TopCountries)
	sortedCountries := make([]TopCountries, 0)
	sortedWithoutNull := make([]TopCountries, 0)
	used := make(map[string]bool)

	for _, athlete := range athletes {
		if athlete.Year == year {
			countries[athlete.Country] = TopCountries{
				Country: athlete.Country,
				Gold:    countries[athlete.Country].Gold + athlete.Gold,
				Silver:  countries[athlete.Country].Silver + athlete.Silver,
				Bronze:  countries[athlete.Country].Bronze + athlete.Bronze,
				Total:   countries[athlete.Country].Total + athlete.Total,
			}
		}
	}

	if len(countries) == 0 {
		http.Error(*w, "invalid year", 404)
		(*w).WriteHeader(404)

		return make([]TopCountries, 0)
	}

	for _, athlete := range athletes {
		if !used[athlete.Country] {
			sortedCountries = append(sortedCountries, countries[athlete.Country])
			used[athlete.Country] = true
		}
	}

	sort.Slice(sortedCountries, func(i, j int) bool {
		if sortedCountries[i].Gold == sortedCountries[j].Gold {
			if sortedCountries[i].Silver == sortedCountries[j].Silver {
				if sortedCountries[i].Bronze == sortedCountries[j].Bronze {
					return sortedCountries[i].Country < sortedCountries[j].Country
				}

				return sortedCountries[i].Bronze > sortedCountries[j].Bronze
			}

			return sortedCountries[i].Silver > sortedCountries[j].Silver
		}

		return sortedCountries[i].Gold > sortedCountries[j].Gold
	})

	for _, country := range sortedCountries {
		if country.Country != "" {
			sortedWithoutNull = append(sortedWithoutNull, country)
		}
	}

	if len(sortedWithoutNull) < limit {
		return sortedWithoutNull
	}

	return sortedWithoutNull[:limit]
}

func Handle(w http.ResponseWriter, r *http.Request) {
	urlStr := "http://" + r.Host + r.URL.String()
	u, err := url.Parse(urlStr)

	if err != nil {
		log.Fatal(err)
	}

	q := u.Query()
	ParseJSON(&w)

	if len(q["name"]) > 0 {
		info := SearchByName(&w, q["name"][0], "")

		if info.Athlete != "" {
			infoJSON, _ := json.Marshal(info)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			_, _ = fmt.Fprintln(w, string(infoJSON))
		}
	} else if len(q["sport"]) > 0 {
		limit := 3

		if len(q["limit"]) > 0 {
			limit, err = strconv.Atoi(q["limit"][0])

			if err != nil {
				http.Error(w, "invalid limit", 400)
				w.WriteHeader(400)
			}
		}

		top := topInSports(&w, q["sport"][0], limit)

		if len(top) > 0 {
			topJSON, _ := json.Marshal(top)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			_, _ = fmt.Fprintln(w, string(topJSON))
		}
	} else if len(q["year"]) > 0 {
		limit := 3

		if len(q["limit"]) > 0 {
			limit, err = strconv.Atoi(q["limit"][0])

			if err != nil {
				http.Error(w, "invalid limit", 400)
				w.WriteHeader(400)
			}
		}

		year, _ := strconv.Atoi(q["year"][0])
		top := TopCountriesPerYear(&w, year, limit)

		if len(top) > 0 {
			currentTopJSON, _ := json.Marshal(top)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			_, _ = fmt.Fprintln(w, string(currentTopJSON))
		}
	}
}

func main() {
	port := flag.Int("port", 8000, "port string")
	data := flag.String("data", "", "data string")
	flag.Parse()
	sdata = *data
	http.HandleFunc("/", Handle)
	log.Fatal(http.ListenAndServe("localhost:"+strconv.Itoa(*port), nil))
}
