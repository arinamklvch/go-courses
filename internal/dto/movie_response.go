package dto

type Movie struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Year     int    `json:"year"`
	Director string `json:"director"`
}

type MovieResponse struct {
	ID int `json:"id"`
	*Movie
	Error string `json:"error,omitempty"`
}
