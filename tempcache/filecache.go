package tempcache

import (
	//"path/filepath"
	"fmt"
		//"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/sunminghong/goutils"
	log "github.com/sunminghong/goutils/logger"
)

type File struct {
	file *os.File
	sync.Mutex
}

type FileCache struct {
	dir string

	file *File

	filenamePrefix      string
	filename            string
	maxBytes            int64
	checkDurationMinute int

	processHandle processFunc

	in chan []byte

	lock sync.Mutex
}

type processFunc func(filename string, data []byte) bool

func NewFileTempCache(
	dir string,
	filenamePrefix string,
	process processFunc,
	checkDurationMinute int,
	maxMB int) *FileCache {

	var lock sync.Mutex

	fc := &FileCache{
		dir:                 dir,
		processHandle:       process,
		maxBytes:            int64(maxMB) * 1024 * 1024,
		filename:            filenamePrefix,
		checkDurationMinute: checkDurationMinute,
		filenamePrefix:      filenamePrefix,

		file: &File{},
		in:   make(chan []byte, 20),
		lock: lock,
	}

	fc.init()

	return fc
}

func (fc *FileCache) init() {
	/*
		go func() {
			ticker := time.NewTicker(time.Duration(fc.checkDurationMinute) * time.Minute)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					fc.checkMax()
				}
			}
		}()
	*/

	go func() {
		for {
			select {
			case in := <-fc.in:
				fc._append(in)
			}
		}
	}()

	log.Debug("xxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
	fc.ReadAndProcess()
	log.Debug("yyyyyyyyyyyyyyyyyyyyyyyyyyyyy")

	////fc.openCache()
}

func (fc *FileCache) checkMax() {
	fi, err := fc.file.file.Stat()
	if err != nil {
		return
	}

	//如果大于最大字节数，就分文件
	if fi.Size() >= fc.maxBytes {
		log.Debug("checkMax:", fi.Size())
		fc.rename()
		fc.openCache()
	}
}

func (fc *FileCache) rename() {
	fc.file.Lock()
	fc.file.file.Close()

	now := utils.Strftime(time.Now(), "060102150405")

	filename := fc.filename
	oldfile := fmt.Sprintf("%s/%s%s.che", fc.dir, filename, now)
	os.Rename(filename, oldfile)
	fc.file.Unlock()

}

func (fc *FileCache) openCache() {
	file, err := os.OpenFile(fmt.Sprintf("%s/%s.che", fc.dir, fc.filename), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	var lock sync.Mutex
	fc.file = &File{
		file,
		lock,
	}
}

func (fc *FileCache) ReadAndProcess() {
	log.Debug("readandprocess!!!!!!!!!!", fc.dir)

	fc.lock.Lock()
	defer fc.lock.Unlock()

	tmps, err := utils.ListDir(fc.dir, "che")
	if err != nil {
		log.Error("read offset tmp error:", err)
		return
	}

	if fc.file.file != nil {
		fc.file.file.Close()
	}

	for _, tmp := range tmps {
		log.Debug("aaaaaaaaaaaaaaaaaaaaaaaa:", tmp, fc.filenamePrefix)
		if strings.Index(tmp, fc.filenamePrefix) != 0 {
			log.Debug("filenameprefix is not exists:", tmp, fc.filenamePrefix)
			continue
		}
		//先更名

		dofile := fmt.Sprintf("%s/%s_doing.che", fc.dir, tmp[:len(tmp)-4])
		os.Rename(fc.dir+"/"+tmp, dofile)

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

		rel := fc.processHandle(dofile, []byte{})

		log.Debug("dddddddddddddddddddddddd")
		if rel {
			os.Remove(dofile)
		}
	}

	fc.openCache()
}

func (fc *FileCache) Append(data []byte) {
	fc.in <- data
}

func (fc *FileCache) _append(data []byte) {
	fc.file.Lock()
	if _, err := fc.file.file.Write(data); err != nil {
		log.Fatal(err)
	}
	fc.file.Unlock()
}
