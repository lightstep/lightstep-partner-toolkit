package flags

import (
	"go.uber.org/zap"
	"sync"
)

type FlagManager struct {
	flags map[string]*Flag

	sync.Mutex
}

var Manager *FlagManager

func init() {
	Manager = NewFlagManager()
}

func NewFlagManager() *FlagManager {
	return &FlagManager{flags: make(map[string]*Flag)}
}

func (fm *FlagManager) Clear() {
	fm.Lock()
	fm.flags = make(map[string]*Flag)
	fm.Unlock()
}

func (fm *FlagManager) GetFlags() map[string]*Flag {
	fm.Lock()
	defer fm.Unlock()
	return fm.flags
}

func (fm *FlagManager) LoadFlags(configFlags []Flag, logger *zap.Logger) {
	fm.Lock()
	defer fm.Unlock()

	for _, f := range configFlags {
		flag := f
		flag.Setup(logger)
		fm.flags[flag.Name] = &flag
	}
}

func (fm *FlagManager) FlagCount() int {
	fm.Lock()
	defer fm.Unlock()
	return len(fm.flags)
}

func (fm *FlagManager) GetFlag(name string) *Flag {
	fm.Lock()
	defer fm.Unlock()
	return fm.flags[name]
}
