package container

import (
	// internal package
	"phonebook/internal/phonebook"
	repo "phonebook/internal/phonebook/repository"

	// thirdparty package
	"github.com/go-kit/kit/log"
)

type Container struct {
	PhoneBook *phonebook.Service
}

func CreateContainer(actor string, logger log.Logger) *Container {
	svc := phonebook.NewService(repo.NewPostgres(), actor, logger)
	return &Container{
		PhoneBook: svc,
	}
}
