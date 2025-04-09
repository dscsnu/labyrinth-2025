package state

import "labyrinth/internal/jwtauth"

type State struct {
	JwtSession *jwtauth.JWTSession
}

func NewState() *State {

	return &State{}

}

func (s *State) WithJwtSession(secretKey []byte) *State {

	s.JwtSession = jwtauth.NewJWTSession(secretKey)
	return s

}
