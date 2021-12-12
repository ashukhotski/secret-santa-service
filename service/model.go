// model.go
package service

import (
	"context"
)

const (
	ResponseTypeEphemeral string = "ephemeral"
	ResponseTypeInChannel string = "in_channel"
)

type Participant struct {
	Address          *string `bson:"addresss"`
	ChannelId        *string `bson:"channelId"`
	EnterpriseId     *string `bson:"enterpriseId"`
	IsHost           bool    `bson:"isHost"`
	IsMatched        bool    `bson:"isMatched"`
	ResponseUrl      string  `bson:"responseUrl"`
	TeamId           *string `bson:"teamId"`
	UserId           string  `bson:"userId"`
	UserName         string  `bson:"userName"`
	YourMatchAddress *string `bson:"yourMatchAddress"`
	YourMatchId      *string `bson:"yourMatchId"`
	YourMatchName    *string `bson:"yourMatchName"`
}

type SlackMessage struct {
	ResponseType string `json:"response_type"`
	Text string `json:"text"`
}

type SlackRequest struct {
	ChannelId      *string `json:"channel_id"`
	ChannelName    *string `json:"channel_name"`
	Command        *string `json:"command"`
	EnterpriseId   *string `json:"enterprise_id"`
	EnterpriseName *string `json:"enterprise_name"`
	ResponseUrl    string  `json:"response_url"`
	TeamDomain     *string `json:"team_domain"`
	TeamId         *string `json:"team_id"`
	Text           *string `json:"text"`
	Token          *string `json:"token"`
	TriggerId      *string `json:"trigger_id"`
	UserId         string  `json:"user_id"`
	UserName       string  `json:"user_name"`
}

type SecretSantaRepository interface {
	Check(ctx context.Context, chid *string, eid *string, tid *string, y int) bool
	CountAllParticipants(ctx context.Context, chid *string, eid *string, tid *string, y int) (int64, error)
	CountMatchedParticipants(ctx context.Context, chid *string, eid *string, tid *string, y int) (int64, error)
	GetAllParticipants(ctx context.Context, chid *string, eid *string, tid *string, y int) ([]Participant, error)
	GetUnmatchedParticipants(ctx context.Context, chid *string, eid *string, tid *string, y int) ([]Participant, error)
	GetParticipantById(ctx context.Context, chid *string, eid *string, tid *string, uid string, y int) (*Participant, error)
	RegisterParticipant(ctx context.Context, p *Participant, y int) error
	UpdateParticipantMatch(ctx context.Context, match *Participant, p *Participant, y int) error
}
