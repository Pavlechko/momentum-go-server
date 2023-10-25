package models

type BackgroundData struct {
	Pexels   FrontendBackgroundImageResponse
	Unsplash FrontendBackgroundImageResponse
}

type Src struct {
	Original string `json:"landscape"`
}

type Urls struct {
	Regular string `json:"regular"`
}

type SourceUrl struct {
	Image string `json:"html"`
}

type Photographer struct {
	Name string `json:"username"`
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

type UnsplashImageResponse struct {
	Alt          string       `json:"alt_description"`
	Image        Urls         `json:"urls"`
	Photographer Photographer `json:"user"`
	SourceUrl    SourceUrl    `json:"links"`
}

type FrontendBackgroundImageResponse struct {
	Photographer string `json:"photographer"`
	Image        string `json:"image"`
	Alt          string `json:"alt"`
	Source       string `json:"source"`
	SourceUrl    string `json:"source_url"`
}
