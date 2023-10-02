package mongoutils

import (
	"time"
)

// BackupModel backup implementation
type BackupModel struct {
	Checksum   string     `bson:"checksum" json:"checksum"`
	LastBackup *time.Time `bson:"last_backup" json:"last_backup"`
}

func (model *BackupModel) ToMap() map[string]any {
	panic("overide ToMap method for backup")
}

func (model *BackupModel) CanBackup() bool {
	return len(model.ToMap()) > 0
}

func (model *BackupModel) MD5() string {
	if !model.CanBackup() {
		return ""
	}
	cs := NewChecksum(model.ToMap())
	return cs.MD5()
}

func (model *BackupModel) SetChecksum(v string) {
	model.Checksum = v
}

func (model *BackupModel) GetChecksum() string {
	return model.Checksum
}

func (model *BackupModel) NeedBackup() bool {
	return model.LastBackup == nil
}

func (model *BackupModel) MarkBackup() {
	t := time.Now()
	model.LastBackup = &t
}

func (model *BackupModel) UnMarkBackup() {
	model.LastBackup = nil
}
