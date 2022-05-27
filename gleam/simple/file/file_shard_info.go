package file

import (
	"bytes"
	"dailylib/gleam/simple"
	"encoding/gob"
	"fmt"
	"github.com/chrislusf/gleam/filesystem"
	"io"
	"log"
	"os"
)

type FileShardInfo struct {
	Config    map[string]string
	FileName  string
	FileType  string
	HasHeader bool
	Fields    []string
}

var (
	registeredMapperReadShard = simple.RegisterMapper(readShard)
)

func init() {
	gob.Register(FileShardInfo{})
}

func encodeShardInfo(shardInfo *FileShardInfo) []byte {
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	if err := enc.Encode(shardInfo); err != nil {
		log.Fatal("encode shard info:", err)
	}
	return network.Bytes()
}

func readShard(row []interface{}) error {
	encodedShardInfo := row[0].([]byte)
	return decodeShardInfo(encodedShardInfo).ReadSplit()
}

func (ds *FileShardInfo) ReadSplit() error {

	// println("opening file", ds.FileName)
	fr, err := filesystem.Open(ds.FileName)
	if err != nil {
		return fmt.Errorf("Failed to open file %s: %v", ds.FileName, err)
	}
	defer fr.Close()

	reader, err := ds.NewReader(fr)
	if err != nil {
		return fmt.Errorf("Failed to read file %s: %v", ds.FileName, err)
	}
	if ds.HasHeader {
		reader.ReadHeader()
	}

	for {
		row, err := reader.Read()
		if err != nil {
			if err != io.EOF {
				log.Printf("ds.ReadSplit() Failed to read from file %s: %v", ds.FileName, err)
			}
			break
		}
		row.WriteTo(os.Stdout)
	}

	return err
}

func decodeShardInfo(encodedShardInfo []byte) *FileShardInfo {
	network := bytes.NewBuffer(encodedShardInfo)
	dec := gob.NewDecoder(network)
	var p FileShardInfo
	if err := dec.Decode(&p); err != nil {
		log.Fatal("decode shard info", err)
	}
	return &p
}