package util

import (
	"path"
	"sync"

	"github.com/fsnotify/fsnotify"
)

var (
	fsWatcher            *fsnotify.Watcher
	faEventActionMapping = struct {
		sync.RWMutex
		actionMap map[string]func(in fsnotify.Event)
	}{}
	initFsWatcherFlag = sync.Once{}
)

func newWatcher() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		Log().Panic("FileNotify init new watcher failed [%v]", err)
	}
	fsWatcher = watcher
	faEventActionMapping.actionMap = make(map[string]func(in fsnotify.Event), 0)
}

func getFsEventAction(filename string) func(in fsnotify.Event) {
	faEventActionMapping.RLock()
	defer faEventActionMapping.RUnlock()
	actionEvent, ok := faEventActionMapping.actionMap[filename]
	if !ok {
		return nil
	}
	return actionEvent
}

func FileNotify() {
	defer fsWatcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-fsWatcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write ||
					event.Op&fsnotify.Create == fsnotify.Create {
					eventAction := getFsEventAction(event.Name)
					if eventAction == nil {
						continue
					}
					eventAction(event)
				}
			case err, ok := <-fsWatcher.Errors:
				if !ok {
					return
				}
				Log().Info("FileNotify fsWatcher error: [%v]", err)
			}
		}
	}()
	<-done
}

func AddFsNotify(filepath string, run func(in fsnotify.Event)) {
	if fsWatcher == nil {
		initFsWatcherFlag.Do(func() {
			newWatcher()
			go FileNotify()
		})
	}
	filename := path.Base(filepath)
	err := fsWatcher.Add(filename)
	if err != nil {
		Log().Error("AddFsNotify add new watcher failed %v", err)
	}
	// 注册事件
	faEventActionMapping.Lock()
	defer faEventActionMapping.Unlock()
	faEventActionMapping.actionMap[filename] = run
}
