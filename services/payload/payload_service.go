package payload

type payloadService struct {
	Payload map[string]*Data
}

func NewPayloadService() Service {
	return &payloadService{
		Payload: make(map[string]*Data, 0),
	}
}

func (s payloadService) Set(key string, payload *Data) {
	s.Payload[key] = payload
}

func (s payloadService) Get(key string) (*Data, bool) {
	val, found := s.Payload[key]
	return val, found
}

func (s payloadService) Remove(key string) {
	delete(s.Payload, key)
}
