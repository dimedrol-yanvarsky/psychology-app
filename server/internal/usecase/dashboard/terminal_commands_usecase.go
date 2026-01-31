package dashboard

import (
	"strings"
)

// AvailableCommands содержит список поддерживаемых команд терминала
var AvailableCommands = []CommandDescription{
	{Name: "help", Description: "Показать доступные команды и их назначение"},
	{Name: "block user", Description: "Блокировка пользователя"},
	{Name: "delete user", Description: "Удаление пользователя"},
	{Name: "delete account", Description: "Удаление аккаунта и связанных данных"},
	{Name: "change user data", Description: "Обновление имени, почты или статуса пользователя"},
}

// TerminalCommandsUseCase - use case для обработки терминальных команд
type TerminalCommandsUseCase struct{}

// NewTerminalCommandsUseCase создает новый экземпляр TerminalCommandsUseCase
func NewTerminalCommandsUseCase() *TerminalCommandsUseCase {
	return &TerminalCommandsUseCase{}
}

// Execute обрабатывает команду админ-терминала
func (uc *TerminalCommandsUseCase) Execute(input TerminalCommandInput) TerminalCommandOutput {
	normalized := normalizeCommand(input.Command)
	if normalized == "" {
		return TerminalCommandOutput{
			Status:  "error",
			Message: "Команда не может быть пустой",
			Command: "",
		}
	}

	commandKey := strings.ToLower(normalized)
	switch commandKey {
	case "help":
		return TerminalCommandOutput{
			Status:   "success",
			Message:  "Список доступных команд",
			Command:  normalized,
			Commands: AvailableCommands,
		}
	case "block user":
		return TerminalCommandOutput{
			Status:  "success",
			Command: normalized,
			Message: "Команда блокировки пользователя ожидает параметры целевого пользователя",
		}
	case "delete user":
		return TerminalCommandOutput{
			Status:  "success",
			Command: normalized,
			Message: "Команда удаления пользователя ожидает параметры целевого пользователя",
		}
	case "delete account":
		return TerminalCommandOutput{
			Status:  "success",
			Command: normalized,
			Message: "Команда удаления аккаунта ожидает параметры целевого аккаунта",
		}
	case "change user data":
		return TerminalCommandOutput{
			Status:  "success",
			Command: normalized,
			Message: "Команда изменения данных пользователя ожидает параметры для обновления",
		}
	default:
		return TerminalCommandOutput{
			Status:  "error",
			Command: normalized,
			Message: "Команда не найдена",
		}
	}
}

func normalizeCommand(value string) string {
	cleaned := strings.TrimSpace(value)
	if cleaned == "" {
		return ""
	}
	return strings.Join(strings.Fields(cleaned), " ")
}
