package events

import (
	"context"
	"fmt"
)

type RefreshTokenCreatedEvent struct {
	RefreshTokenID string
	AccessTokenID  string
}

func PublishRefreshTokenCreatedEvent(ctx context.Context, refreshTokenID, accessTokenID string) {
	event := &RefreshTokenCreatedEvent{
		RefreshTokenID: refreshTokenID,
		AccessTokenID:  accessTokenID,
	}

	fmt.Println("Publishing RefreshTokenCreated event:", event)

	// Publish the event using a pub/sub system or other mechanism
}


