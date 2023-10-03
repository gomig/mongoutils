package mongoutils

import "time"

type SoftDelete interface {
	// SoftDelete set deleted_at field to current date
	SoftDelete()
	// Restore set deleted_at field to nil
	Restore()
	// IsDeleted check if item soft deleted
	IsDeleted() bool
}

type SoftDeleteModel struct {
	DeletedAt *time.Time `bson:"deleted_at" json:"deleted_at"`
}

func (model *SoftDeleteModel) SoftDelete() {
	t := time.Now()
	model.DeletedAt = &t
}

func (model *SoftDeleteModel) Restore() {
	model.DeletedAt = nil
}

func (model *SoftDeleteModel) IsDeleted() bool {
	return model.DeletedAt != nil
}
