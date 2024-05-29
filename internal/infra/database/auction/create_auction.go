package auction

import (
	"context"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}
type AuctionRepository struct {
	Collection *mongo.Collection
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	repo := &AuctionRepository{
		Collection: database.Collection("auctions"),
	}
	go repo.autoCloseExpiredAuctions()
	return repo
}

func (ar *AuctionRepository) CreateAuction(
	ctx context.Context,
	auctionEntity *auction_entity.Auction) *internal_error.InternalError {
	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}
	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internal_error.NewInternalServerError("Error trying to insert auction")
	}

	return nil
}

func getAuctionDuration() time.Duration {
	auctionDuration := os.Getenv("AUCTION_DURATION")
	duration, err := time.ParseDuration(auctionDuration)
	if err != nil {
		return time.Minute * 30 // Padrão para 30 minutos
	}
	return duration
}

func (ar *AuctionRepository) autoCloseExpiredAuctions() {
	ticker := time.NewTicker(time.Second * 5) // Verifica a cada 5s
	defer ticker.Stop()

	for range ticker.C {
		logger.Info("Checking for expired auctions")
		ar.closeExpiredAuctions()
	}
}

func (ar *AuctionRepository) closeExpiredAuctions() {
	ctx := context.Background()
	auctionDuration := getAuctionDuration()
	now := time.Now().Unix()

	filter := bson.M{
		"status":    auction_entity.Active,
		"timestamp": bson.M{"$lt": now - int64(auctionDuration.Seconds())},
	}
	update := bson.M{
		"$set": bson.M{"status": auction_entity.Completed},
	}

	_, err := ar.Collection.UpdateMany(ctx, filter, update)
	if err != nil {
		logger.Error("Error trying to update expired auctions", err)
	}
}
