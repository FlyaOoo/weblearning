package main

import (
	"bytes"
	"encoding/csv"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
)

type Post struct {
	Id      int
	Content string
	Author  string
}

// 使用Post的指针是为了能够映射到同一个post而不是同一个post的不同副本
var PostById map[int]*Post
var PostByAuthor map[string][]*Post

func storeWithMemory(post Post) {
	PostById[post.Id] = &post
	PostByAuthor[post.Author] = append(PostByAuthor[post.Author], &post)
}

func storeWithCSV(postSlice []Post) {
	csvFile, err := os.Create("./posts.csv")
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	for _, post := range postSlice {
		line := []string{strconv.Itoa(post.Id), post.Content, post.Author}
		if err := writer.Write(line); err != nil {
			panic(err)
		}
	}
	// 缓冲区可能存在的数据写入文件
	writer.Flush()
}

func storeWithByte(data interface{}, filename string) {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(data)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(filename, buffer.Bytes(), 0600)
	if err != nil {
		panic(err)
	}
}

func readByteFile(data interface{}, filename string) {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	buffer := bytes.NewBuffer(raw)
	dec := gob.NewDecoder(buffer)
	err = dec.Decode(data)
	if err != nil {
		panic(err)
	}
}

func main() {
	PostById = make(map[int]*Post)
	PostByAuthor = make(map[string][]*Post)

	posts := []Post{
		{Id: 1, Content: "Hello World", Author: "liyao"},
		{Id: 2, Content: "Web With Go", Author: "michal"},
		{Id: 3, Content: "Python Learning", Author: "guido"},
		{Id: 4, Content: "World Hello", Author: "liyao"},
	}

	// test for store with memory
	for _, post := range posts {
		storeWithMemory(post)
	}
	fmt.Println(PostById[1])
	fmt.Println(PostById[2])
	for _, post := range PostByAuthor["liyao"] {
		fmt.Println(post)
	}
	for _, post := range PostByAuthor["guido"] {
		fmt.Println(post)
	}

	// test for store with csv
	//fmt.Println(os.Getwd())
	storeWithCSV(posts)
	file, err := os.Open("posts.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
		}
		fmt.Println(record)
	}

	// test for store with bytes
	storeWithByte(posts, "posts")
	var postRead []Post
	readByteFile(&postRead, "posts")
	fmt.Println(postRead)
}
