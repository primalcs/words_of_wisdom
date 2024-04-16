package hashcash

import (
	"time"

	"github.com/primalcs/words_of_wisdom/internal/utils"

	"github.com/primalcs/words_of_wisdom/pkg/interfaces"
)

type Builder struct {
	ZerosCount        int64
	ChallengeDuration time.Duration
}

func (builder *Builder) GenerateRandomChallenge(currentTime time.Time, resource string) (interfaces.POWChallenge, string) {
	randomId := utils.GenerateRandomString()
	return &Challenge{
		Version:    1,
		ZerosCount: builder.ZerosCount,
		Date:       currentTime,
		Resource:   resource,
		Rand:       randomId,
		Counter:    0,
	}, randomId
}

func (builder *Builder) GenerateChallengeById(currentTime time.Time, resource, randomId string) interfaces.POWChallenge {
	return &Challenge{
		Version:    1,
		ZerosCount: builder.ZerosCount,
		Date:       currentTime,
		Resource:   resource,
		Rand:       randomId,
		Counter:    0,
	}
}

func (builder *Builder) GetChallengeDuration() time.Duration {
	return builder.ChallengeDuration
}
