package service

import (
	"context"
	"go-restful/model"
	"go-restful/model/entity"
	"go-restful/repo"
)

type TodoService interface {
	Find(ctx context.Context) ([]entity.Todo, error)
	FindByID(ctx context.Context, id uint) (entity.Todo, error)
	Create(ctx context.Context, name string) (entity.Todo, error)
	Update(ctx context.Context, id uint, name string) error
	Delete(ctx context.Context, id uint) error
}

type todoService struct {
	todoRepositorty repo.TodoRepository
}

func NewTodoService() TodoService {
	db := model.GetDB()
	return &todoService{
		todoRepositorty: repo.NewTodoRepository(db),
	}
}

func (service *todoService) Find(ctx context.Context) ([]entity.Todo, error) {
	return service.todoRepositorty.Find(ctx)
}

func (service *todoService) FindByID(ctx context.Context, id uint) (entity.Todo, error) {
	return service.todoRepositorty.FindByID(ctx, id)
}

func (service *todoService) Create(ctx context.Context, name string) (entity.Todo, error) {
	todo := entity.NewTodo(name)
	err := service.todoRepositorty.Create(ctx, &todo)
	return todo, err
}

func (service *todoService) Update(ctx context.Context, id uint, name string) error {
	updateData := entity.Todo{Name: name}
	return service.todoRepositorty.Update(ctx, id, updateData)
}

func (service *todoService) Delete(ctx context.Context, id uint) error {
	return service.todoRepositorty.Delete(ctx, id)
}
