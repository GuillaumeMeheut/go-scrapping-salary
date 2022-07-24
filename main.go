package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"regexp"
	"strconv"
	"strings"
)

func formatText(str string) string {
	return regexp.MustCompile("[a-zA-Z'àÀ ]").ReplaceAllString(str, "")
}

func calculAverage(salaries []string) {
	// average := int

	for _, v := range salaries {
		f := strings.Split(v, "€")

		e, errq := strconv.Unquote(f[0])
		if errq != nil {
			fmt.Println("error on delete quote")
			panic(errq)
		}
		i, err := strconv.Atoi(e)
		if err != nil {
			// ... handle error
			panic(err)
		}
		fmt.Println(i)
		fmt.Println(f[0])
	}

}

func main() {

	search := "nextjs"
	salaries := []string{}

	// Create a Collector specifically for Shopify
	c := colly.NewCollector(colly.AllowedDomains("fr.indeed.com"))

	//contents
	c.OnHTML("ul.jobsearch-ResultsList table.jobCard_mainContent td.resultContent div.salaryOnly", func(e *colly.HTMLElement) {
		content := e.ChildText("div.salary-snippet-container > div.attribute_snippet")

		a := strings.Contains(content, "Jusqu")
		b := strings.Contains(content, "partir")
		yearly := strings.Contains(content, "par an")
		if yearly && !a && !b {
			text := formatText(content)
			salaries = append(salaries, text)
		}
	})

	//next button
	c.OnHTML("a[aria-label=Suivant]", func(e *colly.HTMLElement) {
		next_page := e.Request.AbsoluteURL(e.Attr("href"))
		c.Visit(next_page)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println(r.URL.String())
	})

	// Start the collector
	c.Visit("https://fr.indeed.com/jobs?q=" + search)

	fmt.Println("All known salaries:")
	for _, s := range salaries {
		fmt.Println("\t", s)
	}
	fmt.Println("Collected", len(salaries), "salary")
	calculAverage(salaries)

}
