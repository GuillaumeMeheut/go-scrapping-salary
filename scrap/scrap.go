package scrap

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

var (
	formatRegex, _ = regexp.Compile(`[^\dàÀ]+`)
)

func formatText(c string) int {
	a := strings.Contains(c, "Jusqu")
	b := strings.Contains(c, "partir")
	yearly := strings.Contains(c, "par an")
	if yearly && !a && !b {
		text := formatRegex.ReplaceAllString(c, "")
		return calculAverage(text)
	}
	return 0
}

func calculAverage(str string) int {
	tokens := strings.Split(strings.ToLower(str), "à")
	switch len(tokens) {
	case 1:
		value, err := strconv.Atoi(tokens[0])
		if err != nil {
			fmt.Println(err)
			return 0
		}
		return value
	case 2:
		min, err := strconv.Atoi(tokens[0])
		if err != nil {
			fmt.Println(err)
			return 0
		}
		max, err := strconv.Atoi(tokens[1])
		if err != nil {
			fmt.Println(err)
			return 0
		}
		return (min + max) / 2
	}
	fmt.Println("No value")
	return 0
}

func calculTotalAverage(salaries []int) int {
	avg := 0
	//Il faut delete les 0 avant
	for _, v := range salaries {
		if v != 0 {
			avg += v
		}
	}
	return avg / len(salaries)
}

func FetchData(search string) int {

	var salaries []int

	// Create a Collector specifically for Shopify
	c := colly.NewCollector(colly.AllowedDomains("fr.indeed.com"))

	//Get the content
	c.OnHTML("ul.jobsearch-ResultsList table.jobCard_mainContent td.resultContent div.salaryOnly", func(e *colly.HTMLElement) {
		content := e.ChildText("div.salary-snippet-container > div.attribute_snippet")

		t := formatText(content)
		if t != 0 {
			salaries = append(salaries, t)
		}

	})

	//On next button we go to the next page
	c.OnHTML("a[aria-label=Suivant]", func(e *colly.HTMLElement) {
		next_page := e.Request.AbsoluteURL(e.Attr("href"))
		c.Visit(next_page)
	})

	//Each request we print the URL
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(r.URL.String())
	})

	// Start the collector to that link
	c.Visit("https://fr.indeed.com/jobs?q=" + search)

	totalAvgSalary := calculTotalAverage(salaries)

	fmt.Println("Collected", len(salaries), "salary")
	fmt.Printf("Le salaire moyen pour le terme %v est de %v", search, totalAvgSalary)
	return totalAvgSalary
}
