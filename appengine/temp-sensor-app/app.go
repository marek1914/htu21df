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

const EntityName = "TempAndHumidity"

type tempAndHumidityRecord struct {
  ClientId string
  RecordedTimestampMs int64
  TempDegreesC float64
  PercentRelativeHumidity float64
  SensorName string `datastore:",noindex"`
  Debug string `datastore:",noindex"`
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
	key := datastore.NewIncompleteKey(c, EntityName, nil)
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

  response := &htu21df.UploadResponse{}
  if err = buildResponse(c, request, response); err != nil {
    respondWith500(w, c, err, err.Error())
    return
  }

  c.Infof("response %v", response.String())
  responseData, err := proto.Marshal(response)
  if checkErr(w, c, err, "Unable to encode response proto.") {
    return
  }
  w.Write(responseData)
}

func buildResponse(c appengine.Context, request *htu21df.UploadRequest, response *htu21df.UploadResponse) error {
  newKeys, err := saveRecords(c, request)
  if err != nil {
    return err
  }

  unsavedIds, err := verifySaved(c, request.DataAndClientId, newKeys)
  if err != nil {
    return err
  }

  response.NumSaved = proto.Int(len(newKeys))
  response.UnsavedClientIds = unsavedIds
  return nil;
}

func saveRecords(c appengine.Context, request *htu21df.UploadRequest) ([]*datastore.Key, error) {
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
    return nil, err
  }
  return newKeys, nil
}

func verifySaved(c appengine.Context, clientRecords []*htu21df.DataAndClientId, newKeys []*datastore.Key) ([]string, error) {
  numSaved := int32(len(newKeys))
  numRecords := int32(len(clientRecords))

  if numRecords != numSaved && numSaved > 0 {
    // find unsaved client ids
    unsavedIds := make([]string, numRecords - numSaved)
    savedRecords := make([]*tempAndHumidityRecord, numSaved)
    if err := datastore.GetMulti(c, newKeys, savedRecords); err != nil {
      c.Infof("Unable to query saved records")
      return nil, err
    }
    clientIds := make(map[string]bool)
    for _, r := range savedRecords {
      clientIds[r.ClientId] = true
    }
    for i, dataAndId := range clientRecords {
      if !clientIds[dataAndId.GetClientId()] {
        unsavedIds[i] = dataAndId.GetClientId()
      }
    }
    c.Infof("Found unsaved ids %v", unsavedIds)
    return unsavedIds, nil
  }
  return nil, nil
}