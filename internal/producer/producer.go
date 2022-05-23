package producer

type Producer interface {
	Posts(domain string, offset int, count int) ([]Message, error)
}

type Message interface {
	Recipient() string
	HTML() (string, error)
}
