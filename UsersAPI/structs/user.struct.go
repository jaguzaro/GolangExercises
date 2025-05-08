package structs

type RandomUserResponse struct {
	Results []User `json:"results"`
	Info    Info   `json:"info"`
}

type User struct {
	Gender   string   `json:"gender"`
	Name     Name     `json:"name"`
	Location Location `json:"location"`
	Email    string   `json:"email"`
	Phone    string   `json:"phone"`
}

type Name struct {
	Title string `json:"title"`
	First string `json:"first"`
	Last  string `json:"last"`
}

type Location struct {
	Street  Street `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	Country string `json:"country"`
}

type Street struct {
	Number int    `json:"number"`
	Name   string `json:"name"`
}

type Info struct {
	Seed    string `json:"seed"`
	Results int    `json:"results"`
	Page    int    `json:"page"`
	Version string `json:"version"`
}
