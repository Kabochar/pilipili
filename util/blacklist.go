package util

import (
	"bufio"
	"os"
	"strconv"
	"sync"

	"github.com/fsnotify/fsnotify"
)

var blacklists struct {
	sync.RWMutex
	data map[uint]struct{}
}

func init() {
	blacklists.data = make(map[uint]struct{})
}

func WatchBlackList() {
	loadBlacklist()
	AddFsNotify(os.Getenv("BLACK_LIST_FILE_PATH"), onBlacklistChange)
}

func onBlacklistChange(in fsnotify.Event) {
	const writeOrCreateMask = fsnotify.Write | fsnotify.Create
	if in.Op&writeOrCreateMask == 0 {
		return
	}
	loadBlacklist()
}

func loadBlacklist() {
	filePath := os.Getenv("BLACK_LIST_FILE_PATH")
	fp, err := os.Open(filePath)
	if err != nil {
		Log().Error("open file failed %v", err)
		return
	}
	defer fp.Close()

	data := make(map[uint]struct{}, 20)
	f := bufio.NewReader(fp)
	for {
		line, _, err := f.ReadLine()
		if err != nil {
			break
		}
		uid, err := strconv.ParseInt(string(line), 10, 64)
		if err != nil {
			continue
		}
		if uid <= 0 {
			continue
		}
		data[uint(uid)] = struct{}{}
	}
	if len(data) == 0 {
		return
	}

	blacklists.Lock()
	defer blacklists.Unlock()
	blacklists.data = data
}

func CheckInBlackList(uid uint) bool {
	blacklists.RLock()
	defer blacklists.RUnlock()
	_, ok := blacklists.data[uid]
	return ok
}
