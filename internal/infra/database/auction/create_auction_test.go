package auction

import (
	"context"
	"fullcycle-auction_go/configuration/database/mongodb"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"os"
	"testing"
	"time"
)

func TestAutoCloseExpiredAuctions(t *testing.T) {
	// Configurar a conexão com o banco de dados de teste
	os.Setenv("MONGODB_URL", "mongodb://admin:admin@mongodb:27017/auction_test?authSource=admin")
	os.Setenv("MONGODB_DB", "auction_test")
	db, err := mongodb.NewMongoDBConnection(context.Background())
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer db.Client().Disconnect(context.Background())

	auctionRepo := NewAuctionRepository(db)

	// Limpar coleções antes do teste
	db.Collection("auctions").Drop(context.Background())

	// Configurar variável de ambiente para duração do leilão
	os.Setenv("AUCTION_DURATION", "1s")

	// Criar leilão de teste
	auctionEntity := &auction_entity.Auction{
		Id:          "test_auction",
		ProductName: "Test Product",
		Category:    "Test Category",
		Description: "Test Description",
		Condition:   auction_entity.New,
		Status:      auction_entity.Active,
		Timestamp:   time.Now(),
	}
	auctionRepo.CreateAuction(context.Background(), auctionEntity)

	// Esperar até o leilão expirar
	time.Sleep(2 * time.Second)

	// Executar a função de fechamento automático
	auctionRepo.closeExpiredAuctions()

	// Verificar se o leilão foi fechado
	updatedAuction, _ := auctionRepo.FindAuctionById(context.Background(), "test_auction")
	if updatedAuction == nil {
		t.Fatalf("Failed to find auction: %v", err)
	}

	if updatedAuction.Status != auction_entity.Completed {
		t.Fatalf("Expected auction to be completed, but got %v", updatedAuction.Status)
	}
}
