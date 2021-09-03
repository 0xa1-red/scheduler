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
	// items := map[uuid.UUID]string{
	// 	itemIDs[0]: "barracks",
	// 	itemIDs[1]: "granary",
	// }

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
				minutes := time.Duration(-1 * (rand.Intn(600) + 1)) // nolint
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
	start := time.Now()
	_, errors := schedule.Collect(context.Background())
	if len(errors) > 0 {
		log.Panicf("%+v", errors)
	}
	// for i, message := range queue {
	// 	fmt.Printf("%d - %s: %s\n", i+1, message.Timestamp.Format(time.RFC1123), items[message.ItemID])
	// }
	fmt.Printf("[*] Collection took %s\n\n", time.Since(start))
}
