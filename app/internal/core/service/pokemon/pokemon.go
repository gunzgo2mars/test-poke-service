package pokemon

import (
	"context"
	"math/rand"
	"strconv"

	"go.uber.org/zap"

	"github.com/gunzgo2mars/test-poke-service/app/internal/core/model"
	"github.com/gunzgo2mars/test-poke-service/app/internal/repository/ext"
	"github.com/gunzgo2mars/test-poke-service/app/pkg/logger"
)

type IPokemonService interface {
	RequestPokemonInfo(
		ctx context.Context,
		name string,
	) (*model.PokemonSchema, error)
	RequestPokemonAbilities(
		ctx context.Context,
		name string,
	) (*model.PokemonAbilitiesSchema, error)
	RandomPokemon(ctx context.Context) (*model.PokemonSchema, error)
}

type pokemonService struct {
	extRepo ext.IExtRepository
}

func New(extRepo ext.IExtRepository) IPokemonService {
	return &pokemonService{
		extRepo: extRepo,
	}
}

func (s *pokemonService) RequestPokemonInfo(
	ctx context.Context,
	name string,
) (*model.PokemonSchema, error) {
	cacheResult, err := s.extRepo.CacheGetPokemonData(ctx, name)

	if err != nil || cacheResult == nil {
		repoResult, err := s.extRepo.GetPokemonInfo(ctx, name)
		if err != nil {
			return nil, err
		}
		if err := s.extRepo.CacheSetPokemonData(ctx, repoResult); err != nil {
			return nil, err
		}

		logger.Info(
			ctx,
			"Get data from DB.",
			zap.String("type", "service"),
		)

		return repoResult, nil
	}

	logger.Info(
		ctx,
		"Get data from Cache.",
		zap.String("type", "service"),
	)

	return cacheResult, nil
}

func (s *pokemonService) RequestPokemonAbilities(
	ctx context.Context,
	name string,
) (*model.PokemonAbilitiesSchema, error) {
	cacheResult, err := s.extRepo.CacheGetPokemonAbilities(ctx, name)
	if err != nil || cacheResult == nil {
		repoResult, err := s.extRepo.GetPokemonAbilities(ctx, name)
		if err != nil {
			return nil, err
		}
		if err := s.extRepo.CacheSetPokemonAbilities(ctx, name, repoResult); err != nil {
			return nil, err
		}

		logger.Info(
			ctx,
			"Get data from DB.",
			zap.String("type", "service"),
		)

		return repoResult, nil
	}

	logger.Info(
		ctx,
		"Get data from Cache.",
		zap.String("type", "service"),
	)
	return cacheResult, nil
}

func (s *pokemonService) RandomPokemon(ctx context.Context) (*model.PokemonSchema, error) {
	randomNumber := rand.Intn(100) + 1
	return s.extRepo.GetPokemonInfo(ctx, strconv.Itoa(randomNumber))
}
