package cashbook

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type payment struct {
	Id          string  `bson:"id"`
	Amount      float64 `bson:"amount"`
	Description string  `bson:"description"`
	Payer       string  `bson:"payer"`
}

func (p *payment) toModel() *Payment {
	return &Payment{p.Id, p.Amount, p.Description, p.Payer}
}

type cashbookDocument struct {
	ID           primitive.ObjectID `bson:"_id"`
	TripName     string             `bson:"tripName"`
	Payments     []*payment         `bson:"payments"`
	Participants []string           `bson:"participants"`
}

func (c *cashbookDocument) toModel() *Cashbook {
	payments := make([]*Payment, 0)
	for _, p := range c.Payments {
		payments = append(payments, p.toModel())
	}
	return &Cashbook{
		Id:       c.ID.Hex(),
		TripName: c.TripName,
		Payments: payments,
		People:   c.Participants,
	}
}

func ConnectToDatabase(ctx context.Context) *mongo.Client {
	log.Default().Println("connecting to database")

	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	dbUri := os.Getenv("CASHBOOK_DB_URI")
	if dbUri == "" {
		panic("no CASHBOOK_DB_URI was passed")
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(dbUri).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(err)
	}
	return client
}

func DisconnectFromDatabase(ctx context.Context, client *mongo.Client) {
	log.Default().Println("disconnecting from database")

	if err := client.Disconnect(ctx); err != nil {
		panic(err)
	}
}

func (s *Server) createNewCashbook(ctx context.Context, tripName string) (*Cashbook, error) {
	log.Default().Println("creating a new cashbook", tripName)

	coll := s.db.Database("cashbook").Collection("cashbooks")
	if coll == nil {
		err := s.db.Database("cashbook").CreateCollection(ctx, "cashbooks")
		if err != nil {
			panic("could not create collection")
		}
		coll = s.db.Database("cashbook").Collection("cashbooks")
	}
	result, err := coll.InsertOne(ctx, cashbookDocument{
		ID:           primitive.NewObjectID(),
		TripName:     tripName,
		Payments:     make([]*payment, 0),
		Participants: make([]string, 0),
	})
	if err != nil {
		panic("could not create cashbook")
	}

	log.Default().Println("created cashbook with id", result.InsertedID)

	id, _ := result.InsertedID.(primitive.ObjectID)
	return &Cashbook{id.Hex(), tripName, make([]*Payment, 0), make([]string, 0)}, nil
}

func (s *Server) findById(ctx context.Context, id string) (*Cashbook, error) {
	log.Default().Println("find cashbook by id", id)

	objectId, _ := primitive.ObjectIDFromHex(id)

	coll := s.db.Database("cashbook").Collection("cashbooks")

	result := coll.FindOne(ctx, bson.M{"_id": objectId})
	if result.Err() != nil {
		log.Default().Println("could find a cashbook with id", id, result.Err().Error())
		return nil, errors.New("not-found")
	}

	foundCashbook := cashbookDocument{}
	decodeError := result.Decode(&foundCashbook)
	if decodeError != nil {
		log.Default().Println("could not decode cashbook", decodeError.Error())
	} else {
		log.Default().Println("foundCashbook", foundCashbook)
	}

	return foundCashbook.toModel(), nil
}

func (s *Server) createNewPayment(ctx context.Context, cashbookId string, p *Payment) (*Cashbook, error) {
	log.Default().Println("add payment", p, "to cashbook with id", cashbookId)

	objectId, _ := primitive.ObjectIDFromHex(cashbookId)

	coll := s.db.Database("cashbook").Collection("cashbooks")

	p.Id = uuid.NewString()
	updateResult, err := coll.UpdateByID(ctx, objectId, bson.D{{"$push", bson.D{{"payments", p}}}})
	if err != nil {
		log.Default().Println("could not update document", err.Error())
		return nil, err
	}
	log.Default().Println(updateResult.ModifiedCount, "cashbooks were updated")

	result := coll.FindOne(ctx, bson.M{"_id": objectId})
	if result.Err() != nil {
		log.Default().Println("could find a cashbook with id", objectId, result.Err().Error())
		return nil, errors.New("not-found")
	}

	foundCashbook := cashbookDocument{}
	decodeError := result.Decode(&foundCashbook)
	if decodeError != nil {
		log.Default().Println("could not decode cashbook", decodeError.Error())
	} else {
		log.Default().Println("foundCashbook", foundCashbook)
	}

	return foundCashbook.toModel(), nil
}

func (s *Server) createNewParticipant(ctx context.Context, cashbookId string, participantName string) (*Cashbook, error) {
	log.Default().Println("add participant", participantName, "to cashbook with id", cashbookId)

	objectId, _ := primitive.ObjectIDFromHex(cashbookId)

	coll := s.db.Database("cashbook").Collection("cashbooks")

	updateResult, err := coll.UpdateByID(ctx, objectId, bson.D{{"$push", bson.D{{"participants", participantName}}}})
	if err != nil {
		log.Default().Println("could not update document", err.Error())
		return nil, err
	}
	log.Default().Println(updateResult.ModifiedCount, "cashbooks were updated")

	result := coll.FindOne(ctx, bson.M{"_id": objectId})
	if result.Err() != nil {
		log.Default().Println("could find a cashbook with id", objectId, result.Err().Error())
		return nil, errors.New("not-found")
	}

	foundCashbook := cashbookDocument{}
	decodeError := result.Decode(&foundCashbook)
	if decodeError != nil {
		log.Default().Println("could not decode cashbook", decodeError.Error())
	} else {
		log.Default().Println("foundCashbook", foundCashbook)
	}

	return foundCashbook.toModel(), nil
}
