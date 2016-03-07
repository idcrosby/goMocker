package main

import (
     "encoding/json"
     "fmt"
     "math/rand"
     "net/http"
     "os"
     "time"
)

var routes map[string]Resource
var swagger SwaggerDoc

func main() {

     // Load json schema
     file, err := os.Open("swaggerSample2.json")
     if err != nil {
          fmt.Println("Cannot find schema file. Exiting...")
          return
     }
     decoder := json.NewDecoder(file)
     swagger = SwaggerDoc{}
     err = decoder.Decode(&swagger)
     if err != nil {
          panic(err)
     } else {
          fmt.Printf("Read in swagger conf:\n%v\n", swagger)
     }

     // Determine path
     m := swagger.Paths.(map[string]interface{})
     for key,value := range m {
          fmt.Printf("Found path for: %s\n", key)
          m2 := value.(map[string]interface{})
          for k2, v2 := range m2 {
               path, ok := v2.(SwagPath)
               if !ok {
                    fmt.Printf("Unable to cast method: %s to SwaggerPath\n%v\n", k2, v2)
               }
               fmt.Printf("listening for: %s\n", path.summary)
          }

          // decoder = json.NewDecoder(value)
          // SwaggerPath path = SwaggerPath
          // err = decoder.Decode(&path)
          // if err != nil {
          //      panic(err)
          // }
     }

     http.HandleFunc(swagger.BasePath, rootHandler)
     http.ListenAndServe(":8080", nil)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
     resource, ok := routes[r.URL.Path]

     if !ok {
          w.WriteHeader(404)
          return
     }

     valid := false
     for _, m := range resource.methods {
          if m == r.Method {
               valid = true
               break
          }
     }

     if !valid {
          w.WriteHeader(404)
          return
     }

     w.WriteHeader(200)
     body, err := json.Marshal(generateResponse(resource.response))
     if err != nil {
          panic(err)
     }
     w.Write(body)

}

func generateResponse(r interface{}) map[string]interface{} {
     rand.Seed(time.Now().UTC().UnixNano())

     // Access the data's underlying interface
     m := r.(map[string]interface{})

     var result map[string]interface{}
     result = make(map[string]interface{})

     for key,value := range m {
          var valueType string
          // determine type of value
          switch value.(type) {
               case int:
                    // valueType = "number"
                    result[key] = rand.Int()
               case float64:
                    result[key] = rand.Float64()
               case string:
                    result[key] = "random string...."
                    // valueType = "string"
               case bool:
                    // valueType = "boolean"
                    temp := rand.Intn(1)
                    result[key] = (temp == 1)
               case nil:
                    // valueType = "null"
                    result[key] = nil
               case map[string]interface{}:
                    // valueType = "object"
                    result[key] = generateResponse(value)
               case []interface{}:
                    // valueType = "array"
                    // TODO implement
                    result[key] = nil
               default:
                    // valueType = "unknown"
                    fmt.Printf("Unknown type for key: %s\n", key)
          }
          result[key] = valueType
     }
     return result
}

type Resource struct {
     path string
     methods []string
     response interface{}
}

/**
  Swagger Model
  */

type SwaggerDoc struct {
     Swagger string
     Info struct {
          Version string
          Title string
          Description string
          TermsOfService string
          Contact struct {
               Name string
               Url string
          }
          License struct {
               Name string
               Url string
          }
     }
     Host string
     BasePath string
     Schemas []string
     Consumes []string
     Produces []string
     Paths interface{}
     Definitions interface{}
}

type SwaggerPath struct {
     Method string           `json:"-"`
     Description string
     OperationId string
     Produces []string
     Parameters []struct {
          Name string
          In string
          Description string
          Required bool
          PType string         `json:"type"`
          Items interface{}
          CollectionFormat string
          Format string
     }
}

type SwagPath struct {
     tags []string
     summary string
}