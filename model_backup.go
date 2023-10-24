package mongoutils

import "time"

// BackupModel backup implementation
type BackupModel struct {
	Checksum   string     `bson:"checksum" json:"-"`
	LastBackup *time.Time `bson:"last_backup" json:"-"`
}

func (model *BackupModel) SetChecksum(v string) {
	model.Checksum = v
}

func (model BackupModel) GetChecksum() string {
	return model.Checksum
}

func (model BackupModel) NeedBackup() bool {
	return model.LastBackup == nil
}

func (model *BackupModel) MarkBackup() {
	t := time.Now()
	model.LastBackup = &t
}

func (model *BackupModel) UnMarkBackup() {
	model.LastBackup = nil
}
