package storage

type IStorage[Entity any] interface {
	Create(Entity) error
	Update(Entity) error
	ReadAll() ([]Entity, error)
	Delete(Entity) error
}
