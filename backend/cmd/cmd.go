package cmd

import (
	"fmt"
	"labyrinth/internal/cache"
	"labyrinth/internal/channel"
	"labyrinth/internal/controllers"
	"labyrinth/internal/router"
	"labyrinth/internal/state"
	"log"
	"net"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func Execute() {

	// load environment variables
	if err := godotenv.Load(".env"); err != nil {

		log.Fatal(err)

	}

	host, err := net.ResolveTCPAddr("tcp", ":3100")
	if err != nil {
		log.Fatal("Failed to resolve address:", err)
	}

	cm := cache.NewManager()

	defaultTTL := 10 * time.Minute
	cleanUpInterval := 20 * time.Minute
	entities := []string{cache.Team, cache.TeamLevel, cache.TeamLevelProgress, cache.TeamUserIDIndex}

	for _, entity := range entities {
		cm.InitCache(entity, defaultTTL, cleanUpInterval)
	}

	r := router.NewRouter().
		WithServerConfig(router.ServerConfig{Host: host}).
		WithState(state.NewState().
			WithJwtSession([]byte(os.Getenv("JWT_SESSION_KEY"))).
			WithChannelPool(channel.NewChannelPool()).
			WithPostgresDriver(os.Getenv("POSTGRES_URL")).
			WithCacheManager(cm),
		)
	controllers.HandleAll(r)
	if err := r.Run(); err != nil {

		log.Fatal("Server killed due to an error: ", err)

	}
	fmt.Println("Labyrinth server is up and running at port :3100")
}
