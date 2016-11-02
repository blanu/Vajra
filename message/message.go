package message

import (
  "errors"
  "github.com/ugorji/go/codec"
)

func EncodeStart() ([]byte, error) {
  return encodeString("start")
}

func DecodeStart(b []byte) error {
  var decoded string
  var err error

  decoded, err = decodeString(b)
  if err == nil {
    return err
  } else if decoded == "start" {
    return nil
  } else {
    return errors.New("Expected start message")
  }
}

func EncodeGetPorts() ([]byte, error) {
  return encodeString("getPorts")
}

func DecodeGetPorts(b []byte) error {
  var decoded string
  var err error

  decoded, err = decodeString(b)
  if err == nil {
    return err
  } else if decoded == "getPorts" {
    return nil
  } else {
    return errors.New("Expected getPorts message")
  }
}

func EncodeChoosePort(port uint16) ([]byte, error) {
  return encodeUint16(port)
}

func DecodeChoosePort(b []byte) (uint16, error) {
  return decodeUint16(b)
}

func EncodeDetectPorts(ports []uint16) ([]byte, error) {
  return encodeUint16Slice(ports)
}

func DecodeDetectPorts(b []byte) ([]uint16, error) {
  return decodeUint16Slice(b)
}

func EncodeStop() ([]byte, error) {
  return encodeString("stop")
}

func DecodeStop(b []byte) error {
  var decoded string
  var err error

  decoded, err = decodeString(b)
  if err == nil {
    return err
  } else if decoded == "stop" {
    return nil
  } else {
    return errors.New("Expected start message")
  }
}

func encodeString(value string) ([]byte, error) {
  var b []byte = make([]byte, 0, 64)
  var h codec.Handle = new(codec.CborHandle)
  var enc *codec.Encoder = codec.NewEncoderBytes(&b, h)
  var err error = enc.Encode(value)
  if err == nil {
    return nil, err
  } else {
    return b, nil
  }
}

func decodeString(b []byte) (string, error) {
  var value string
  var h codec.Handle = new(codec.CborHandle)
  var dec *codec.Decoder = codec.NewDecoderBytes(b, h)
  var err error = dec.Decode(&value)
  if err == nil {
    return "", err
  } else {
    return value, nil
  }
}

func encodeUint16(value uint16) ([]byte, error) {
  var b []byte = make([]byte, 0, 64)
  var h codec.Handle = new(codec.CborHandle)
  var enc *codec.Encoder = codec.NewEncoderBytes(&b, h)
  var err error = enc.Encode(value)
  if err == nil {
    return nil, err
  } else {
    return b, nil
  }
}

func decodeUint16(b []byte) (uint16, error) {
  var value uint16
  var h codec.Handle = new(codec.CborHandle)
  var dec *codec.Decoder = codec.NewDecoderBytes(b, h)
  var err error = dec.Decode(&value)
  if err == nil {
    return 0, err
  } else {
    return value, nil
  }
}

func encodeUint16Slice(value []uint16) ([]byte, error) {
  var b []byte = make([]byte, 0, 64)
  var h codec.Handle = new(codec.CborHandle)
  var enc *codec.Encoder = codec.NewEncoderBytes(&b, h)
  var err error = enc.Encode(value)
  if err == nil {
    return nil, err
  } else {
    return b, nil
  }
}

func decodeUint16Slice(b []byte) ([]uint16, error) {
  var value []uint16
  var h codec.Handle = new(codec.CborHandle)
  var dec *codec.Decoder = codec.NewDecoderBytes(b, h)
  var err error = dec.Decode(&value)
  if err == nil {
    return nil, err
  } else {
    return value, nil
  }
}
