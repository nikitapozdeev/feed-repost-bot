package storage

type Storage interface {
	Get(clientId int64) ([]Subscription, error)
	Add(s Subscription) error
	Update(s Subscription) error
	Delete(id int64) error

	Close() error
}

type Subscription struct {
	ID       int64
	ClientID int64
	FeedLink string
	Updated  int64
}
