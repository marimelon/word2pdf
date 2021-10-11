package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func fileIsExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func export(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed) // 405
		w.Write([]byte("POST Only"))
		return
	}

	file, fileInfo, err := r.FormFile("file")

	if err != nil {
		log.Println("ファイルアップロードを確認できませんでした。")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer file.Close()

	ext := filepath.Ext(fileInfo.Filename)
	saveFile, err := ioutil.TempFile("", "w2p*"+ext)
	if err != nil {
		log.Println("サーバ側でファイル確保できませんでした。")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer os.Remove(saveFile.Name())
	defer saveFile.Close()

	_, err = io.Copy(saveFile, file)
	if err != nil {
		log.Println(err)
		log.Println("アップロードファイルの書き込みに失敗しました。")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	outDir, _ := ioutil.TempDir("", "w2p")
	if err != nil {
		log.Println(err)
		log.Println("一時ディレクトリの作成に失敗しました。")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(outDir)

	word := new(Word)
	log.Println("input file: " + saveFile.Name())
	log.Println("output dir: " + outDir)

	//PDF変換
	outFilePath, err := word.Export(saveFile.Name(), outDir)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("output file: " + outFilePath)

	outFile, err := ioutil.ReadFile(outFilePath)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Length", strconv.Itoa(len(outFile)))
	w.Header().Set("Content-Disposition", `attachment; filename="`+filepath.Base(outFilePath)+`"`)
	w.Write(outFile)
}

func root(w http.ResponseWriter, r *http.Request) {
	html := `
	<!DOCTYPE html PUBLIC "-//W3C//DTD HTML 3.2 Final//EN">
    <html>
    <body>
    <form ENCTYPE="multipart/form-data" method="post" action="/upload">
    <input name="file" type="file"/>
    <input type="submit" value="upload"/>
    </form>
    </body>
    </html>
	`
	fmt.Fprintf(w, html)
}

func main() {
	port := "8000"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}
	http.HandleFunc("/", root)
	http.HandleFunc("/upload", export)
	log.Println("Server is listening on port " + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
