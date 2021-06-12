GoLang - Creating Web Services with Go

1 - Handling HTTP Requests

    Basic Handlers 

        http.Handle / http.HttpHandleFunc
        ServeMux
        http.ListenAndServe(port, nil)
        http.ListenAndServeTLS(addr, certFile, keyFile, handler) error

    Json

        json.Marshal - func Marshal(v interface{}) ([]byte, error)
        json.Unmarshal - func Unmarshal(data byte[], v interface{}) error

    Request
    URL Path
    Middleware
    CORS

2 - Persistent Data
3 - Websockets
4 - Templating