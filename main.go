package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"
	"time"
)

type Article struct {
	Name    string
	Tags    []string
	Date    string
	Content string
}

const pageName = "learning"
const dateBegin = 10

func mustLoadTags() map[string]int {
	frequentTags := `
{
    "quia": 15,
    "molestiae": 7,
    "voluptates": 4,
    "dignissimos": 5,
    "occaecati": 3,
    "doloribus": 3,
    "ad": 2,
    "laudantium": 1,
    "aut": 1,
    "nulla": 4,
    "nobis": 9,
	"sunt": 6
}
`
	var tags map[string]int
	err := json.Unmarshal([]byte(frequentTags), &tags)
	if err != nil {
		log.Fatal(err)
	}
	return tags
}

func createArticles(tags map[string]int) []*Article {
	min := func() int {
		n := 0
		for _, v := range tags {
			if v > n {
				n = v
			}
		}
		return n
	}()
	var articles []*Article

	body := `
Nostrum molestiae repellendus quidem distinctio debitis et fugiat
Et deserunt voluptatum omnis voluptas optio mollitia officia et
Enim ducimus autem laudantium
Aut porro qui veniam ut aperiam
Est quo nulla nobis sunt impedit earum est.  

### Subtitle  

Qui modi et veniam voluptas maiores quas.  
Omnis ipsa molestiae ad.  
Omnis placeat fuga ut.  
Saepe eos id quae et quod dolore officiis.  
Similique voluptatem iure sit aut rerum et debitis.  
Odio omnis aperiam dolor expedita et in aspernatur sequi.
`

	for i := 0; i < min+1; i++ {
		articles = append(articles, &Article{
			Content: fmt.Sprintf("%v", body),
			Date:    time.Now().AddDate(0, 0, -(i + dateBegin)).Format(time.RFC3339),
			Name:    fmt.Sprintf("%v-%v.md", pageName, i),
		})
	}
	for t, n := range tags {
		arts := articles[0 : n+1]
		for _, a := range arts {
			a.Tags = append(a.Tags, t)
		}
	}
	return articles
}
func stringifyTags(ts []string) []string {
	var s = ""
	for _, t := range ts {
		s += fmt.Sprintf("%+q,", t)
	}
	return []string{strings.TrimSuffix(s, ",")}

}
func cleanWorkingDir() {
	dir := fmt.Sprintf("content/%v", pageName)
	err := os.RemoveAll(dir)
	if err != nil {
		log.Fatal(err)
	}
	err = os.Mkdir(dir, fs.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}
func writeArticles(articles []*Article) {
	cleanWorkingDir()
	for _, a := range articles {
		t := strings.Split(a.Name, ".")
		t = strings.Split(t[0], "-")
		title := fmt.Sprintf("%v %v", t[0], t[1])
		tags := stringifyTags(a.Tags)
		str := fmt.Sprintf("---\ntitle: %q\ndate: %v\ntags: %v\n---\n%v", title, a.Date, tags, a.Content)

		f, err := os.Create(fmt.Sprintf("content/%v/%v", pageName, a.Name))
		defer func() {
			err := f.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()
		if err != nil {
			log.Fatal(err)
		}
		_, err = f.WriteString(str)
		if err != nil {
			log.Fatal(err)
		}

	}
}

func main() {
	tags := mustLoadTags()
	articles := createArticles(tags)
	writeArticles(articles)
}
