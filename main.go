package main

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
    "net/http"
	"os"
)


type orders struct {
	Order []order `json:"order"`
}
type order struct {
	Id string     `json:"id"`
	OrderItems []string `json:"orderItems"`
	Pincode string `json:"pincode"`
	Status string `json:"status"`
}


func readJSONFile(filename string)(data orders ,err error){
	jsonFile, err := os.Open(filename)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Successfully Opened File")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	fmt.Println(string(byteValue))

	err = json.Unmarshal([]byte(byteValue), &data)
	if err != nil {
		fmt.Println(err)
	}

	return
}

func writeToFile(data orders , filename string)(err error){

	dataBytes, err := json.Marshal(data)
	if err != nil {
		logrus.Error(err)
	}

	err = ioutil.WriteFile(filename, dataBytes, 0644)
	if err != nil {
		logrus.Error(err)
	}
	return
}


func getAllOrders()(orders  orders,err error){
	// fetch the orders present in the order.json file
	orders,err = readJSONFile("./orders.json")
	return
}

func addNewOrder(newOrder order )(err error ){

	orders,err := getAllOrders()
	if err != nil {
		fmt.Println("ERROR : Unable to get orders",err)
		return
	}

	orders.Order = append(orders.Order,newOrder)

	err = writeToFile(orders,"./orders.json")
	if err != nil {
		fmt.Println("ERROR : Error while adding a new order",err)
	}
    return
}

func changeOrderStatus(details map[string]string)(err error){
	id := details["id"]
	status := details["status"]
	var modifiedOrders orders
	orders,err := getAllOrders()
	if err != nil {
		fmt.Println("ERROR : Unable to get orders",err)
		return
	}

	for _,val := range orders.Order {
		if val.Id == id{
			val.Status = status
		}
		modifiedOrders.Order = append(modifiedOrders.Order,val)
	}
	err = writeToFile(modifiedOrders,"./orders.json")
	if err != nil {
		fmt.Println("ERROR : Error while adding a new order",err)
	}
	return
}


func handleOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		w.WriteHeader(http.StatusOK)
		orders,err := getAllOrders()
		if err != nil {
			fmt.Println("ERROR : Unable to get orders",err)
			return
		}
		b, err := json.Marshal(orders)
		if err != nil {
			panic(err)
		}
		w.Write([]byte(b))
	case "POST":
		var newOrder order
		decoder := json.NewDecoder(r.Body)
		decoder.Decode(&newOrder)

		err := addNewOrder(newOrder)
		if err != nil {
			fmt.Println("ERROR : Error while adding the new order.",err)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"message": "post called"}`))
	case "PUT":
		var details map[string]string
		decoder := json.NewDecoder(r.Body)
		decoder.Decode(&details)

		err := changeOrderStatus(details)
		if err != nil {
			fmt.Println("ERROR : Error while changing the order status.",err)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(`{"message": "put called"}`))
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "not found"}`))
	}
}

func main() {
	http.HandleFunc("/order", handleOrders)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
