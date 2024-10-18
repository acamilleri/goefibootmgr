package goefibootmgr

import (
	"regexp"
	"strings"
)

// BootManager holds all info returned from the efibootmgr command
type BootManager struct {
	// BootEntry stored in the EFI BootCurrent variable
	BootCurrent *BootEntry
	// BootEntry stored in the EFI BootNext variable
	BootNext *BootEntry
	// Slice containing BootEntrys in order they will boot
	BootOrder []BootEntry
	// Slice containing all detected boot entries
	BootEntries []BootEntry

	executor ExecCommand
}

// BootManagerOption type to override default init value of
// BootManager type
type BootManagerOption func(manager *BootManager)

// WithCustomExecutor override default with a custom executor
func WithCustomExecutor(executor ExecCommand) BootManagerOption {
	return func(manager *BootManager) {
		manager.executor = executor
	}
}

// NewBootManager initialize the boot manager
func NewBootManager(opts ...BootManagerOption) (*BootManager, error) {
	manager := &BootManager{executor: defaultExecutor{}}

	// override default values
	for _, opt := range opts {
		opt(manager)
	}

	// initialize manager values
	err := bootInfo(manager)
	if err != nil {
		return nil, err
	}
	return manager, nil
}

// bootInfo runs the efibootmanager command and returns an BootInfo struct
// containing all the info that was returned
func bootInfo(manager *BootManager) error {
	out, err := manager.executor.Output("efibootmgr")
	if err != nil {
		return err
	}

	bootCurrentRe := regexp.MustCompile(`BootCurrent: ([0-9a-fA-F]{4})`)
	bootNextRe := regexp.MustCompile(`BootNext: ([0-9a-fA-F]{4})`)
	bootOrderRe := regexp.MustCompile(`BootOrder: ([0-9a-fA-F]{4}(?:,[0-9a-fA-F]{4})*)`)
	bootEntryRe := regexp.MustCompile(`Boot([0-9a-fA-F]{4})(\*?)\s+(.*)`)

	lines := strings.Split(string(out), "\n")

	var bootEntryMap = map[uint16]BootEntry{}

	// First find boot entries
	for _, line := range lines {
		if match := bootEntryRe.FindStringSubmatch(line); match != nil {
			e := BootEntry{
				Num:    hexStringToBootNum(match[1]),
				Active: match[2] == "*",
				Label:  match[3],
			}
			bootEntryMap[e.Num] = e
			manager.BootEntries = append(manager.BootEntries, e)
		}
	}

	// Now parse the rest of the boot info
	for _, line := range lines {
		if match := bootCurrentRe.FindStringSubmatch(line); match != nil {
			num := hexStringToBootNum(match[1])
			entry := bootEntryMap[num]
			manager.BootCurrent = &entry
		} else if match := bootNextRe.FindStringSubmatch(line); match != nil {
			num := hexStringToBootNum(match[1])
			entry := bootEntryMap[num]
			manager.BootNext = &entry
		} else if match := bootOrderRe.FindStringSubmatch(line); match != nil {
			numList := match[1]
			numSlice := strings.Split(numList, ",")
			for _, num := range numSlice {
				num := hexStringToBootNum(num)
				manager.BootOrder = append(manager.BootOrder, bootEntryMap[num])
			}
		}
	}

	return nil
}

// SetBootOrder sets the EFI BootOrder variable using a list of args containing
// boot numbers
func (manager *BootManager) SetBootOrder(bo ...uint16) (err error) {
	if len(bo) == 0 {
		// Delete BootOrder
		err = manager.executor.Run("efibootmgr", "-O")
	} else {
		s := bootnumToHexString(bo[0])
		for _, o := range bo[1:] {
			s += "," + bootnumToHexString(o)
		}

		err = manager.executor.Run("efibootmgr", "-o", s)
	}

	return
}

// SetBootNext sets the EFI BootNext variable to the bootnum passed to it
func (manager *BootManager) SetBootNext(b uint16) (err error) {
	return manager.executor.Run("efibootmgr", "-n", bootnumToHexString(b))
}

// DeleteBootNext deletes the EFI BootNext variable
func (manager *BootManager) DeleteBootNext() error {
	return manager.executor.Run("efibootmgr", "-N")
}
