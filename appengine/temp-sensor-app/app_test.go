package temp_sensor_app

import (
  //"code.google.com/p/goprotobuf/proto"

  //"bytes"
  "errors"
  _ "fmt"
  _ "net/http"
  "net/http/httptest"
  "testing"

  "appengine"
  "appengine/aetest"
  //"appengine/datastore"

  //"htu21df"
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
  if test == nil {
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
  a.IsNil(testKey.Parent())
  a.IsTrue(testKey.Incomplete())
}