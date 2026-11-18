package data

type Movie struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Year     int    `json:"year"`
	Director string `json:"director"`
}

type MovieRequestResult struct {
	ID    int
	Movie Movie
	Error error
}

type MovieResponse struct {
	ID int `json:"id"`
	*Movie
	Error string `json:"error,omitempty"`
}
