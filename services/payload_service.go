package services

type payloadService struct {
	Payload map[string]*PayloadData
}

func NewPayload() Payload {
	return &payloadService{
		Payload: make(map[string]*PayloadData, 0),
	}
}

func (s payloadService) Set(key string, payload *PayloadData) {
	s.Payload[key] = payload
}

func (s payloadService) Get(key string) (*PayloadData, bool) {
	val, found := s.Payload[key]
	return val, found
}

func (s payloadService) Remove(key string) {
	delete(s.Payload, key)
}
