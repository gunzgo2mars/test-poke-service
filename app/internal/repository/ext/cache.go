package ext

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/gunzgo2mars/test-poke-service/app/internal/core/model"
	"github.com/gunzgo2mars/test-poke-service/app/pkg/logger"
)

const (
	pokemonKey        = "POKEMON|%s"
	pokemonAbilityKey = "POKEMON|AB|%s"
)

func (r *extRepository) CacheSetPokemonData(
	ctx context.Context,
	stampData *model.PokemonSchema,
) error {
	redisKey := fmt.Sprintf(pokemonKey, stampData.Name)

	jsonData, err := json.Marshal(stampData)
	if err != nil {
		logger.Error(
			ctx,
			fmt.Sprintf("Error[SET-CACHE]: %s", err.Error()),
			zap.String("type", "repository"),
		)
		return err
	}

	if err := r.cache.HSet(ctx, redisKey, "info", jsonData).Err(); err != nil {
		logger.Error(
			ctx,
			fmt.Sprintf("Error[SET-CACHE]: %s", err.Error()),
			zap.String("type", "repository"),
		)
		return err
	}

	if err := r.cache.Expire(ctx, redisKey, 10*time.Minute).Err(); err != nil {
		logger.Error(
			ctx,
			fmt.Sprintf("Error[SET-CACHE-EXPIRE]: %s", err.Error()),
			zap.String("type", "repository"),
		)
		return err
	}

	return nil
}

func (r *extRepository) CacheGetPokemonData(
	ctx context.Context,
	name string,
) (*model.PokemonSchema, error) {
	redisKey := fmt.Sprintf(pokemonKey, name)

	var pokemonSchema *model.PokemonSchema

	result, err := r.cache.HGet(ctx, redisKey, "info").Result()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(result), &pokemonSchema)
	if err != nil {
		return nil, err
	}

	return pokemonSchema, nil
}

func (r *extRepository) CacheSetPokemonAbilities(
	ctx context.Context,
	name string,
	stampData *model.PokemonAbilitiesSchema,
) error {
	redisKey := fmt.Sprintf(pokemonAbilityKey, name)

	jsonData, err := json.Marshal(stampData)
	if err != nil {
		logger.Error(
			ctx,
			fmt.Sprintf("Error[SET-CACHE]: %s", err.Error()),
			zap.String("type", "repository"),
		)
		return err
	}

	if err := r.cache.HSet(ctx, redisKey, "info", jsonData).Err(); err != nil {
		logger.Error(
			ctx,
			fmt.Sprintf("Error[SET-CACHE]: %s", err.Error()),
			zap.String("type", "repository"),
		)
		return err
	}

	if err := r.cache.Expire(ctx, redisKey, 10*time.Minute).Err(); err != nil {
		logger.Error(
			ctx,
			fmt.Sprintf("Error[SET-CACHE-EXPIRE]: %s", err.Error()),
			zap.String("type", "repository"),
		)
		return err
	}

	return nil
}

func (r *extRepository) CacheGetPokemonAbilities(
	ctx context.Context,
	name string,
) (*model.PokemonAbilitiesSchema, error) {
	redisKey := fmt.Sprintf(pokemonAbilityKey, name)

	var pokemonSchema *model.PokemonAbilitiesSchema

	result, err := r.cache.HGet(ctx, redisKey, "info").Result()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(result), &pokemonSchema)
	if err != nil {
		return nil, err
	}

	return pokemonSchema, nil
}
