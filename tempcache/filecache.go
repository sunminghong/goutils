package tempcache

import (
	"path/filepath"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/sunminghong/goutils"
	log "github.com/sunminghong/goutils/logger"
)

type File struct {
	sync.Mutex

	file *os.File
}

type FileCache struct {
	dir string
	dirCache string
	dirBak string

	file *File

	filenamePrefix      string
	filename            string
	maxBytes            int64
	intervalMinuteFlushData int

	processHandle processFunc

	in chan []byte
	quit chan bool
  stop bool

	lock sync.Mutex

  isSyncing bool
}

func NewFileTempCache( dir string, filenamePrefix string) *FileCache {
	var lock sync.Mutex

	fc := &FileCache{
		dir:                 dir,
    dirCache:            dir + "/cache",
    dirBak:              dir + "/bak",
		filename:            filenamePrefix,
		filenamePrefix:      filenamePrefix,

		file: &File{},
		in:   make(chan []byte, 20),
		quit: make(chan bool),
		lock: lock,
    stop: false,

    isSyncing: false,
	}

	return fc
}

func (fc *FileCache) Init( process processFunc, intervalMinuteFlushData int) {
  fc.processHandle=       process
  fc.intervalMinuteFlushData = intervalMinuteFlushData

  os.MkdirAll(fc.dirCache, os.ModeDir)
  os.MkdirAll(fc.dirBak, os.ModeDir)

	fc.init()
}

func (fc *FileCache) init() {
  //启动写入数据缓存goroutine
	go func() {
		for {
			select {
			case in ,ok := <-fc.in:
        if ok {
				fc._append(in)
      } else {
        fc.quit <-true
        return
      }

			}
		}
	}()

  //启动 定时上报handler callback
  go func() {
   // return
		ticker := time.NewTicker(time.Duration(fc.intervalMinuteFlushData) * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
        log.Debug("===============================timeticker filecache sync callback-------")
        go fc.SyncCache()
			}
		}
	}()

	log.Debug("xxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
	fc.SyncCache()
	log.Debug("yyyyyyyyyyyyyyyyyyyyyyyyyyyyy")

	fc.openCache()
}

func (fc *FileCache) openCache() {
	now := utils.Strftime(time.Now(), "060102150405")
	file, err := os.OpenFile(fmt.Sprintf("%s/%s_%s.che", fc.dirCache, fc.filename, now), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	var lock sync.Mutex
	fc.file = &File{
		lock,
		file,
	}
}

func (fc *FileCache) SyncCache() {
  if fc.isSyncing {
    log.Info("syncCache is running!!!!!!!!!!", fc.dirCache)
    return
  }
  fc.isSyncing = true
	log.Info("syncCache startting", fc.dirCache)

	fc.lock.Lock()
	defer func() {
    fc.lock.Unlock()
    fc.isSyncing = false
  }()

	tmps, err := utils.ListDir(fc.dirCache, "che")
	if err != nil {
		log.Error("read offset tmp error:", err)
		return
	}

	if fc.file.file != nil {
		fc.file.file.Close()
    fc.openCache()
	}

	for _, tmp := range tmps {
		log.Debug("aaaaaaaaaaaaaaaaaaaaaaaa:", tmp, fc.filenamePrefix)
		if strings.Index(tmp, fc.filenamePrefix) != 0 {
			log.Debug("filenameprefix is not exists:", tmp, fc.filenamePrefix)
			continue
		}

		//先更名
    fi, err := os.Stat(fc.dirCache+"/"+tmp)
    if err != nil {
      continue
    }

    if fi.Size() == 0 {
			os.Remove(fc.dirCache+"/"+tmp)
      continue
    }

    var dofile string
    if strings.Index(tmp,"_doing") == -1 {
      dofile = fmt.Sprintf("%s/%s_doing.che", fc.dirCache, tmp[:len(tmp)-4])
      os.Rename(fc.dirCache+"/"+tmp, dofile)
    } else {
      dofile = fmt.Sprintf("%s/%s", fc.dirCache, tmp)
    }

		/*
			fi, err := os.Open(dofile)
			if err != nil {
				log.Error("read "+tmp+" error:", err)
				continue
			}
			fi.Seek(0, 0)
			content, _ := ioutil.ReadAll(fi)
			log.Debug("ccccccccccccccccccccccccccc", len(content))
			rel := fc.processHandle(dofile, content)

		log.Debug("bbbbbbbbbbbbbbbbbbbbbbbbb")
		*/
    oldfile := strings.Replace(dofile, "_doing","", -1)

    rel := fc.processHandle(dofile, filepath.Base(oldfile)[:len(filepath.Base(oldfile))-4], []byte{})

		log.Debug("dddddddddddddddddddddddd")
		if rel {
      bakfile := fmt.Sprintf("%s/bak_%s.che", fc.dirBak, tmp[:len(tmp)-4])
      os.Rename(dofile, bakfile)
		}
	}

}

func (fc *FileCache) Quit() {
  log.Info("filecache start quitting...")
  fc.stop = true
  close(fc.in)
  <-fc.quit
  fc.file.file.Close()
  log.Info("filecache quited")
}

func (fc *FileCache) Append(data []byte) {
  if !fc.stop {
    fc.in <- data
  }
}

func (fc *FileCache) _append(data []byte) {
	fc.lock.Lock()
	//fc.file.Lock()
	if _, err := fc.file.file.Write(data); err != nil {
		log.Fatal(err)
	}
	//fc.file.Unlock()
  fc.lock.Unlock()
}
