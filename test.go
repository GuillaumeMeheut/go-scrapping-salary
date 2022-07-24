package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

var (
	formatRegex = regexp.MustCompile(`[^\dàÀ]+`)
)

type Salary struct {
	Text     string
	Min, Max int
}

func (s *Salary) String() string {
	return fmt.Sprintf("%s, avg = %d", s.Text, s.Avg())
}

func (s *Salary) Avg() int {
	return (s.Min + s.Max) / 2
}

func NewSalary(str string) (s *Salary, err error) {
	result := formatRegex.ReplaceAllString(str, "")
	tokens := strings.Split(strings.ToLower(result), "à")
	switch len(tokens) {
	case 1:
		min, err := strconv.Atoi(tokens[0])
		if err != nil {
			return nil, err
		}

		return &Salary{str, min, min}, nil
	case 2:
		min, err := strconv.Atoi(tokens[0])
		if err != nil {
			return nil, err
		}
		max, err := strconv.Atoi(tokens[1])
		if err != nil {
			return nil, err
		}
		return &Salary{str, min, max}, nil
	}

	return nil, fmt.Errorf(">= 2 tokens: %q", str)
}

func calculAverage(salaries []*Salary) {
	fmt.Println(salaries)
	// for _, salary := range salaries {
	// 	fmt.Println(salary.String())
	// 	fmt.Println(salary.Avg())
	// }
}

func main() {

	search := "nextjs"
	salaries := []*Salary{}

	// Create a Collector specifically for Shopify
	c := colly.NewCollector(colly.AllowedDomains("fr.indeed.com"))

	//contents
	c.OnHTML("ul.jobsearch-ResultsList table.jobCard_mainContent td.resultContent div.salaryOnly", func(e *colly.HTMLElement) {
		content := e.ChildText("div.salary-snippet-container > div.attribute_snippet")

		a := strings.Contains(content, "Jusqu")
		b := strings.Contains(content, "partir")
		yearly := strings.Contains(content, "par an")
		if yearly && !a && !b {
			sal, err := NewSalary(content)
			if err != nil {
				fmt.Println(err)
				return
			}
			salaries = append(salaries, sal)
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
