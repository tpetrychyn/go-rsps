package util

import (
	"bytes"
	"compress/bzip2"
	"fmt"
	"io"
	"log"
	"os"
)

type CacheLoader struct{}

const INDEX_SIZE = 6
const SECTOR_HEAD_SIZE = 8
const SECTOR_SIZE = 520

var rootCache *os.File

func LoadCache() *CacheLoader {
	file, err := os.Open("./definitions/cache/main_file_cache.dat")
	if err != nil {
		fmt.Print(err)
	}
	//rootCache := bytes.NewBuffer(nil)
	rootCache = file
	io.Copy(rootCache, file)
	//file.Close()

	file, err = os.Open("./definitions/cache/main_file_cache.idx1")
	if err != nil {
		fmt.Print(err)
	}
	idxCache := bytes.NewBuffer(nil)
	io.Copy(idxCache, file)
	//file.Close()

	indexId := 1
	position := int64(indexId * INDEX_SIZE)

	buffer := make([]byte, INDEX_SIZE)
	file.Seek(position, 0)
	file.Read(buffer)

	buf := bytes.NewBuffer(buffer)
	length := getMedium(buf)
	id := getMedium(buf)

	log.Printf("%+v %+v", length, id)

	position = int64(id * SECTOR_HEAD_SIZE)

	data := make([]byte, length)
	next := id
	offset := 0

	for chunk := 0; offset < length; chunk++ {
		read := length - offset

		//readsector
		sector := readSector(next, data, offset, read)

		next = sector.NextIndexId
		offset += read
	}

	log.Printf("data %+v", data)

	buf = bytes.NewBuffer(data)
	length = getMedium(buf)
	compressedLength := getMedium(buf)

	log.Printf("length %d, compressedlength %d", length, compressedLength)
	h := []byte("h")[0]
	one := []byte("1")[0]
	bzipData := []byte{h, one}
	bzipData = append(bzipData, data...)
	bReader := bytes.NewBuffer(bzipData)
	bzipReader := bzip2.NewReader(bReader)
	total := []byte{0,0}
	bzipReader.Read(total)

	return &CacheLoader{}
}

type CacheSector struct {
	IndexId     int
	Chunk       int
	NextIndexId int
	CacheId     int
}

func readSector(sectorId int, data []byte, offset, length int) *CacheSector {
	position := int64(sectorId * SECTOR_SIZE)

	buffer := make([]byte, length*SECTOR_HEAD_SIZE)

	rootCache.Seek(position, 0)
	rootCache.Read(buffer)
	buf := bytes.NewBuffer(buffer)
	log.Printf("%+v", buf)

	return decodeSector(buf, data, offset, length)
}

func decodeSector(buffer *bytes.Buffer, data []byte, offset, length int) *CacheSector {
	id := getShort(buffer)
	chunk := getShort(buffer)
	nextIndexId := getMedium(buffer)
	cacheId, _ := buffer.ReadByte()
	log.Printf("id %d chunk %d nextIndexId %d cacheId %d", id, chunk, nextIndexId, cacheId)
	buffer.Read(data)
	return &CacheSector{
		IndexId:     id,
		Chunk:       chunk,
		NextIndexId: nextIndexId,
		CacheId:     int(cacheId),
	}
}

func getShort(buf *bytes.Buffer) int {
	hi, _ := buf.ReadByte()
	lo, _ := buf.ReadByte()
	return int((hi << 8) | lo)
}
func getMedium(buf *bytes.Buffer) int {
	hi, _ := buf.ReadByte()
	lo, _ := buf.ReadByte()
	l, _ := buf.ReadByte()
	return int((hi<<8|lo)<<8 | (l & 0xFF))
}
