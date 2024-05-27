package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Ids struct {
	Num int64
	ID  string
}

func (product Ids) ToString() string {
	return fmt.Sprintf("name: %d\nID: %d", product.Num, product.ID)
}

func GetDB() (db *sql.DB, err error) {
	// dbDriver := "mysql"
	// dbUser := "root"
	// dbPass := "password"
	// dbName := "IDs"
	dbDriver := os.Getenv("DRIVER")
	dbUser := os.Getenv("USER")
	dbPass := os.Getenv("PASSWORD")
	dbName := os.Getenv("DBNAME")
	db, err = sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	return
}

func resondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)

}

type ProductModel struct {
	Db *sql.DB
}

func (ProductModel ProductModel) FindAllEmpty() (product []Ids, err error) {
	rows, err := ProductModel.Db.Query("select num, ID from id_produced where inserted=0")
	if err != nil {
		return nil, err
	} else {
		var products []Ids
		for rows.Next() {
			var num int64
			var id string
			err2 := rows.Scan(&num, &id)
			if err2 != nil {
				return nil, err2
			} else {
				product := Ids{
					Num: num,
					ID:  id,
				}
				products = append(products, product)
			}
		}
		return products, nil
	}
}

func FindAll(response http.ResponseWriter, request *http.Request) {
	db, err := GetDB()
	if err != nil {
		resondWithError(response, http.StatusBadRequest, err.Error())
	} else {
		productModel := ProductModel{
			Db: db,
		}
		products, err2 := productModel.FindAllEmpty()
		if err2 != nil {
			resondWithError(response, http.StatusBadRequest, err2.Error())
		} else {
			respondWithJson(response, http.StatusOK, products)
		}
	}
}

func UpdateStatusID(response http.ResponseWriter, request *http.Request) {
	fmt.Println("POST REQUEST INSERTION")
	vars := mux.Vars(request)
	db, err := GetDB()
	if err != nil {
		resondWithError(response, http.StatusBadRequest, err.Error())
	} else {
		productModel := ProductModel{
			Db: db,
		}
		err2 := productModel.UpdateStatusIDEmpty(vars["id"])
		if err2 != nil {
			resondWithError(response, http.StatusBadRequest, err2.Error())
		} else {
			respondWithJson(response, http.StatusOK, "")
		}
	}
}

func (ProductModel ProductModel) UpdateStatusIDEmpty(ID string) error {
	query := fmt.Sprintf("UPDATE id_produced SET inserted = 1 WHERE ID = %s", ID)
	fmt.Println(query)
	_, err := ProductModel.Db.Query(query)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (ProductModel ProductModel) GenerateIDsRequestEmpty(Qs string) error {
	query := fmt.Sprintf("call Insert_new_id(%s)", Qs)
	fmt.Println(query)
	_, err := ProductModel.Db.Query(query)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func GenerateIDsRequest(response http.ResponseWriter, request *http.Request) {
	type Quantita struct {
		Quantita   string `json:Quantita`
		ReqCompany string `json:Requester`
	}
	var prova Quantita
	err := json.NewDecoder(request.Body).Decode(&prova)
	//il body contiene l'organizzazione richiedente (e.g. Org1MSP) e la quantita'
	//di tronchi per cui si richiede la generazione
	fmt.Printf("\n")
	fmt.Println("AUDIT INSERTION")
	fmt.Printf("POST REQUEST FROM " + prova.ReqCompany + " for " + prova.Quantita + " IDs ")
	fmt.Printf("\n")
	db, err := GetDB()
	if err != nil {
		resondWithError(response, http.StatusBadRequest, err.Error())
	} else {
		productModel := ProductModel{
			Db: db,
		}
		err2 := productModel.GenerateIDsRequestEmpty(prova.Quantita)
		if err2 != nil {
			resondWithError(response, http.StatusBadRequest, err2.Error())
		} else {
			respondWithJson(response, http.StatusOK, "")
		}
	}
}

func main() {
	serverPort := os.Getenv("SERVERPORT")
	fmt.Printf("server listening on port %s ", serverPort)
	router := mux.NewRouter()
	router.HandleFunc("/api/product/audit", GenerateIDsRequest).Methods("POST")
	router.HandleFunc("/api/product/findall", FindAll).Methods("GET")
	router.HandleFunc("/api/product/inserted/{id}", UpdateStatusID).Methods("POST")
	err := http.ListenAndServe(":"+serverPort, router)

	if err != nil {
		fmt.Println(err)
	}
}
