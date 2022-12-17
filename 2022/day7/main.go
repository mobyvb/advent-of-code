package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"mobyvb.com/advent/common"
	"mobyvb.com/advent/common/filesystem"
)

func main() {
	ld, err := common.OpenFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	// first, create the filesystem structure by navigating through the steps in the input
	rootDir := filesystem.NewDirectory("/", nil)
	currentDir := rootDir
	ld.EachF(func(currentLine string) {
		switch currentLine[:3] {
		case "$ c": // "$ cd <dirname>"
			dirname := strings.TrimPrefix(currentLine, "$ cd ")
			switch dirname {
			case "/":
				currentDir = rootDir
			case "..":
				currentDir = currentDir.Parent
			default:
				currentDir = currentDir.GetDirectory(dirname)
			}
		case "$ l": // "$ ls"
			// no-op
		case "dir": // "dir <dirname>" (output of ls)
			dirname := strings.TrimPrefix(currentLine, "dir ")
			currentDir.AddDirectory(dirname)
		default: // "<filesize> <filename>" (output of ls)
			parts := strings.Split(currentLine, " ")
			fileSize, err := strconv.Atoi(parts[0])
			if err != nil {
				panic(err)
			}
			fileName := parts[1]
			currentDir.AddFile(fileName, fileSize)
		}
	})

	// part 1: find all directories with total size of at most 100000, and sum them
	totalPart1Size := 0
	maxSize := 100000
	rootSize := rootDir.Size()
	if rootSize <= maxSize {
		totalPart1Size += rootSize
	}
	rootDir.Walk(func(d *filesystem.Directory) {
		dirSize := d.Size()
		if dirSize <= maxSize {
			totalPart1Size += dirSize
		}
	})
	fmt.Println(totalPart1Size)

	// part 2: find the smallest directory that is able to be deleted to free up enough space for the update
	totalSpace := 70000000
	neededSpace := 30000000
	spaceRemaining := totalSpace - rootSize
	smallestToDelete := rootSize
	rootDir.Walk(func(d *filesystem.Directory) {
		dirSize := d.Size()
		if spaceRemaining+dirSize > neededSpace &&
			dirSize < smallestToDelete {
			smallestToDelete = dirSize
		}
	})
	fmt.Println(smallestToDelete)
}
