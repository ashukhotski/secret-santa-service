// mongo.go
package service

import (
	"context"
	"errors"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ServiceRepo struct {
	client *mongo.Client
	dbName string
}

func NewServiceRepo(connString string, dbName string) (*ServiceRepo, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(connString))
	if err != nil {
		return nil, err
	}

	err = client.Connect(context.TODO())
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	return &ServiceRepo{
		client,
		dbName,
	}, nil
}

func (r *ServiceRepo) CountMatchedParticipants(ctx context.Context, chid *string, eid *string, tid *string, y int) (int64, error) {
	eidValue := ""
	if eid != nil {
		eidValue = *eid
	}
	tidValue := ""
	if tid != nil {
		tidValue = *tid
	}
	chidValue := ""
	if chid != nil {
		chidValue = *chid
	}
	cName := eidValue + "_" + tidValue + "_" + chidValue + "_" + strconv.Itoa(y)
	collection := r.client.Database(r.dbName).Collection(cName)

	filter := bson.M{"isMatched": true}

	count, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *ServiceRepo) CountAllParticipants(ctx context.Context, chid *string, eid *string, tid *string, y int) (int64, error) {
	eidValue := ""
	if eid != nil {
		eidValue = *eid
	}
	tidValue := ""
	if tid != nil {
		tidValue = *tid
	}
	chidValue := ""
	if chid != nil {
		chidValue = *chid
	}
	cName := eidValue + "_" + tidValue + "_" + chidValue + "_" + strconv.Itoa(y)
	collection := r.client.Database(r.dbName).Collection(cName)

	filter := bson.M{}

	count, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *ServiceRepo) GetAllParticipants(ctx context.Context, chid *string, eid *string, tid *string, y int) ([]Participant, error) {
	eidValue := ""
	if eid != nil {
		eidValue = *eid
	}
	tidValue := ""
	if tid != nil {
		tidValue = *tid
	}
	chidValue := ""
	if chid != nil {
		chidValue = *chid
	}
	cName := eidValue + "_" + tidValue + "_" + chidValue + "_" + strconv.Itoa(y)
	collection := r.client.Database(r.dbName).Collection(cName)

	options := options.Find()

	filter := bson.M{}

	var results []Participant

	cur, err := collection.Find(context.TODO(), filter, options)
	if err != nil {
		return nil, err
	}

	for cur.Next(context.TODO()) {
		var i Participant
		err := cur.Decode(&i)
		if err != nil {
			return nil, err
		}

		results = append(results, i)
	}

	return results, nil
}

func (r *ServiceRepo) GetUnmatchedParticipants(ctx context.Context, chid *string, eid *string, tid *string, y int) ([]Participant, error) {
	eidValue := ""
	if eid != nil {
		eidValue = *eid
	}
	tidValue := ""
	if tid != nil {
		tidValue = *tid
	}
	chidValue := ""
	if chid != nil {
		chidValue = *chid
	}
	cName := eidValue + "_" + tidValue + "_" + chidValue + "_" + strconv.Itoa(y)
	collection := r.client.Database(r.dbName).Collection(cName)

	options := options.Find()

	filter := bson.M{"isMatched": false}

	var results []Participant

	cur, err := collection.Find(context.TODO(), filter, options)
	if err != nil {
		return nil, err
	}

	for cur.Next(context.TODO()) {
		var i Participant
		err := cur.Decode(&i)
		if err != nil {
			return nil, err
		}

		results = append(results, i)
	}

	return results, nil
}

func (r *ServiceRepo) GetParticipantById(ctx context.Context, chid *string, eid *string, tid *string, uid string, y int) (*Participant, error) {
	eidValue := ""
	if eid != nil {
		eidValue = *eid
	}
	tidValue := ""
	if tid != nil {
		tidValue = *tid
	}
	chidValue := ""
	if chid != nil {
		chidValue = *chid
	}
	cName := eidValue + "_" + tidValue + "_" + chidValue + "_" + strconv.Itoa(y)
	collection := r.client.Database(r.dbName).Collection(cName)

	options := options.Find()

	filter := bson.M{"userId": uid}

	var results []Participant

	cur, err := collection.Find(context.TODO(), filter, options)
	if err != nil {
		return nil, err
	}

	for cur.Next(context.TODO()) {
		var i Participant
		err := cur.Decode(&i)
		if err != nil {
			return nil, err
		}

		results = append(results, i)
	}

	if len(results) < 1 {
		return nil, errors.New("no such participant")
	}

	return &results[0], nil
}

func (r *ServiceRepo) RegisterParticipant(ctx context.Context, p *Participant, y int) error {
	mod := mongo.IndexModel{
		Keys:    bson.M{"userId": 1},
		Options: options.Index().SetUnique(true),
	}

	eidValue := ""
	if p.EnterpriseId != nil {
		eidValue = *p.EnterpriseId
	}
	tidValue := ""
	if p.TeamId != nil {
		tidValue = *p.TeamId
	}
	chidValue := ""
	if p.ChannelId != nil {
		chidValue = *p.ChannelId
	}
	cName := eidValue + "_" + tidValue + "_" + chidValue + "_" + strconv.Itoa(y)
	collection := r.client.Database(r.dbName).Collection(cName)
	_, _ = collection.Indexes().CreateOne(ctx, mod)

	_, err := collection.InsertOne(ctx, p)
	if err != nil {
		return err
	}
	return nil
}

func (r *ServiceRepo) UpdateParticipantMatch(ctx context.Context, match *Participant, p *Participant, y int) error {
	eidValue := ""
	if p.EnterpriseId != nil {
		eidValue = *p.EnterpriseId
	}
	tidValue := ""
	if p.TeamId != nil {
		tidValue = *p.TeamId
	}
	chidValue := ""
	if p.ChannelId != nil {
		chidValue = *p.ChannelId
	}
	cName := eidValue + "_" + tidValue + "_" + chidValue + "_" + strconv.Itoa(y)
	collection := r.client.Database(r.dbName).Collection(cName)

	filter := bson.M{
		"userId": p.UserId,
	}

	update := bson.D{
		{"$set", bson.D{
			{"isMatched", true},
			{"yourMatchAddress", match.Address},
			{"yourMatchId", match.UserId},
			{"yourMatchName", match.UserName},
		}},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}
