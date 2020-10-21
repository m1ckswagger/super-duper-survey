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
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
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

func initCatalogs() {
	ocaq := NewRegularCatalog("Organizational Culture Assessment Questionaire", "data/ocaq.txt", "data/ocaq.json", []string{"Yes", "No"})
	ocaq.Description = `The Organizational Culture Assessment Questionnaire (OCAQ) is based on the work of Dr. Talcott
	Parsons, a sociologist at Harvard. Parsons developed a framework and theory of action in social
	systems. He argued that all organizations must carry out four crucial functions if they are to survive for
	any substantial length of time. We have labeled these four functions managing change, achieving
	goals, coordinating teamwork, and building a strong culture. One aspect of the way in which
	organizations achieve their goals is especially important, yet often neglected. This factor has been
	made into a separate, fifth scale: customer orientation.`
	catalogs = append(catalogs, ocaq)
	sheff := NewRegularCatalog("Sheffield Culture Survey", "data/sheffield.txt", "data/sheffield.json", []string{
		"Strongly disagree",
		"Disagree",
		"Neutral",
		"Agree",
		"Strongly Agree",
	})
	sheff.Description = `The questions related to the key elements of a generative culture as follows:
	1) Shared sense of purpose
	2) Structured systems and resources to achieve the purpose
	3) Mindfulness
	4) Processes and mindset for learning and continuous improvement`
	catalogs = append(catalogs, sheff)
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
	cata.Description = `The purpose of the OCAI is to assess six key dimensions of organizational culture. In completing the
	instrument, you will be providing a picture of how your organization operates and the values that
	characterize it. No right or wrong answers exist for these questions, just as there is no right or wrong
	culture. Every organization will most likely produce a different set of responses. Therefore, be as
	accurate as you can in responding to the questions so that your resulting cultural diagnosis will be as
	precise as possible. `
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
