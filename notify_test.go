package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var fichaTest = DatosBasicos{
	bson.NewObjectId(),
	"Degue",
	"123",
	time.Now(),
	"Victor Samuel", "Mosquera Artamonov", "CedulaCiudadania", 1087998004,
	3207922369,
	"Masculino",
	"CO",
	"Dosquebradas",
	time.Unix(600000000, 0),
	"Risaralda",
	"Comuna 1",
	"Santa Monica",
	"",
	"",
	"CabeceraMunicipa",
	"Medicos Generales",
	"Contributivo",
	"Nueva EPS",
	"Otro",
	false, false, false, false, false, false, false, false, false, false, false,
	true,
	"Risaralda", "Pereira", "Crr 11 bis No. 1-05", time.Unix(1512491920, 0),
	time.Unix(1512491920, 0),
	"Sospechoso",
	false,
	time.Unix(0, 0),
	"Vivo",
	time.Unix(0, 0),
	0,
	"",
}

func TestDatosBasicosHandler(t *testing.T) {
	var dataBody bytes.Buffer
	err := json.NewEncoder(&dataBody).Encode(fichaTest)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, "/datosbasicos", &dataBody)
	if err != nil {
		t.Fatal(err)
	}

	session, err := mgo.Dial(os.Getenv("MONGODB_URI"))
	defer session.Close()
	if err != nil {
		t.Fatal(err)
	}

	handlerWithServer := server{session, "test"}
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(handlerWithServer.datosBasicosHandler)

	handler.ServeHTTP(res, req)

	if status := res.Code; status != http.StatusOK {
		t.Errorf(`handler returned wrong status code: expected %v got %v`,
			http.StatusOK, status)
	}

	mapper := DatosBasicos{}
	id := bson.ObjectIdHex(res.Body.String())
	err = session.DB("test").C("datosBasicos").FindId(id).One(&mapper)
	if err != nil {
		t.Fatal(id, err)
		return
	}

	if fichaTest.NombresPaciente != mapper.NombresPaciente {
		t.Errorf("Expected %v got %v", fichaTest.NombresPaciente, mapper.NombresPaciente)
	}

	err = session.DB("test").C("datosBasicos").DropCollection()
	if err != nil {
		t.Fatal(err)
	}

}

func TestActivateCorsMiddleware(t *testing.T) {
	req, err := http.NewRequest(http.MethodOptions, "/datosbasicos", nil)
	req.Header.Set(`Origin`, `test.vsmd`)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(activateCors(func(w http.ResponseWriter, r *http.Request) {
		return
	}))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf(`handler returned wrong status code: expected %v got %v`,
			http.StatusOK, status)
	}

	if headerAllowMethod := rr.Header().
		Get("Access-Control-Allow-Methods"); headerAllowMethod != "POST, OPTIONS" {
		t.Errorf("Wrong methods allowed %v", headerAllowMethod)
	}

	if headerAllowOrigin := rr.Header().
		Get("Access-Control-Allow-Origin"); headerAllowOrigin != "test.vsmd" {
		t.Errorf("Wrong origins allowed: %v", headerAllowOrigin)
	}

	if headerContentType := rr.Header().
		Get("Content-Type"); headerContentType != "application/json" {
		t.Errorf("Wrong origins allowed: %v", headerContentType)
	}
}

func TestDataBaseConection(t *testing.T) {
	session, err := mgo.Dial(os.Getenv("MONGODB_URI"))
	defer session.Close()

	if err != nil {
		t.Fatal("Failed conection to database", err)
	}

	c := session.DB("test").C("people")
	err = c.Insert(&Usuario{"Victor", "Mosquera", "victorsamuelmd", "3207922369"})
	if err != nil {
		t.Error(err)
	}

	var result Usuario
	err = c.Find(bson.M{"nombres": "Victor"}).One(&result)
	if err != nil {
		t.Error(err)
	}

	err = c.DropCollection()
	if err != nil {
		t.Error(err)
	}
}
