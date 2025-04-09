package cmd

import (
	"labyrinth/internal/controllers"
	"labyrinth/internal/router"
	"labyrinth/internal/state"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Execute() {

	// load environment variables
	if err := godotenv.Load(".env"); err != nil {

		log.Fatal(err)

	}

	r := router.NewRouter().
		WithState(state.NewState().
			WithJwtSession([]byte(os.Getenv("JWT_SESSION_KEY"))),
		)

	controllers.HandleAll(r)
	if err := r.Run(); err != nil {

		log.Fatal("Server killed due to an error: ", err)

	}

}
