package errors

import (
	"errors"
	"fmt"
)

// Стандартные ошибки Go для проверки типов
var (
	// Используем errors из стандартной библиотеки
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
	ErrInvalidInput  = errors.New("invalid input")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrForbidden     = errors.New("forbidden")
	ErrInternal      = errors.New("internal error")
	ErrDatabase      = errors.New("database error")
	ErrExternalAPI   = errors.New("external API error")
	ErrTimeout       = errors.New("timeout error")
)

// Error представляет собой расширенную типизированную ошибку
type Error struct {
	// Оригинальная ошибка
	err error
	// Тип ошибки для проверки
	errorType error
	// Сообщение об ошибке
	message string
	// Дополнительный контекст в виде ключ-значение
	fields map[string]interface{}
}

// Error реализует интерфейс error
func (e *Error) Error() string {
	if e.message != "" {
		if e.err != nil {
			return fmt.Sprintf("%s: %s", e.message, e.err.Error())
		}
		return e.message
	}
	return e.err.Error()
}

// Unwrap позволяет работать с errors.Is и errors.As
func (e *Error) Unwrap() error {
	return e.err
}

// Is проверяет тип ошибки, интегрируется с errors.Is
func (e *Error) Is(target error) bool {
	if e.errorType != nil && errors.Is(e.errorType, target) {
		return true
	}
	return errors.Is(e.err, target)
}

// WithField добавляет дополнительную информацию к ошибке
func (e *Error) WithField(key string, value interface{}) *Error {
	if e.fields == nil {
		e.fields = make(map[string]interface{})
	}
	e.fields[key] = value
	return e
}

// Fields возвращает все поля ошибки
func (e *Error) Fields() map[string]interface{} {
	return e.fields
}

// GetType возвращает тип ошибки
func (e *Error) GetType() error {
	return e.errorType
}

// Конструкторы различных типов ошибок

// NewNotFound создает ошибку типа "не найдено"
func NewNotFound(entity string, id interface{}) *Error {
	return &Error{
		errorType: ErrNotFound,
		message:   fmt.Sprintf("%s with ID %v not found", entity, id),
		fields: map[string]interface{}{
			"entity": entity,
			"id":     id,
		},
	}
}

// NewDatabaseError создает ошибку базы данных
func NewDatabaseError(err error) *Error {
	return &Error{
		err:       err,
		errorType: ErrDatabase,
		message:   "database error: ",
	}
}

// NewInvalidInput создает ошибку невалидного ввода
func NewInvalidInput(details string, fields map[string]interface{}) *Error {
	e := &Error{
		errorType: ErrInvalidInput,
		message:   fmt.Sprintf("invalid input: %s", details),
	}
	if fields != nil {
		e.fields = fields
	}
	return e
}

// NewInternal создает внутреннюю ошибку
func NewInternal(err error, message string) *Error {
	return &Error{
		err:       err,
		errorType: ErrInternal,
		message:   message,
	}
}

// NewExternalAPIError создает ошибку внешнего API
func NewExternalAPIError(err error, service string, details string) *Error {
	return &Error{
		err:       err,
		errorType: ErrExternalAPI,
		message:   fmt.Sprintf("error from %s: %s", service, details),
		fields: map[string]interface{}{
			"service": service,
		},
	}
}

// Wrap оборачивает существующую ошибку с дополнительным контекстом
func Wrap(err error, message string) *Error {
	if err == nil {
		return nil
	}

	// Если это уже наша ошибка, копируем её тип
	var typedErr *Error
	if errors.As(err, &typedErr) {
		return &Error{
			err:       typedErr,
			errorType: typedErr.errorType,
			message:   message,
			fields:    typedErr.fields,
		}
	}

	// Иначе создаем новую ошибку
	return &Error{
		err:     err,
		message: message,
	}
}

// Is проверяет тип ошибки
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As приводит ошибку к указанному типу
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}
