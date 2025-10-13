// Package logger логер.
package logger

import "go.uber.org/zap"

// MakeLogger конструктор логера приложения.
func MakeLogger() (*zap.SugaredLogger, error) {
	log, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	return log.Sugar(), nil
}

// MakeNop конструктор логера ничего не делающего. Для тестов.
func MakeNop() (*zap.SugaredLogger, error) {
	log := zap.NewNop()
	return log.Sugar(), nil
}
