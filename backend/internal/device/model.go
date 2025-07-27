package device

type Device struct {
	ID        string   `bson:"_id,omitempty"`
	Name      string   `bson:"name"`
	Type      string   `bson:"type"`
	Status    string   `bson:"status"`
	IPAddress string   `bson:"ip_address"`
	LastSeen  int64    `bson:"last_seen"`
	Tags      []string `bson:"tags,omitempty"`
}
