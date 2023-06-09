package storage

type ITemporaryStorage[Entity any] interface {
	Create(Entity) error
	Delete(Entity) error
	GetByUniqueKey(string) (Entity, error)
}
