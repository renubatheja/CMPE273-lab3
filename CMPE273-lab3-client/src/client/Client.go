package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"consistentHashing"
	"log"
	"strings"
	"strconv"
	"github.com/julienschmidt/httprouter"
)

// Variable Declarations 
// A ring representing the nodes on a circle hashed using Consistent Hashing Algorithm
var Ring *consistentHashing.HashRing
// An array of my server instances i.e. servers running on ports 3000, 3001 and 3002 
var MyRESTServers []string
// A map to store generated seed data. This data set would be shared to three server instances, using Consistent Hashing Algorithm
var DataMap = make(map[string]string)


//Main Function
func main() {
	fmt.Println("Starting the client")
	MyRESTServers = []string{"127.0.0.1:3000",
	                            "127.0.0.1:3001",
	                            "127.0.0.1:3002"}
	
	Ring = consistentHashing.New(MyRESTServers)

    mux := httprouter.New()

    //handler to service GET ALL request
    mux.GET("/keys", GetAllKeys)
    
    //handler to service GET request
    mux.GET("/keys/:id", GetKey)
    
    //handler to service PUT request
    mux.PUT("/keys/:id/:value", PutKey)
    
    server := http.Server{
            Addr:        "0.0.0.0:8080",
            Handler: mux,
    }
    server.ListenAndServe()
}


//Main function for Client, if Client is to be run as a Simple Standalone GO application 
func mainStandAlone() {
	MyRESTServers = []string{"127.0.0.1:3000",
	                            "127.0.0.1:3001",
	                            "127.0.0.1:3002"}
	
	Ring = consistentHashing.New(MyRESTServers)
	
	//Generate seed-data to be distributed across the servers
	GenerateSeedData()	
	
	//Distribute keys
	for key, value := range DataMap {
    	fmt.Println("Distributing/Saving this <key-valye> pair : <",key,"-",value,">")
    	PutOperation(key, value)
	}
	
	//Try to retrive few keys 
	fmt.Println("Retrieving FEW keys now!")
	GetOperation("2")
	GetOperation("4")
	GetOperation("6")

	//Try to retrive few keys 
	fmt.Println("Retrieving ALL keys now!")
	GetAllOperation("127.0.0.1:3000")
	GetAllOperation("127.0.0.1:3001")
	GetAllOperation("127.0.0.1:3002")

}


/*
* PutKey function - Function to service PUT REST request on Client
* This will internally send a PUT request to the Server for storing the key and value pair
*/
func PutKey(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
    // Grab id
    id := p.ByName("id")
    // Grab value
    value := p.ByName("value")
    
	
	server,_ := Ring.GetNode(id)
	fmt.Println("The <key-value> pair to be saved : <",id,"-",value,">")
	fmt.Println("The server for this key : ",server) //127.0.0.1:3001
	
	//make a corresponding PUT request to the server here
	url := "http://" + server + "/keys/" + id + "/" + value
	
	client := &http.Client{}
	request, err := http.NewRequest("PUT", url, strings.NewReader(""))
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("The Response is:", contents)
		rw.WriteHeader(200)	
	}		
	fmt.Println("Key saved successfully!")
	fmt.Println("------------------------------------------------------------")	
	
}


/*
* GetAllKey function - Function to service GET ALL request on Client
* This will internally send a GET ALL request to the Server for fetching all the key and value pairs
*/
func GetAllKeys(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	var server string
	results := ""
	for index := 0; index < len(MyRESTServers); index++ {
		server = MyRESTServers[index]
		fmt.Println("Getting all keys from this server : ",server) //127.0.0.1:3001
		
		//make a corresponding GET ALL request to the server here
		url := "http://" + server + "/keys"
	
	
		response, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		} else {
			defer response.Body.Close()
			contents, err := ioutil.ReadAll(response.Body)
			results = results + string(contents) + "\n "
			if err != nil {
				log.Fatal(err)
			}
			
		}
	}
	// Write content-type, statuscode, payload
	fmt.Println("The Response is:", results)
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(200)
	fmt.Fprintf(rw, "%s", results)
	fmt.Println("All Keys retrieved successfully!")
	fmt.Println("------------------------------------------------------------")	
}


/*
* GetKey function - Function to service GET REST request on Client
* This will internally send a GET request to the Server for fetching the key and value pair
*/
func GetKey(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	// Grab id
    id := p.ByName("id")
	
	server,_ := Ring.GetNode(id)
	fmt.Println("Getting key from this server : ",server) //127.0.0.1:3001
	
	//make a corresponding GET request to the server here
	url := "http://" + server + "/keys/" + id
	
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("The Response is : %s\n", string(contents));
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(200)
		fmt.Fprintf(rw, "%s", contents)		
	}

	fmt.Println("Key retrieved successfully!")
	fmt.Println("------------------------------------------------------------")	

}


//Function to generate seed data for storage
//Keys are : 1 to 26
//Values are : Unicode characters : a - z
func GenerateSeedData() {
	value := "abcdefghijklmnopqrstuvwxyz"
    // Split after an empty string to get all letters.
    result := strings.SplitAfter(value, "")
    var key string
    j := 1
    for i := range(result) {
		// Get letter and save it in the map
		letter := result[i]
		key = strconv.Itoa(j)
		DataMap[key] = letter
		j++ 
    }
}


/*
* PutOperation function - Function to support PUT operation and shard the data set into three server instances
*/
func PutOperation(key string, value string) {
    // Grab id
    id := key
	
	server,_ := Ring.GetNode(id)
	fmt.Println("The <key-value> pair to be saved : <",key,"-",value,">") 
	fmt.Println("The server for this key : ",server) 
	
	//make a corresponding PUT request to the server here
	url := "http://" + server + "/keys/" + id + "/" + value
	
	client := &http.Client{}
	request, err := http.NewRequest("PUT", url, strings.NewReader(""))
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	} else {
		defer response.Body.Close()
		_, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
	}		
	fmt.Println("Key saved successfully!")
	fmt.Println("------------------------------------------------------------")		
}


/*
* GetOperation function - Function to support GET operation and retrieve keys from sharded the data set into three server instances
*/
func GetOperation(key string) {
	// Grab id
    id := key
	
	server,_ := Ring.GetNode(id)
	fmt.Println("Retrieving key from this server : ",server) //127.0.0.1:3001
	
	//make a corresponding GET request to the server here
	url := "http://" + server + "/keys/" + id
	
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("The Response is : %s\n", string(contents));
	}

	fmt.Println("Key retrieved successfully!")
	fmt.Println("------------------------------------------------------------")
	
}

/*
* GetAllOperation function - Function to support GET ALL operation and retrieve ALL keys from any particular server instances
*/
func GetAllOperation(server string) {
	fmt.Println("Retrieving ALL keys from this server : ",server) 
	
	//make a corresponding GET request to the server here
	url := "http://" + server + "/keys"
	
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("The Response is : %s\n", string(contents));
	}

	fmt.Println("Keys retrieved successfully!")
	fmt.Println("------------------------------------------------------------")

}