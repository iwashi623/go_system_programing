package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

func dumpChumk(chunk io.Reader) {
	var length int32
	binary.Read(chunk, binary.BigEndian, &length)
	buffer := make([]byte, 4)
	chunk.Read(buffer)
	fmt.Printf("chunk '%v' (%d bytes)\n", string(buffer), length)
}

// PNG ファイルはバイナリフォーマット。先頭の8バイトがシグニチャ(固定のバイト列)となっている。
// 1チャンク32バイト
func readChunks(file *os.File) []io.Reader {
	var chunks []io.Reader

	// 最初の8バイトを飛ばす
	file.Seek(8, 0)
	var offset int64 = 8

	for {
		var length int32
		err := binary.Read(file, binary.BigEndian, &length)
		if err == io.EOF {
			break
		}
		chunks = append(chunks, io.NewSectionReader(file, offset, int64(length)+12))

		// offsetの値を次のチャンクの先頭に変更
		offset, _ = file.Seek(int64(length+8), 1)
	}

	return chunks
}

func main() {
	file, err := os.Open("PNG_transparency_demonstration_1.png")
	if err != nil {
		panic(err)
	}

	defer file.Close()
	chunks := readChunks(file)
	for _, chunk := range chunks {
		dumpChumk(chunk)
	}
}
