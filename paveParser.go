package goarstream

import (
	"bufio"
	"encoding/binary"
	"errors"
)

const (
	skippedFrameBytes = 14
	frameSig          = "PaVE"
	frameSize         = 68
)

type frame struct {
	Signature            string `json:"signature"`
	Version              uint8  `json:"version"`
	VideoCodec           uint8  `json:"video_codec"`
	HeaderSize           uint16 `json:"header_size"`
	PayloadSize          uint32 `json:"payload_size"`
	EncodedStreamWidth   uint16 `json:"encoded_stream_width"`
	EncodedStreamHeight  uint16 `json:"encoded_stream_height"`
	DisplayWidth         uint16 `json:"display_width"`
	DisplayHeight        uint16 `json:"display_height"`
	FrameNumber          uint32 `json:"frame_number"`
	Timestamp            uint32 `json:"time_stamp"`
	TotalChunks          uint8  `json:"total_chunks"`
	ChunkIndex           uint8  `json:"chunk_index"`
	FrameType            uint8  `json:"frame_type"`
	Control              uint8  `json:"control"`
	StreamBytePositionLW uint32 `json:"stream_byte_positionLW"`
	StreamBytePositionUW uint32 `json:"StreamBytePositionUW"`
	StreamID             uint16 `json:"StreamID"`
	TotalSlices          uint8  `json:"TotalSlices"`
	SliceIndex           uint8  `json:"SliceIndex"`
	Header1Size          uint8  `json:"Header1Size"`
	Header2Size          uint8  `json:"Header2Size"`
	Reserved2            uint8  `json:"Reserved2"`
	AdvertisedSize       uint32 `json:"AdvertisedSize"`
	Payload              []byte `json:"payload"`
}

func (f *frame) parse(reader *bufio.Reader) error {
	var err error
	endianess := binary.LittleEndian

	f.Signature, _ = f.getSig(reader)
	if !isPaVE(f.Signature) {
		return errors.New("Invalid signature: " + f.Signature)
	}
	binary.Read(reader, endianess, &f.Version)
	binary.Read(reader, endianess, &f.VideoCodec)
	binary.Read(reader, endianess, &f.HeaderSize)
	binary.Read(reader, endianess, &f.PayloadSize)
	binary.Read(reader, endianess, &f.EncodedStreamWidth)
	binary.Read(reader, endianess, &f.EncodedStreamHeight)
	binary.Read(reader, endianess, &f.DisplayWidth)
	binary.Read(reader, endianess, &f.DisplayHeight)
	binary.Read(reader, endianess, &f.FrameNumber)
	binary.Read(reader, endianess, &f.Timestamp)
	binary.Read(reader, endianess, &f.TotalChunks)
	binary.Read(reader, endianess, &f.ChunkIndex)
	binary.Read(reader, endianess, &f.FrameType)
	binary.Read(reader, endianess, &f.Control)
	binary.Read(reader, endianess, &f.StreamBytePositionLW)
	binary.Read(reader, endianess, &f.StreamBytePositionUW)
	binary.Read(reader, endianess, &f.StreamID)
	binary.Read(reader, endianess, &f.TotalSlices)
	binary.Read(reader, endianess, &f.SliceIndex)
	binary.Read(reader, endianess, &f.Header1Size)
	binary.Read(reader, endianess, &f.Header2Size)
	err = f.skip(reader, 2)
	binary.Read(reader, endianess, &f.AdvertisedSize)
	err = f.skip(reader, 12)

	// stupid kludge for https://projects.ardrone.org/issues/show/159
	err = f.skip(reader, int(f.HeaderSize-64))

	f.Payload = make([]byte, f.PayloadSize)
	reader.Read(f.Payload)
	return err
}

func (f *frame) skip(reader *bufio.Reader, n int) error {
	var err error
	for i := 0; i < n; i++ {
		_, err = reader.ReadByte()
	}
	return err
}

func (f *frame) getSig(reader *bufio.Reader) (string, error) {
	buff := make([]byte, 4)
	_, err := reader.Read(buff)
	if err != nil {
		return "", err
	}
	return string(buff), nil
}

func isPaVE(s string) bool {
	return s == frameSig
}
