package mongoutils

import "time"

func parseAsBackup(v any) (Backup, bool) {
	i, ok := v.(Backup)
	return i, ok
}

type Backup interface {
	// ToMap get model as map for backup
	// return nil or empty map to skip backup
	ToMap() map[string]any
	// CanBackup check if ToMap method not nil
	CanBackup() bool
	// MD5 calculate md5 checksum for model data
	// Returns empty string if CanBackup return false
	MD5() string
	// SetChecksum set model md5 checksum
	SetChecksum(string)
	// GetChecksum get model md5 checksum
	GetChecksum() string
	// NeedBackup check if record need backup
	NeedBackup() bool
	// MarkBackup set backup state to current date
	MarkBackup()
	// UnMarkBackup set backup state to nil
	UnMarkBackup()
}

// BackupModel backup implementation
type BackupModel struct {
	Checksum   string     `bson:"checksum" json:"checksum"`
	LastBackup *time.Time `bson:"last_backup" json:"last_backup"`
}

func (model BackupModel) ToMap() map[string]any {
	panic("overide ToMap method for backup")
}

func (model BackupModel) CanBackup() bool {
	return len(model.ToMap()) > 0
}

func (model BackupModel) MD5() string {
	if !model.CanBackup() {
		return ""
	}
	cs := NewChecksum(model.ToMap())
	return cs.MD5()
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
