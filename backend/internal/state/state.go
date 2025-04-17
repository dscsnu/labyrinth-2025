package state

import (
	"labyrinth/internal/cache"
	"labyrinth/internal/channel"
	"labyrinth/internal/database"
	"labyrinth/internal/jwtauth"
	"log"
)

type State struct {
	JwtSession *jwtauth.JWTSession
	DB         *database.PostgresDriver
	ChanPool   *channel.ChannelPool
	CM         *cache.CacheManager
}

func NewState() *State {

	return &State{}

}

func (s *State) WithJwtSession(secretKey []byte) *State {

	s.JwtSession = jwtauth.NewJWTSession(secretKey)
	return s

}

func (s *State) WithPostgresDriver(connectionURL string) *State {

	dbconn, err := database.CreatePostgresDriver(connectionURL)
	if err != nil {

		log.Fatal(err)

	}

	s.DB = dbconn

	return s

}

func (s *State) WithChannelPool(chanpool *channel.ChannelPool) *State {

	s.ChanPool = chanpool
	return s

}

func (s *State) WithCacheManager(CM *cache.CacheManager) *State {

	s.CM = CM
	return s

}
