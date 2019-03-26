package tempcache


type processFunc func(filename string, orginFilename string ,data []byte) bool

type ITempCache interface{
  Init(processHandle processFunc, intervalMinuteFlushData int)
  Append(data []byte)
  SyncCache()
  Quit()
}

