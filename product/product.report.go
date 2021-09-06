package product

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"path"
	"text/template"
	"time"
)

type ProductReportFilter struct {
	NameFilter         string `json:"productName"`
	ManufacturerFilter string `json:"manufacturer"`
	SKUFilter          string `json:"sku"`
}

const reportTemplateName = "report.gotmpl"

func productReportHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var productFilter ProductReportFilter

		// decodifica o conteudo do body em json
		errDecoder := json.NewDecoder(r.Body).Decode(&productFilter)
		if errDecoder != nil {
			log.Println(errDecoder)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		products, errSearchForProductData := SearchForProductData(productFilter)
		if errSearchForProductData != nil {
			log.Println(errSearchForProductData)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// define o template a partir de um arquivo
		// e adiciona uma função chamada "mod"
		t := template.New(reportTemplateName).Funcs(template.FuncMap{"mod": func(i, x int) bool { return i%x == 0 }})
		t, errParseFile := t.ParseFiles(path.Join("templates", reportTemplateName))
		if errParseFile != nil {
			log.Println(errParseFile)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// preenche o template com 1 produto
		var tmpl bytes.Buffer
		if len(products) > 0 {
			errTemplateExecute := t.Execute(&tmpl, products)
			if errTemplateExecute != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		//define um reader para retornar no response
		rdr := bytes.NewReader(tmpl.Bytes())
		w.Header().Set("Content-Disposition", "Attachment")
		http.ServeContent(w, r, "report.html", time.Now(), rdr)

	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
