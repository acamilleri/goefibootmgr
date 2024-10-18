package goefibootmgr

// BootEntry represents and EFI boot number entry
type BootEntry struct {
	Num    uint16
	Active bool
	Label  string
}

// ActivateEntry activates the given boot entry
func (manager *BootManager) ActivateEntry(entry *BootEntry) error {
	err := manager.executor.Run("efibootmgr", "-b", bootnumToHexString(entry.Num), "-a")
	if err != nil {
		entry.Active = true
	}

	return err
}

// DeactivateEntry delete the given boot entry
func (manager *BootManager) DeactivateEntry(entry *BootEntry) error {
	err := manager.executor.Run("efibootmgr", "-b", bootnumToHexString(entry.Num), "-A")
	if err != nil {
		entry.Active = false
	}

	return err
}

// DeleteEntry delete the given boot entry
func (manager *BootManager) DeleteEntry(entry *BootEntry) error {
	return manager.executor.Run("efibootmgr", "-b", bootnumToHexString(entry.Num), "-B")
}
