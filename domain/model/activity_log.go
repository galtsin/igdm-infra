package model

import "time"

// Создание/Удаление/Остановка/Запуск аккаунта
// Изменение учетных данных: пароль, прокси
// Попытки авторизации и статус авторизации (challenge, two factor)
// Дата планового полного обхода ящика

type ActivityLog struct {
	AccountID string
	Log       string
	CreatedAt time.Time
}

func NewActivityLog(accountID string) ActivityLog {
	return ActivityLog{
		AccountID: accountID,
		CreatedAt: time.Now(),
	}
}

func (al *ActivityLog) SetLog(log string) {
	al.Log = log
}
