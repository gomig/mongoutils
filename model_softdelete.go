package mongoutils

import "time"

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

func (model SoftDeleteModel) IsDeleted() bool {
	return model.DeletedAt != nil
}
