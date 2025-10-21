package usecase

import (
	"context"
	"time"

	"3xui-bot/internal/core"
	"3xui-bot/internal/ports"
)

type UserUseCase struct {
	userRepo ports.UserRepo
	clock    ports.Clock
}

func NewUserUseCase(userRepo ports.UserRepo, clock ports.Clock) *UserUseCase {

	return &UserUseCase{
		userRepo: userRepo,
		clock:    clock,
	}
}

func (uc *UserUseCase) CreateUser(ctx context.Context, dto CreateUserDTO) (*core.User, error) {
	existingUser, err := uc.userRepo.GetUserByTelegramID(ctx, dto.TelegramID)
	if err == nil && existingUser != nil {

		return existingUser, nil
	}

	newUser := &core.User{
		TelegramID:   dto.TelegramID,
		Username:     dto.Username,
		FirstName:    dto.FirstName,
		LastName:     dto.LastName,
		LanguageCode: dto.LanguageCode,
		HasTrial:     false,
		IsBlocked:    false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = uc.userRepo.CreateUser(ctx, newUser)
	if err != nil {

		return nil, err
	}

	return newUser, nil
}

func (uc *UserUseCase) GetUser(ctx context.Context, telegramID int64) (*core.User, error) {

	return uc.userRepo.GetUserByTelegramID(ctx, telegramID)
}

func (uc *UserUseCase) GetUserByID(ctx context.Context, userID int64) (*core.User, error) {

	return uc.userRepo.GetUserByID(ctx, userID)
}

func (uc *UserUseCase) UpdateUser(ctx context.Context, user *core.User) error {
	user.UpdatedAt = time.Now()

	return uc.userRepo.UpdateUser(ctx, user)
}

func (uc *UserUseCase) ActivateTrial(ctx context.Context, userID int64) (bool, error) {
	user, err := uc.userRepo.GetUserByID(ctx, userID)
	if err != nil {

		return false, err
	}

	if user.HasTrial {

		return false, ErrUserTrialAlreadyUsed
	}

	err = uc.userRepo.MarkTrialAsUsed(ctx, userID)
	if err != nil {

		return false, err
	}

	return true, nil
}
