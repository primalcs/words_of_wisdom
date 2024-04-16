package interfaces

import "time"

type POWChallengeBuilder interface {
	GenerateRandomChallenge(currentTime time.Time, resource string) (POWChallenge, string)
	GenerateChallengeById(currentTime time.Time, resource, randomId string) POWChallenge
	GetChallengeDuration() time.Duration
}

type POWChallenge interface {
	Solve(int64) error
	Verify() (bool, error)
	GetResource() string // to check client's id
	GetRand() string     // to check client's token if applicable
	GetDate() time.Time  // to check if challenge expired
	ToJSON() (string, error)
	FromJSON(string) error
}
