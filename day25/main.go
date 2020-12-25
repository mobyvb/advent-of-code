package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	pubKey1, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}
	rfid1 := &RFID{
		pubKey: pubKey1,
	}
	rfid1.SetLoopSize()
	pubKey2, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}
	rfid2 := &RFID{
		pubKey: pubKey2,
	}
	rfid2.SetLoopSize()

	part1(rfid1, rfid2)
	part2(rfid1, rfid2)
}

type RFID struct {
	pubKey   int
	loopSize int
}

func (r *RFID) SetLoopSize() {
	value := 1
	subject := 7
	loopSize := 0
	for value != r.pubKey {
		value = transformOnce(value, subject)
		loopSize++
	}
	r.loopSize = loopSize
}

// loop size must be set
func (r *RFID) GetEncryptionKey(r2 *RFID) int {
	value := 1
	for i := 0; i < r.loopSize; i++ {
		value = transformOnce(value, r2.pubKey)
	}
	return value

}

func transformOnce(value, subject int) int {
	value *= subject
	value %= 20201227
	return value
}

func part1(rfid1, rfid2 *RFID) {
	fmt.Println(rfid1.GetEncryptionKey(rfid2))
}

func part2(rfid1, rfid2 *RFID) {
	// no part 2!
}
