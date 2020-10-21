package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Catalog struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type RegularCatalog struct {
	Catalog
	Questions []Question `json:"questions"`
}

type Question struct {
	ID       int      `json:"id"`
	Question string   `json:"question"`
	Options  []string `json:"options"`
}

type OCAICatalog struct {
	Catalog
	Categories []OCAICategory `json:"categories"`
}

type OCAICategory struct {
	ID    int          `json:"id"`
	Title string       `json:"title"`
	A     OCAIQuestion `json:"a"`
	B     OCAIQuestion `json:"b"`
	C     OCAIQuestion `json:"c"`
	D     OCAIQuestion `json:"d"`
}
type OCAIQuestion struct {
	Question  string `json:"question"`
	Now       int    `json:"now"`
	Preferred int    `json:"preferred"`
}

func getAllCatalogs() []*RegularCatalog {
	return catalogs
}

func getCatalogByID(id int) (*RegularCatalog, error) {
	for _, c := range getAllCatalogs() {
		if c.ID == id {
			return c, nil
		}
	}
	return nil, fmt.Errorf("could not get catalog by id %d", id)
}

func getOCAICatalog() *OCAICatalog {
	return ocai
}

func NewRegularCatalog(name, in, out string, options []string) *RegularCatalog {
	inFile, err := os.OpenFile(in, os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	var i int
	cata := new(RegularCatalog)
	cata.Name = name
	for scanner.Scan() {
		question := scanner.Text()
		i++
		cata.Questions = append(cata.Questions, Question{ID: i, Question: question, Options: options})
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	content, _ := json.MarshalIndent(cata, "", " ")
	ioutil.WriteFile(out, content, 0644)
	cata.ID = len(catalogs)
	return cata
}

func (c *OCAICatalog) WriteJSON(filepath string) {
	content, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile(filepath, content, 0644)
}

func generateOCAICatalog() *OCAICatalog {
	cata := new(OCAICatalog)
	cata.Name = "Organizational Culture Assessment Instrument"
	cata.Categories = append(cata.Categories, OCAICategory{
		ID:    1,
		Title: "Dominant Characteristics",
		A: OCAIQuestion{
			Question: "The organization is a very personal place. It is like an extended family. People seem to share a lot of themselves.",
		},
		B: OCAIQuestion{
			Question: "The organization is a very dynamic entrepreneurial place. People are	willing to stick their necks out and take risks.",
		},
		C: OCAIQuestion{
			Question: "The organization is very results oriented. A major concern is with getting the job done. People are very competitive and achievement	oriented.",
		},
		D: OCAIQuestion{
			Question: "The organization is a very controlled and structured place. Formal procedures generally govern what people do.",
		},
	})

	cata.Categories = append(cata.Categories, OCAICategory{
		ID:    2,
		Title: "Organizational Leadership",
		A: OCAIQuestion{
			Question: "The leadership in the organization is generally considered to exemplify mentoring, facilitating, or nurturing.",
		},
		B: OCAIQuestion{
			Question: "The leadership in the organization is generally considered to exemplify	entrepreneurship, innovating, or risk taking.",
		},
		C: OCAIQuestion{
			Question: "The leadership in the organization is generally considered to exemplify a no-nonsense, aggressive, results-oriented focus.",
		},
		D: OCAIQuestion{
			Question: "The leadership in the organization is generally considered to exemplify coordinating, organizing, or smooth-running efficiency.",
		},
	})
	cata.Categories = append(cata.Categories, OCAICategory{
		ID:    3,
		Title: "Management of Employees",
		A: OCAIQuestion{
			Question: "The management style in the organization is characterized by teamwork, consensus, and participation.",
		},
		B: OCAIQuestion{
			Question: "The management style in the organization is characterized by individual risk-taking, innovation, freedom, and uniqueness.",
		},
		C: OCAIQuestion{
			Question: "The management style in the organization is characterized by hard-driving competitiveness, high demands, and achievement.",
		},
		D: OCAIQuestion{
			Question: "The management style in the organization is characterized by security of employment, conformity, predictability, and stability in relationships.",
		},
	})
	cata.Categories = append(cata.Categories, OCAICategory{
		ID:    4,
		Title: "Organization Glue",
		A: OCAIQuestion{
			Question: "The glue that holds the organization together is loyalty and mutual trust. Commitment to this organization runs high.",
		},
		B: OCAIQuestion{
			Question: "The glue that holds the organization together is commitment to innovation and development. There is an emphasis on being on the cutting edge.",
		},
		C: OCAIQuestion{
			Question: "The glue that holds the organization together is the emphasis on achievement and goal accomplishment. Aggressiveness and winning	are common themes",
		},
		D: OCAIQuestion{
			Question: "The glue that holds the organization together is formal rules and policies. Maintaining a smooth-running organization is important.",
		},
	})
	cata.Categories = append(cata.Categories, OCAICategory{
		ID:    5,
		Title: "Strategic Emphases",
		A: OCAIQuestion{
			Question: "The organization emphasizes human development. High trust, openness, and participation persist.",
		},
		B: OCAIQuestion{
			Question: "The organization emphasizes acquiring new resources and creating new	challenges. Trying new things and prospecting for opportunities are	valued.",
		},
		C: OCAIQuestion{
			Question: "The organization emphasizes competitive actions and achievement.	Hitting stretch targets and winning in the marketplace are dominant.",
		},
		D: OCAIQuestion{
			Question: "The organization emphasizes permanence and stability. Efficiency, control and smooth operations are important.",
		},
	})
	cata.Categories = append(cata.Categories, OCAICategory{
		ID:    6,
		Title: "Criteria of Success",
		A: OCAIQuestion{
			Question: "The organization defines success on the basis of the development of human resources, teamwork, employee commitment, and concern for people.",
		},
		B: OCAIQuestion{
			Question: "The organization defines success on the basis of having the most unique or newest products. It is a product leader and innovator.",
		},
		C: OCAIQuestion{
			Question: "The organization defines success on the basis of winning in the marketplace and outpacing the competition. Competitive market leadership is key.",
		},
		D: OCAIQuestion{
			Question: "The organization defines success on the basis of efficiency.	Dependable delivery, smooth scheduling and low-cost production are critical.",
		},
	})
	return cata
}
