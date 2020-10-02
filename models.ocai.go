package main

import (
	"encoding/json"
	"io/ioutil"
)

type CatalogOCAI struct {
	Categories []CategoryOCAI `json:"categories"`
}

type CategoryOCAI struct {
	ID    int          `json:"id"`
	Title string       `json:"title"`
	A     QuestionOCAI `json:"a"`
	B     QuestionOCAI `json:"b"`
	C     QuestionOCAI `json:"c"`
	D     QuestionOCAI `json:"d"`
}
type QuestionOCAI struct {
	Question  string `json:"question"`
	Now       int    `json:"now"`
	Preferred int    `json:"preferred"`
}

func (c *CatalogOCAI) WriteJSON(filepath string) {
	content, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile(filepath, content, 0644)
}

func generateOCAICatalog() *CatalogOCAI {
	cata := new(CatalogOCAI)

	cata.Categories = append(cata.Categories, CategoryOCAI{
		ID:    1,
		Title: "Dominant Characteristics",
		A: QuestionOCAI{
			Question: "The organization is a very personal place. It is like an extended family. People seem to share a lot of themselves.",
		},
		B: QuestionOCAI{
			Question: "The organization is a very dynamic entrepreneurial place. People are	willing to stick their necks out and take risks.",
		},
		C: QuestionOCAI{
			Question: "The organization is very results oriented. A major concern is with getting the job done. People are very competitive and achievement	oriented.",
		},
		D: QuestionOCAI{
			Question: "The organization is a very controlled and structured place. Formal procedures generally govern what people do.",
		},
	})

	cata.Categories = append(cata.Categories, CategoryOCAI{
		ID:    2,
		Title: "Organizational Leadership",
		A: QuestionOCAI{
			Question: "The leadership in the organization is generally considered to exemplify mentoring, facilitating, or nurturing.",
		},
		B: QuestionOCAI{
			Question: "The leadership in the organization is generally considered to exemplify	entrepreneurship, innovating, or risk taking.",
		},
		C: QuestionOCAI{
			Question: "The leadership in the organization is generally considered to exemplify a no-nonsense, aggressive, results-oriented focus.",
		},
		D: QuestionOCAI{
			Question: "The leadership in the organization is generally considered to exemplify coordinating, organizing, or smooth-running efficiency.",
		},
	})
	cata.Categories = append(cata.Categories, CategoryOCAI{
		ID:    3,
		Title: "Management of Employees",
		A: QuestionOCAI{
			Question: "The management style in the organization is characterized by teamwork, consensus, and participation.",
		},
		B: QuestionOCAI{
			Question: "The management style in the organization is characterized by individual risk-taking, innovation, freedom, and uniqueness.",
		},
		C: QuestionOCAI{
			Question: "The management style in the organization is characterized by hard-driving competitiveness, high demands, and achievement.",
		},
		D: QuestionOCAI{
			Question: "The management style in the organization is characterized by security of employment, conformity, predictability, and stability in relationships.",
		},
	})
	cata.Categories = append(cata.Categories, CategoryOCAI{
		ID:    4,
		Title: "Organization Glue",
		A: QuestionOCAI{
			Question: "The glue that holds the organization together is loyalty and mutual trust. Commitment to this organization runs high.",
		},
		B: QuestionOCAI{
			Question: "The glue that holds the organization together is commitment to innovation and development. There is an emphasis on being on the cutting edge.",
		},
		C: QuestionOCAI{
			Question: "The glue that holds the organization together is the emphasis on achievement and goal accomplishment. Aggressiveness and winning	are common themes",
		},
		D: QuestionOCAI{
			Question: "The glue that holds the organization together is formal rules and policies. Maintaining a smooth-running organization is important.",
		},
	})
	cata.Categories = append(cata.Categories, CategoryOCAI{
		ID:    5,
		Title: "Strategic Emphases",
		A: QuestionOCAI{
			Question: "The organization emphasizes human development. High trust, openness, and participation persist.",
		},
		B: QuestionOCAI{
			Question: "The organization emphasizes acquiring new resources and creating new	challenges. Trying new things and prospecting for opportunities are	valued.",
		},
		C: QuestionOCAI{
			Question: "The organization emphasizes competitive actions and achievement.	Hitting stretch targets and winning in the marketplace are dominant.",
		},
		D: QuestionOCAI{
			Question: "The organization emphasizes permanence and stability. Efficiency, control and smooth operations are important.",
		},
	})
	cata.Categories = append(cata.Categories, CategoryOCAI{
		ID:    6,
		Title: "Criteria of Success",
		A: QuestionOCAI{
			Question: "The organization defines success on the basis of the development of human resources, teamwork, employee commitment, and concern for people.",
		},
		B: QuestionOCAI{
			Question: "The organization defines success on the basis of having the most unique or newest products. It is a product leader and innovator.",
		},
		C: QuestionOCAI{
			Question: "The organization defines success on the basis of winning in the marketplace and outpacing the competition. Competitive market leadership is key.",
		},
		D: QuestionOCAI{
			Question: "The organization defines success on the basis of efficiency.	Dependable delivery, smooth scheduling and low-cost production are critical.",
		},
	})
	return cata
}
