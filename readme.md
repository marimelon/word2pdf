# word2pdf
Application to convert word files to PDF using Microsoft Office.  
Use multipart/form-data to send files.

Windows only.

## Run

```
go run main.go word.go
```

## Usage

Access `http://127.0.0.1:8000` from your browser.

or 

```
# From curl
curl -F file=@filename http://127.0.0.1:8000 -o out.pdf
```