package mongoutils

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

// parseAsBackup check if value implement backup intraface
func parseAsBackup(v any) (Backup, bool) {
	i, ok := v.(Backup)
	return i, ok
}
