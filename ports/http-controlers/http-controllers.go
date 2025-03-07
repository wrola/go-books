package httpControllers

type Book struct {
	ISBN   string `json:"isbn"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var books = []Book{
	{
		ISBN:   "1234567890",
		Title:  "The Great Gatsby",
		Author: "F. Scott Fitzgerald",
	},
}



