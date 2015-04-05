package temp_sensor_app

import (
  "code.google.com/p/goprotobuf/proto"

  //"bytes"
  "errors"
  _ "fmt"
  _ "net/http"
  "net/http/httptest"
  "testing"

  "appengine"
  "appengine/aetest"
  "appengine/datastore"

  "htu21df"
)

type asserts struct {
  t *testing.T
}

func (a *asserts) assertStrsEqual(actual, expected string) {
  if actual != expected {
    a.t.Errorf("Got %v but expected %v", actual, expected)
  }
}

func (a *asserts) assertIntsEqual(actual, expected int) {
  if actual != expected {
    a.t.Errorf("Got %v but expected %v", actual, expected)
  }
}

func (a *asserts) failOnError(err error) {
  if err != nil {
   a.t.Fatal(err)
  }
}

func (a *asserts) IsTrue(test bool) {
  if !test {
   a.t.Fatal("Expected true value.")
  }
}

func (a *asserts) IsFalse(test bool) {
  if test {
   a.t.Fatal("Expected false value.")
  }
}

func (a *asserts) IsNil(test interface{}) {
  if test != nil {
   a.t.Errorf("Expected nil value, but was %v.", test)
  }
}

func newContext(t *testing.T) aetest.Context {
  a := asserts{t}

  c, err := aetest.NewContext(nil)
  a.failOnError(err)
  return c
}

func TestRoot(t *testing.T) {
  a := asserts{t}
  w := httptest.NewRecorder()
  root(w, nil)
  a.assertStrsEqual(w.Body.String(), "Temp and humidity.")
}

func TestRespondWith400(t *testing.T) {
  a := asserts{t}
  w := httptest.NewRecorder()
  c := newContext(t)
  defer c.Close()

  respondWith400(w, c, errors.New("Test"), "Test Msg")
  a.assertIntsEqual(w.Code, 400)
}

func TestRespondWith500(t *testing.T) {
  a := asserts{t}
  w := httptest.NewRecorder()
  c := newContext(t)
  defer c.Close()

  respondWith500(w, c, errors.New("Test"), "Test Msg")
  a.assertIntsEqual(w.Code, 500)
}

func TestCheckErr(t *testing.T) {
  a := asserts{t}

  w := httptest.NewRecorder()
  c := newContext(t)
  defer c.Close()

  a.IsFalse(checkErr(w, c, nil, "No error"))
  a.IsFalse(w.Code == 500)

  w = httptest.NewRecorder()
  a.IsTrue(checkErr(w, c, errors.New("Error"), "Error!"))
  a.assertIntsEqual(w.Code, 500)
}

func TestNewRecordKey(t *testing.T) {
  a := asserts{t}

  c := newContext(t)
  defer c.Close()

  testKey := newRecordKey(appengine.Context(c))
  a.assertStrsEqual(testKey.Kind(), EntityName)
  //a.IsNil(testKey.Parent())
  a.IsTrue(testKey.Incomplete())
}

func TestSaveRecords(t *testing.T) {
  a := asserts{t}

  c := newContext(t)
  defer c.Close()

  newKeys, err := saveRecords(c, &htu21df.UploadRequest{})
  a.IsNil(err)
  a.assertIntsEqual(len(newKeys), 0)

  req := htu21df.UploadRequest{
    DataAndClientId: []*htu21df.DataAndClientId{
      &htu21df.DataAndClientId{
        ClientId: proto.String("client1"),
        TempAndHumidtyData: &htu21df.TempAndHumidity{
          RecordedTimestampMs: proto.Int64(1411103037254),
          TempDegreesC: proto.Float64(23.95),
          PercentRelativeHumidity: proto.Float64(68.78),
          SensorName: proto.String("htu21df"),
          Debug: proto.String("./temp_db_logger.py"),
        },
      },
      &htu21df.DataAndClientId{
        ClientId: proto.String("client2"),
        TempAndHumidtyData: &htu21df.TempAndHumidity{
          RecordedTimestampMs: proto.Int64(1411104501012),
          TempDegreesC: proto.Float64(23.839),
          PercentRelativeHumidity: proto.Float64(68.547),
          SensorName: proto.String("htu21df"),
          Debug: proto.String("./temp_db_logger.py"),
        },
      },
      &htu21df.DataAndClientId{
        ClientId: proto.String("client3"),
        TempAndHumidtyData: &htu21df.TempAndHumidity{
          RecordedTimestampMs: proto.Int64(1411104611852),
          TempDegreesC: proto.Float64(23.828),
          PercentRelativeHumidity: proto.Float64(68.615),
          SensorName: proto.String("htu21df"),
          Debug: proto.String("./temp_db_logger.py"),
        },
      },
    },
  }

  newKeys, err = saveRecords(c, &req)
  a.IsNil(err)
  a.assertIntsEqual(len(newKeys), 3)

  byClientId := make(map[string]*htu21df.TempAndHumidity)
  for _, dataAndId := range req.DataAndClientId {
    byClientId[*dataAndId.ClientId] = dataAndId.TempAndHumidtyData
  }

  for _, newId := range newKeys {
    rec := &tempAndHumidityRecord{}
    err = datastore.Get(c, newId, rec)
    a.IsNil(err)

    reqRec := byClientId[rec.ClientId]
    a.IsTrue(rec.RecordedTimestampMs == *reqRec.RecordedTimestampMs)
    a.IsTrue(rec.TempDegreesC == *reqRec.TempDegreesC)
    a.IsTrue(rec.PercentRelativeHumidity == *reqRec.PercentRelativeHumidity)
    a.assertStrsEqual(*reqRec.SensorName, rec.SensorName)
    a.assertStrsEqual(*reqRec.Debug, rec.Debug)
  }
}

func TestVerifySaved(t *testing.T) {
  a := asserts{t}

  c := newContext(t)
  defer c.Close()

  // Test empty
  missedIds, err := verifySaved(c, []*htu21df.DataAndClientId{}, []*datastore.Key{})
  a.IsNil(err)
  a.assertIntsEqual(len(missedIds), 0)

  // Test len's agree
  missedIds, err = verifySaved(c,
      []*htu21df.DataAndClientId{&htu21df.DataAndClientId{}, &htu21df.DataAndClientId{}},
      []*datastore.Key{newRecordKey(c), newRecordKey(c)})
  a.IsNil(err)
  a.assertIntsEqual(len(missedIds), 0)

  // Test 3 requested, 1 saved
  sentRecords := []*htu21df.DataAndClientId{
    &htu21df.DataAndClientId{
      ClientId: proto.String("client1"),
      TempAndHumidtyData: &htu21df.TempAndHumidity{
        RecordedTimestampMs: proto.Int64(1411103037254),
        TempDegreesC: proto.Float64(23.95),
        PercentRelativeHumidity: proto.Float64(68.78),
        SensorName: proto.String("htu21df"),
        Debug: proto.String("./temp_db_logger.py"),
      },
    },
    &htu21df.DataAndClientId{
      ClientId: proto.String("client2"),
      TempAndHumidtyData: &htu21df.TempAndHumidity{
        RecordedTimestampMs: proto.Int64(1411104501012),
        TempDegreesC: proto.Float64(23.839),
        PercentRelativeHumidity: proto.Float64(68.547),
        SensorName: proto.String("htu21df"),
        Debug: proto.String("./temp_db_logger.py"),
      },
    },
    &htu21df.DataAndClientId{
      ClientId: proto.String("client3"),
      TempAndHumidtyData: &htu21df.TempAndHumidity{
        RecordedTimestampMs: proto.Int64(1411104611852),
        TempDegreesC: proto.Float64(23.828),
        PercentRelativeHumidity: proto.Float64(68.615),
        SensorName: proto.String("htu21df"),
        Debug: proto.String("./temp_db_logger.py"),
      },
    },
  }

  savedRec := &tempAndHumidityRecord {
    ClientId: "client3",
    RecordedTimestampMs: 1411104611852,
    TempDegreesC: 23.828,
    PercentRelativeHumidity: 68.615,
    SensorName: "htu21df",
    Debug: "./temp_db_logger.py",
  }
  savedKey, err := datastore.Put(c, newRecordKey(c), savedRec)
  a.IsNil(err)

  missedIds, err = verifySaved(c, sentRecords, []*datastore.Key{savedKey})
  a.IsNil(err)
  a.assertIntsEqual(len(missedIds), 2)
  a.assertStrsEqual("client1", missedIds[0])
  a.assertStrsEqual("client2", missedIds[1])
}