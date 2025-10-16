package main


import (
	"net/http"
	"log"
	"fmt"
	"strconv"
)


type database map[string]int

func main() {
	mux := http.NewServeMux()
	db := database{}

	mux.HandleFunc("/list", db.list) // mux диспетчеризирует запросы
	mux.HandleFunc("/find", db.find)
	mux.HandleFunc("/add", db.add)

	if err := http.ListenAndServe("localhost:8080", mux); err != nil { // Второй аргумент имеет тип интерфейса http.Handler с методом ServeHTTP
		log.Fatalf("ERROR server connetion: %v", err)
	}
}

// Каждый обработчик по интерфейсу http.Handler должен иметь тип функции с параметрами запроса и ответа
func (db database) list(w http.ResponseWriter, r *http.Request) {
	for k, v := range db {
		fmt.Fprintf(w, "%s\t%d\n", k, v)
	}
}

func (db database) find(w http.ResponseWriter, r *http.Request) {
	item := r.URL.Query().Get("item")
	price, ok := db[item]
	if !ok {
		fmt.Fprintf(w, "Item does not exist : %s", item)
		return
	}
	fmt.Fprintf(w, "price : %d\n", price)
}

func (db database) add(w http.ResponseWriter, r *http.Request) {
	item  := r.URL.Query().Get("item")
	price := r.URL.Query().Get("price")

	priceint, err := strconv.Atoi(price) // нэйминг страдает, да 
	if err != nil {
		fmt.Fprintf(w, "Impossible price : %s", price)
		return
	}
	db[item] = priceint
	fmt.Fprint(w, "Added")
}