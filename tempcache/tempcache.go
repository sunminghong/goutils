package tempcache


type ITempCache interface{
  //Init(processDelay time.Duration, processHandle func([]byte)bool,  args... interface{})
  //Write(data []byte)
  Append(data []byte)
  ReadAndProcess()
}

