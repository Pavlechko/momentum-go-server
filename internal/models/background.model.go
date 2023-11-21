package models

type Src struct {
	Original string `json:"landscape"`
}

type Urls struct {
	Regular string `json:"regular"`
}

type SourceURL struct {
	Image string `json:"html"`
}

type Photographer struct {
	Name string `json:"name"`
}

type PexelsPhoto struct {
	Photographer string `json:"photographer"`
	Alt          string `json:"alt"`
	Image        Src    `json:"src"`
	SourseURL    string `json:"url"`
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
	SourceURL    SourceURL    `json:"links"`
}

type FrontendBackgroundImageResponse struct {
	Photographer string `json:"photographer"`
	Image        string `json:"image"`
	Alt          string `json:"alt"`
	Source       string `json:"source"`
	SourceURL    string `json:"source_url"`
}

type BackgroundInput struct {
	Source string `json:"source" binding:"required"`
}

var BackgroundProviders = []string{"unsplash.com", "pexels.com"}
