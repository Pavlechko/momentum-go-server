package models

type Src struct {
	Original string `json:"landscape"`
}

type PexelsPhoto struct {
	Photographer string `json:"photographer"`
	Alt          string `json:"alt"`
	Image        Src    `json:"src"`
	SourseUrl    string `json:"url"`
}

type PexelsImageResponse struct {
	Photos   []PexelsPhoto `json:"photos"`
	NextPage string        `json:"next_page"`
	PrevPage string        `json:"prev_page"`
}

type FrontendBackgroundImageResponse struct {
	Photographer string `json:"photographer"`
	Image        string `json:"image"`
	Alt          string `json:"alt"`
	Sourse       string `json:"sourse"`
	SourseUrl    string `json:"sourse_url"`
}
