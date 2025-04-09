package cmd

import (
	"fmt"
	"labyrinth/internal/controllers"
	"labyrinth/internal/router"
	"labyrinth/internal/state"
	"log"
	"net"
	"os"

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

	r := router.NewRouter().
		WithState(state.NewState().
			WithJwtSession([]byte(os.Getenv("JWT_SESSION_KEY"))).
			WithPostgresDriver(os.Getenv("POSTGRES_URL")),
		)

	r.SrvConfig.Host = host

	controllers.HandleAll(r)
	if err := r.Run(); err != nil {

		log.Fatal("Server killed due to an error: ", err)

	}
	fmt.Println("Labyrinth server is up and running at port :3100")
}
