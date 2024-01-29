package events

import (
	"context"
	"fmt"
)

type AccessTokenCreatedEvent struct {
	TokenID string
	UserID  string
	ClientID string
}

func PublishAccessTokenCreatedEvent(ctx context.Context, tokenID, userID, clientID string) {
	event := &AccessTokenCreatedEvent{
		TokenID: tokenID,
		UserID:  userID,
		ClientID: clientID,
	}

	fmt.Println("Publishing AccessTokenCreated event:", event)

	// Publish the event using a pub/sub system or other mechanism
}


