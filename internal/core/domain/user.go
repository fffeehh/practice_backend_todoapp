package domain

import (
	"fmt"
	"regexp"

	core_errors "github.com/fffeehh/practice_backend_todoapp/internal/core/errors"
)

type User struct {
	ID int
	Version int  // Версия может использоваться не только репозиторием, но и уровнем сервиса, транспорта и даже фронтендом

	FullName string
	PhoneNumber *string  // Здесь у нас указатель на строку, потому что FullName у нас обязательный, а вот номер может быть не задан.
	// в таком случае у нас в переменной будет хранится nil, это удобно. А если номер задан, но строка на которую ссылается указатель - значение номера телефона
}

// создание проинициализированного пользователя
func NewUser(
	id int,
	version int,
	fullName string,
	phoneNumber *string,
) User {
	return User{
		ID: id,
		Version: version,
		FullName: fullName,
		PhoneNumber: phoneNumber,
	}
}

// Такая конструкция позволяет перекидывать ответственность за принятие решения об id и version
// на домен пользователя. Транспортному уровню этого делать не нужно, нужно лишь передать сюда то
// что ему пришло
func NewUserUninitialized(
	fullName string,
	phoneNumber *string,
) User {
	return NewUser(
		UninitializedID,
		UninitializedVersion,
		fullName,
		phoneNumber,
		)
}


func (u *User) Validate() error {
	fullNameLength := len([]rune(u.FullName))
	if fullNameLength < 3 || fullNameLength > 100 {
		return fmt.Errorf(
			"invalid `FullName` len: %d: %w",
			fullNameLength,
			core_errors.ErrInvalidArgument,
			)
		// обязательно передаем в строке нашу ошибку ErrInvalidArgument, 
		// чтобы в handler.go она смогла определиться правильно и поставился нужный статус код
	}

	if u.PhoneNumber != nil {
		// не забываем разыменовывать PhoneNumber, потому что у нас он представлен указателем
		phoneNumberLen := len([]rune(*u.PhoneNumber))
		if phoneNumberLen < 10 || phoneNumberLen > 15 {
			return fmt.Errorf(
				"invalid `PhoneNumber` len: %d: %w",
				phoneNumberLen,
				core_errors.ErrInvalidArgument,
				)
		}
		
		// проверяем строку телефона на соответствие регулярному выражению
		re := regexp.MustCompile(`^\+[0-9]+$`)
	
		if !re.MatchString(*u.PhoneNumber) {
			return fmt.Errorf(
				"invalid `PhoneNumber` format: %w",
				core_errors.ErrInvalidArgument,
				)
		}

	}
	return nil
}

type UserPatch struct {
	FullName Nullable[string]
	PhoneNumber Nullable[string]
}

func (p *UserPatch) Validate() error {
	// поля FullName обязательно, поэтому его нельзя удалять
	if p.FullName.Set && p.FullName.Value == nil {
		return fmt.Errorf(
			"`FullName` can't be patched to NULL: %w",
			core_errors.ErrInvalidArgument,
			)
	}
	return nil
}

// применение обновления
func (u *User) ApplyPatch(patch UserPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate user patch: %w", err)
	}

	// после применения обновления, необходимо полностью провалидировать полученного в результате пользователя. 
	// но может быть такое, что в итоге он валидацию не пройдет и придется откатывать изменения. 
	// чтобы такого не было, создаем временную переменную пользователя.
	tmp := *u

	if patch.FullName.Set {
		tmp.FullName = *patch.FullName.Value
	}

	if patch.PhoneNumber.Set {
		tmp.PhoneNumber = patch.PhoneNumber.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched user: %w", err)
	}

	*u = tmp
	
	return nil
}
