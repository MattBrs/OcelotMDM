package command_action

import "go.mongodb.org/mongo-driver/bson/primitive"

type CommandAction struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	Name            string             `bson:"name"`
	Description     string             `bson:"description"`
	RequiredOnlne   bool               `bson:"required_online"`
	DefaultPriority uint               `bson:"default_priority"`
	PayloadRequired bool               `bson:"payload_required"`
	TokenRequired   bool               `bson:"token_required"`
}
