package temp_sensor_app

import (
  "fmt"
  "io/ioutil"
  "net/http"
  
  "code.google.com/p/goprotobuf/proto"
  
  "appengine"
  
  "htu21df"
)

func init() {
  http.HandleFunc("/", root)
  http.HandleFunc("/upload", upload)
}

func root(w http.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, "Temp and humidity.")
}

func checkErr(w http.ResponseWriter, c appengine.Context, err error, msg string) bool {
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    c.Errorf("%v error: %v", msg, err)
    return true
  }
  return false
}

func upload(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  c.Infof("Requested URL: %v", r.URL)
  if r.Body != nil {
    bytes, err := ioutil.ReadAll(r.Body)
    if checkErr(w, c, err, "Failed to read request.") {
      return
    }
    
    c.Infof("Received %v bytes", len(bytes))
    c.Infof("ContentLength %v", r.ContentLength)
    
    request := &htu21df.UploadRequest{}
    err = proto.Unmarshal(bytes, request)
    if checkErr(w, c, err, "Unable to parse proto.") {
      return
    }
    
    c.Infof("request contains %v records", len(request.TempAndHumidtyData))
    
    response := &htu21df.UploadResponse{
      NumSaved:  proto.Int32(int32(len(request.TempAndHumidtyData))),
    }
    c.Infof("response %v", response.String())
    
    responseData, err := proto.Marshal(response)
    if checkErr(w, c, err, "Unable to encode response proto.") {
      return
    }
    w.Write(responseData)
  }
}