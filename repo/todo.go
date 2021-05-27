package repo

import (
	"context"
	"go-restful/model/entity"

	"gorm.io/gorm"
)

type TodoRepository interface {
	Find(ctx context.Context) ([]entity.Todo, error)
	FindByID(ctx context.Context, id uint) (entity.Todo, error)
	Create(ctx context.Context, todo *entity.Todo) error
	Update(ctx context.Context, id uint, updateData entity.Todo) error
	Delete(ctx context.Context, id uint) error
}

type todoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(
	db *gorm.DB,
) TodoRepository {
	return &todoRepository{
		db: db,
	}
}

func (repository *todoRepository) Find(ctx context.Context) ([]entity.Todo, error) {
	var todos = []entity.Todo{}
	result := repository.db.Find(&todos)
	return todos, result.Error
}

func (repository *todoRepository) FindByID(ctx context.Context, id uint) (entity.Todo, error) {
	var todo = entity.Todo{}
	result := repository.db.First(&todo, "id = ?", id)
	return todo, result.Error
}

func (repository *todoRepository) Create(ctx context.Context, todo *entity.Todo) error {
	result := repository.db.Create(todo)
	return result.Error
}

func (repository *todoRepository) Update(ctx context.Context, id uint, updateData entity.Todo) error {
	result := repository.db.Model(entity.Todo{
		Model: gorm.Model{
			ID: id,
		},
	}).Updates(updateData)
	return result.Error
}

func (repository *todoRepository) Delete(ctx context.Context, id uint) error {
	result := repository.db.Delete(&entity.Todo{
		Model: gorm.Model{
			ID: id,
		},
	})
	return result.Error
}
