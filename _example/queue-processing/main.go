package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
	"hq.0xa1.red/axdx/scheduler/internal/platform/database"
	"hq.0xa1.red/axdx/scheduler/internal/platform/database/models"
	"hq.0xa1.red/axdx/scheduler/internal/platform/database/redis"
	"hq.0xa1.red/axdx/scheduler/internal/schedule"
)

func main() {
	// Create 10 random users
	// Add random number of scheduled items for all
	// Process queue
	itemIDs := []uuid.UUID{
		uuid.New(),
		uuid.New(),
	}
	items := map[uuid.UUID]string{
		itemIDs[0]: "barracks",
		itemIDs[1]: "granary",
	}

	rand.Seed(time.Now().UnixNano())

	users := make(map[uuid.UUID]int)
	for i := 0; i < 50; i++ {
		users[uuid.New()] = rand.Intn(150) + 20 // nolint
	}

	redis.SetDatabase(9)
	database.SetBackend(string(database.KindRedis))
	db, err := database.New()
	if err != nil {
		log.Panic(err)
	}

	redis.Flush()
	startGen := time.Now()
	wg := sync.WaitGroup{}
	for id, number := range users {
		for i := 0; i < number; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				minutes := time.Duration(rand.Intn(300) + 1) // nolint
				tt := time.Now().Add(minutes * time.Second)
				item := itemIDs[rand.Intn(2)] // nolint
				message := models.NewMessage(tt, "queue", item, id)
				if err := db.Schedule(context.Background(), message); err != nil {
					log.Panic(err)
				}
			}()
		}
	}
	fmt.Printf("[*] Generating took %s\n\n", time.Since(startGen))
	wg.Wait()

	t := time.NewTicker(time.Second)
	for range t.C {
		start := time.Now()
		queue, errors := schedule.Collect(context.Background())
		log.Printf("[*] Time of collection: %s", time.Now().Format("15:04:05"))
		if len(errors) > 0 {
			log.Panicf("%+v", errors)
		}
		for i, message := range queue {
			fmt.Printf("%d - %s - %s: %s\n", i+1, message.Timestamp.Format("15:04:05"), message.ID.String(), items[message.ItemID])
			if err := db.Acknowledge(context.Background(), message.ID, message.OwnerID); err != nil {
				log.Printf("ERROR: failed to acknowledge message %s", message.ID.String())
			}
		}
		fmt.Printf("[*] Collection took %s\n\n", time.Since(start))
	}

}
