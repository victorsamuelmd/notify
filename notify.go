package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"

	"github.com/gorilla/mux"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	baseDeDatos           = "victorsamuelmd"
	coleccionDatosBasicos = "datosBasicos"
	datosBasicosIdUrl     = "datos"
)

type DatosBasicos struct {
	ID                                bson.ObjectId `json:"id" bson:"_id"`
	NombreEvento                      string        `json:"nombre_evento" bson:"nombre_evento"`
	CodigoEvento                      string        `json:"codigo_evento" bson:"codigo_evento"`
	FechaNotificacion                 time.Time     `json:"fecha_notificacion" bson:"fecha_notificacion"`
	NombresPaciente                   string        `json:"nombres_paciente" bson:"nombres_paciente"`
	ApellidosPaciente                 string        `json:"apellidos_paciente" bson:"apellidos_paciente"`
	TipoIdentificacion                string        `json:"tipo_identificacion" bson:"tipo_identificacion"`
	NumeroIdentificacion              uint64        `json:"numero_identificacion" bson:"numero_identificacion"`
	Telefono                          uint64        `json:"telefono" bson:"telefono"`
	SexoPaciente                      string        `json:"sexo_paciente" bson:"sexo_paciente"`
	PaisOcurrencia                    string        `json:"pais_ocurrencia" bson:"pais_ocurrencia"`
	MunicipioOcurrencia               string        `json:"municipio_ocurrencia" bson:"municipio_ocurrencia"`
	FechaNacimientoPaciente           time.Time     `json:"fecha_nacimiento_paciente" bson:"fecha_nacimiento_paciente"`
	DepartamentoOcurrenciaCaso        string        `json:"departamento_ocurrencia_caso" bson:"departamento_ocurrencia_caso"`
	LocalidadOcurrenciaCaso           string        `json:"localidad_ocurrencia_caso" bson:"localidad_ocurrencia_caso"`
	BarrioOcurrenciaCaso              string        `json:"barrio_ocurrencia_caso" bson:"barrio_ocurrencia_caso"`
	CabeceraCentroRuralOcurrenciaCaso string        `json:"cabecera_centro_rural_ocurrencia_caso" bson:"cabecera_centro_rural_ocurrencia_caso"`
	VeredaZonaOcurrenciaCaso          string        `json:"vereda_zona_ocurrencia_caso" bson:"vereda_zona_ocurrencia_caso"`
	AreaOcurrenciaCaso                string        `json:"area_ocurrencia_caso" bson:"area_ocurrencia_caso"`
	OcupacionPaciente                 string        `json:"ocupacion_paciente" bson:"ocupacion_paciente"`
	TipoRegimenSalud                  string        `json:"tipo_regimen_salud" bson:"tipo_regimen_salud"`
	NombreAdministradoraSalud         string        `json:"nombre_administradora_salud" bson:"nombre_administradora_salud"`
	PertenenciaEtnica                 string        `json:"pertenencia_etnica" bson:"pertenencia_etnica"`
	Discapacitados                    bool          `json:"discapacitados" bson:"discapacitados"`
	Migrantes                         bool          `json:"migrantes" bson:"migrantes"`
	Gestantes                         bool          `json:"gestantes" bson:"gestantes"`
	InfantilCargoIcbf                 bool          `json:"infantil_cargo_icbf" bson:"infantil_cargo_icbf"`
	Desmovilizados                    bool          `json:"desmovilizados" bson:"desmovilizados"`
	VictimasViolenciaArmada           bool          `json:"victimas_violencia_armada" bson:"victimas_violencia_armada"`
	Desplazados                       bool          `json:"desplazados" bson:"desplazados"`
	Carcelarios                       bool          `json:"carcelarios" bson:"carcelarios"`
	Indigentes                        bool          `json:"indigentes" bson:"indigentes"`
	MadresComunitarias                bool          `json:"madres_comunitarias" bson:"madres_comunitarias"`
	CentrosPsiquiatricos              bool          `json:"centros_psiquiatricos" bson:"centros_psiquiatricos"`
	OtrosGruposPoblacionales          bool          `json:"otros_grupos_poblacionales" bson:"otros_grupos_poblacionales"`
	DepartamentoResidencia            string        `json:"departamento_residencia" bson:"departamento_residencia"`
	MunicipioResidencia               string        `json:"municipio_residencia" bson:"municipio_residencia"`
	DireccionResidencia               string        `json:"direccion_residencia" bson:"direccion_residencia"`
	FechaInicioSintomas               time.Time     `json:"fecha_inicio_sintomas" bson:"fecha_inicio_sintomas"`
	FechaConsulta                     time.Time     `json:"fecha_consulta" bson:"fecha_consulta"`
	ClasificacionInicialCaso          string        `json:"clasificacion_inicial_caso" bson:"clasificacion_inicial_caso"`
	Hospitalizado                     bool          `json:"hospitalizado" bson:"hospitalizado"`
	FechaHospitalizacion              time.Time     `json:"fecha_hospitalizacion" bson:"fecha_hospitalizacion"`
	CondicionFinal                    string        `json:"condicion_final" bson:"condicion_final"`
	FechaDefuncion                    time.Time     `json:"fecha_defuncion" bson:"fecha_defuncion"`
	NumeroCertificadoDefuncion        int64         `json:"numero_certificado_defuncion" bson:"numero_certificado_defuncion"`
	CausaBasicaMuerte                 string        `json:"causa_basica_muerte" bson:"causa_basica_muerte"`
}

type server struct {
	DB     *mgo.Session
	dbName string
}
type DatosRegistro struct {
	Nombres       string `json:"nombres"`
	Apellidos     string `json:"apellidos"`
	NombreUsuario string `json:"nombreUsuario"`
	Telefono      string `json:"telefono"`
	Clave         string `json:"clave"`
	RepetirClave  string `json:"repetirClave"`
}

type Usuario struct {
	Nombres       string `json:"nombres" bson:"nombres"`
	Apellidos     string `json:"apellidos" bson:"apellidos"`
	NombreUsuario string `json:"nombreUsuario" bson:"nombreUsuario"`
	Telefono      string `json:"telefono" bson:"telefono"`
}

type Credenciales struct {
	Token         string  `json:"token"`
	NombreUsuario Usuario `json:"nombreUsuario"`
}

func (s *server) datosBasicosHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		basicos := DatosBasicos{}
		err := json.NewDecoder(r.Body).Decode(&basicos)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		basicos.ID = bson.NewObjectId()
		basicos.FechaNotificacion = time.Now()

		db := s.DB.Clone()
		defer db.Close()

		c := db.DB(s.dbName).C(coleccionDatosBasicos)
		err = c.Insert(&basicos)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, `"%s"`, basicos.ID.Hex())
	}
}

func (s *server) mostrarDatosBasicos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(`Content-Type`, `application/json`)
	db := s.DB.Clone()
	defer db.Close()
	vars := mux.Vars(r)
	c := db.DB(s.dbName).C(coleccionDatosBasicos)
	var datosBasicos DatosBasicos
	err := c.FindId(bson.ObjectIdHex(vars[datosBasicosIdUrl])).One(&datosBasicos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	err = json.NewEncoder(w).Encode(&datosBasicos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func (s *server) mostrarDatosBasicosSVG(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(`Content-Type`, `text/html`)
	db := s.DB.Clone()
	defer db.Close()

	temp, err := template.ParseFiles("./datos_basicos.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	c := db.DB(s.dbName).C(coleccionDatosBasicos)
	var datosBasicos DatosBasicos
	err = c.FindId(bson.ObjectIdHex(vars[datosBasicosIdUrl])).One(&datosBasicos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	err = temp.Execute(w, datosBasicos)
}

func activateCors(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "content-type")
		w.Header().Set(`Content-Type`, `application/json`)
		if r.Method == http.MethodOptions {
			return
		}
		f.ServeHTTP(w, r)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalln("$PORT must be set")
	}

	session, err := mgo.Dial(os.Getenv("MONGODB_URI"))
	if err != nil {
		log.Fatal(err.Error())
	}

	handler := server{session, baseDeDatos}

	r := mux.NewRouter()
	r.HandleFunc(fmt.Sprintf(`/verdatosbasicos/{%s}.json`, datosBasicosIdUrl), handler.mostrarDatosBasicos)
	r.HandleFunc(fmt.Sprintf(`/verdatosbasicos/{%s}`, datosBasicosIdUrl), handler.mostrarDatosBasicosSVG)
	r.HandleFunc("/datosbasicos", activateCors(handler.datosBasicosHandler))
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))

	fmt.Printf("Serving in %s", port)

	http.ListenAndServe(fmt.Sprintf(":%s", port), r)
}
