package main

import (
    "fmt"
    "github.com/julienschmidt/httprouter"
    "net/http"
    "encoding/json"    
)

/*
* Request and Response structs
*/

type KeyResponse struct {
	Key string `json:"key"`
	Value string `json:"value"` 
}

var KeyMap = make(map[string]string)

/*
* PutKey function - Function to service PUT REST request
*/

func PutKey(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {

    // Grab id
    id := p.ByName("id")
    // Grab value
    value := p.ByName("value")

	KeyMap[id] = value 
	rw.WriteHeader(200)
	fmt.Println("Saving this <key-valye> pair : <",id,"-",value,">")
	fmt.Println("Key saved successfully!")
	fmt.Println("------------------------------------------------------------")	
	
}

/*
* GetAllKeys function - Function to service GET ALL REST request
*/

func GetAllKeys(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
    // Stub Trips
    var Keys = make([]KeyResponse, len(KeyMap)) //make(map[string]string)
    
    // Fetch keys
    index := 0
    for key, value := range KeyMap {
    	fmt.Println("Fetched <Key-value> pair : <", key, "-", value,">")
    	Keys[index].Key = key
    	Keys[index].Value = value
    	index++
	}
    
	results := "["
	for index := 0; index < len(Keys); index++ {
	    // Marshal provided interface into JSON structure
	    KeyJson, _ := json.Marshal(Keys[index])
		results = results + string(KeyJson) + ",\n "
	}
	results = results + "]"
	// Write content-type, statuscode, payload
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(200)
	fmt.Fprintf(rw, "%s", results)
	fmt.Println("All Keys retrieved successfully!")
	fmt.Println("------------------------------------------------------------")	
}

/*
* GetKey function - Function to service GET REST request
*/

func GetKey(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
    // Stub Trips
    var Key KeyResponse

    // Grab id
    id := p.ByName("id")
    
    // Fetch key
    for key, value := range KeyMap {
    	if(key == id) {
    		Key.Key = key
    		Key.Value = value
    		fmt.Println("Fetched <Key-value> pair : <", key, "-", value,">")
    		break
    	}
	}
    
	// Marshal provided interface into JSON structure
	KeyJson, _ := json.Marshal(Key)
	// Write content-type, statuscode, payload
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(200)
	fmt.Fprintf(rw, "%s", KeyJson)
	fmt.Println("Key retrieved successfully!")
	fmt.Println("------------------------------------------------------------")	
}


/*
* Main function - Setting up of httprouter and REST handlers
*/

func main() {
    mux := httprouter.New()

    //handler to service GET ALL request
    mux.GET("/keys", GetAllKeys)
    
    //handler to service GET request
    mux.GET("/keys/:id", GetKey)
    
    //handler to service PUT request
    mux.PUT("/keys/:id/:value", PutKey)
    
    server := http.Server{
            Addr:        "0.0.0.0:3000",
            Handler: mux,
    }
    server.ListenAndServe()
}