package temp_sensor_app

import (
  "errors"
  "fmt"
  "io/ioutil"
  "net/http"

  "code.google.com/p/goprotobuf/proto"

  "appengine"
  "appengine/datastore"

  "htu21df"
)

type tempAndHumidityRecord struct {
	ClientId string
  RecordedTimestampMs int64
  TempDegreesC float64
  PercentRelativeHumidity float64
  SensorName string
  Debug string
}

func init() {
  http.HandleFunc("/", root)
  http.HandleFunc("/upload", upload)
}

func root(w http.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, "Temp and humidity.")
}

func respondWith400(w http.ResponseWriter, c appengine.Context, err error, msg string) {
  w.WriteHeader(http.StatusBadRequest)
  c.Errorf("%v error: %v", msg, err.Error())
}

func respondWith500(w http.ResponseWriter, c appengine.Context, err error, msg string) {
  w.WriteHeader(http.StatusInternalServerError)
  c.Errorf("%v error: %v", msg, err.Error())
}

func checkErr(w http.ResponseWriter, c appengine.Context, err error, msg string) bool {
  if err != nil {
    respondWith500(w, c, err, msg)
    return true
  }
  return false
}

func newRecordKey(c appengine.Context) *datastore.Key {
	key := datastore.NewIncompleteKey(c, "TempAndHumidity", nil)
  return key
}

func upload(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  c.Infof("Requested URL: %v", r.URL)
	
  if (r.Body == nil) {
    respondWith400(w, c, errors.New("Response body is nil"), "Response body is nil")
  }
	
  bytes, err := ioutil.ReadAll(r.Body)
  if checkErr(w, c, err, "Failed to read request.") {
    return
  }

  request := &htu21df.UploadRequest{}
  err = proto.Unmarshal(bytes, request)
  if checkErr(w, c, err, "Unable to parse proto.") {
    return
  }

  numSaved, err := saveRecords(c, request)
  if checkErr(w, c, err, "Unable put entities.") {
    return
  }

  response := &htu21df.UploadResponse{
    NumSaved:  proto.Int32(numSaved),
  }
  c.Infof("response %v", response.String())

  responseData, err := proto.Marshal(response)
  if checkErr(w, c, err, "Unable to encode response proto.") {
    return
  }
  w.Write(responseData)
}

func saveRecords(c appengine.Context, request *htu21df.UploadRequest) (numSaved int32, err error) {
  numRecords := len(request.DataAndClientId)
  c.Infof("request contains %v records", numRecords)
  keys := make([]*datastore.Key, numRecords)
  records := make([]*tempAndHumidityRecord, numRecords)

  for i, dataAndId := range request.DataAndClientId {
  	keys[i] = newRecordKey(c)
  	records[i] = &tempAndHumidityRecord {
  		ClientId: dataAndId.GetClientId(),
  	  RecordedTimestampMs: dataAndId.TempAndHumidtyData.GetRecordedTimestampMs(),
  	  TempDegreesC: dataAndId.TempAndHumidtyData.GetTempDegreesC(),
  	  PercentRelativeHumidity: dataAndId.TempAndHumidtyData.GetPercentRelativeHumidity(),
  	  SensorName: dataAndId.TempAndHumidtyData.GetSensorName(),
  	  Debug: dataAndId.TempAndHumidtyData.GetDebug(),
  	}
  }

  newKeys, err := datastore.PutMulti(c, keys, records)
  if err != nil {
    return 0, err
  }

  return int32(len(newKeys)), nil
}