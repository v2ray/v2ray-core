package io

import(
  "errors"
)

const (
  SizeSmall = 16
  SizeMedium = 128
  SizeLarge = 512
)

var (
  ErrorNoChannel = errors.New("No suitable channels found.")
)

type VBufferSet struct {
  small chan []byte
  medium chan []byte
  large chan []byte
}

func NewVBufferSet() *VBufferSet {
  bSet := new(VBufferSet)
  bSet.small = make(chan []byte, 128)
  bSet.medium = make(chan []byte, 128)
  bSet.large = make(chan []byte, 128)
  return bSet
}

func (bSet *VBufferSet) detectBucket(size int, strict bool) (chan []byte, error) {
  if strict {
    if size == SizeSmall {
      return bSet.small, nil
    } else if size == SizeMedium {
      return bSet.medium, nil
    } else if size == SizeLarge {
      return bSet.large, nil
    }
  } else {
    if size <= SizeSmall {
      return bSet.small, nil
    } else if size <= SizeMedium {
      return bSet.medium, nil
    } else if size <= SizeLarge {
      return bSet.large, nil
    }
  }
  return nil, ErrorNoChannel
}

func (bSet *VBufferSet) FetchBuffer(minSize int) []byte {
  var buffer []byte
  byteChan, err := bSet.detectBucket(minSize, false)
  if err != nil {
    return make([]byte, minSize)
  }
  select {
    case buffer = <- byteChan:
    default:
    buffer = make([]byte, minSize)
  }
  return buffer
}

func (bSet *VBufferSet) ReturnBuffer(buffer []byte) {
  byteChan, err := bSet.detectBucket(len(buffer), true)
  if err != nil {
    return
  }
  select {
    case byteChan <- buffer:
    default:
  }
}
