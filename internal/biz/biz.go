package biz

import (
	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	NewAuthUsecase,
	NewUserUseCase,
	NewCalendarUseCase,
	NewEventUseCase,
	NewEventHistoryUseCase,
	NewGoogleUseCase,
	NewOpenAIUseCase,
	NewChatUseCase,
)
