package receipt

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/andersonaskm/go_webservices/cors"
)

const receiptPath = "receipts"

func SetupRoutes(apiBasePath string) {
	receiptHandler := http.HandlerFunc(handleReceipts)
	downloadHandler := http.HandlerFunc(handleDownload)

	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, receiptPath), cors.Middleware(receiptHandler))
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, receiptPath), cors.Middleware(downloadHandler))
}

func handleDownload(w http.ResponseWriter, r *http.Request) {
	// verifica se o segmento segue o pattern "/receipts/receiptFileName"
	urlPathSegments := strings.Split(r.URL.Path, fmt.Sprintf("%s/", receiptPath))
	if len(urlPathSegments[1:]) > 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// obtem o nome do arquivo
	fileName := urlPathSegments[1:][0]

	// abre o arquivo
	file, errFileOpen := os.Open(filepath.Join(ReceiptDirectory, fileName))
	if errFileOpen != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer file.Close()

	// identifica o tipo de conteudo
	fHeader := make([]byte, 512) // cria um array de bytes com 512 posições
	file.Read(fHeader)

	fContentType := http.DetectContentType(fHeader)
	stat, errStat := file.Stat()
	if errStat != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// tamanho do arquivo
	fSize := strconv.FormatInt(stat.Size(), 10)

	// define os cabeçalhos de saída do arquivo
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", fContentType)
	w.Header().Set("Content-Length", fSize)

	// redireciona o ponteiro para inicio do arquivo
	file.Seek(0, 0)

	io.Copy(w, file)
}

func handleReceipts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// obter arquivos
	case http.MethodGet:
		receiptList, errGetReceipts := GetReceipts()
		if errGetReceipts != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		j, errJson := json.Marshal(receiptList)
		if errJson != nil {
			log.Fatal(errJson)
		}

		_, err := w.Write(j)
		if err != nil {
			log.Fatal(err)
		}
	// cadastrar arquivo
	case http.MethodPost:
		r.ParseMultipartForm(5 << 20) // 5mb
		file, handler, errFile := r.FormFile("receipt")
		if errFile != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		defer file.Close()

		f, errOpenFile := os.OpenFile(filepath.Join(ReceiptDirectory, handler.Filename), os.O_WRONLY|os.O_CREATE, 0666)
		if errOpenFile != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		defer f.Close()
		io.Copy(f, file)
		w.WriteHeader(http.StatusCreated)
	// CORS
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
