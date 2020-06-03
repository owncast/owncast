package main

type Message struct {
	Author string `json:"author"`
	Body   string `json:"body"`
	Image  string `json:"image"`
	Id     string `json:"id"`
}

func (self *Message) String() string {
	return self.Author + " says " + self.Body
}
