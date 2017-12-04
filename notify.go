package main

import "fmt"
import "os"
import "net/http"
import "encoding/json"
import "encoding/base64"
import "crypto/sha256"
import "time"
import "log"
import "database/sql"
import _ "github.com/lib/pq"

type M map[string]string

func someHandler(w http.ResponseWriter, r *http.Request) {
	dbName := os.Getenv("DATABASE_URL")
	_, err := sql.Open("postgres", dbName)
	if err != nil {
		log.Fatal(err)
	}
}

type DatosBasicos struct {
	NombreEvento                      string `json:"nombreEvento"`
	CodigoEvento                      string `json:"codigoEvento"`
	NombresPaciente                   string `json:"nombres_paciente"`
	ApellidosPaciente                 string `json:"apellidos_paciente"`
	TipoIdentificacion                string `json:"tipo_identificacion"`
	NumeroIdentificacion              int    `json:"numero_identificacion"`
	Telefono                          int    `json:"telefono"`
	SexoPaciente                      string `json:"sexo_paciente"`
	PaisOcurrencia                    string `json:"pais_ocurrencia"`
	MunicipioOcurrencia               string `json:"municipio_ocurrencia"`
	FechaNacimientoPaciente           int64  `json:"fecha_nacimiento_paciente"`
	DepartamentoOcurrenciaCaso        string `json:"departamento_ocurrencia_caso"`
	LocalidadOcurrenciaCaso           string `json:"localidad_ocurrencia_caso"`
	BarrioOcurrenciaCaso              string `json:"barrio_ocurrencia_caso"`
	CabeceraCentroRuralOcurrenciaCaso string `json:"cabecera_centro_rural_ocurrencia_caso"`
	VeredaZonaOcurrenciaCaso          string `json:"vereda_zona_ocurrencia_caso"`
	AreaOcurrenciaCaso                string `json:"area_ocurrencia_caso"`
	OcupacionPaciente                 string `json:"ocupacion_paciente"`
	TipoRegimenSalud                  string `json:"tipo_regimen_salud"`
	NombreAdministradoraSalud         string `json:"nombre_administradora_salud"`
	PertenenciaEtnica                 string `json:"pertenencia_etnica"`
	Discapacitados                    bool   `json:"discapacitados"`
	Migrantes                         bool   `json:"migrantes"`
	Gestantes                         bool   `json:"gestantes"`
	InfantilCargoIcbf                 bool   `json:"infantil_cargo_icbf"`
	Desmovilizados                    bool   `json:"desmovilizados"`
	VictimasViolenciaArmada           bool   `json:"victimas_violencia_armada"`
	Desplazados                       bool   `json:"desplazados"`
	Carcelarios                       bool   `json:"carcelarios"`
	Indigentes                        bool   `json:"indigentes"`
	MadresComunitarias                bool   `json:"madres_comunitarias"`
	CentrosPsiquiatricos              bool   `json:"centros_psiquiatricos"`
	OtrosGruposPoblacionales          bool   `json:"otros_grupos_poblacionales"`
	DepartamentoResidencia            string `json:"departamento_residencia"`
	MunicipioResidencia               string `json:"municipio_residencia"`
	DireccionResidencia               string `json:"direccion_residencia"`
	FechaInicioSintomas               int64  `json:"fecha_inicio_sintomas"`
	FechaConsulta                     int64  `json:"fecha_consulta"`
	ClasificacionInicialCaso          string `json:"clasificacion_inicial_caso"`
	Hospitalizado                     bool   `json:"hospitalizado"`
	FechaHospitalizacion              int64  `json:"fecha_hospitalizacion"`
	CondicionFinal                    string `json:"condicion_final"`
	FechaDefuncion                    int64  `json:"fecha_defuncion"`
	NumeroCertificadoDefuncion        int    `json:"numero_certificado_defuncion"`
	CausaBasicaMuerte                 string `json:"causa_basica_muerte"`
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
	Nombres       string `json:"nombres"`
	Apellidos     string `json:"apellidos"`
	NombreUsuario string `json:"nombreUsuario"`
	Telefono      string `json:"telefono"`
}

type Credenciales struct {
	Token         string  `json:"token"`
	NombreUsuario Usuario `json:"nombreUsuario"`
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(`content-type`, `application/json`)
	json.NewEncoder(w).Encode(&M{"message": "Hello World"})
}

func datosBasicosHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		basicos := &DatosBasicos{}
		data := json.NewDecoder(r.Body)
		data.Decode(basicos)
		fmt.Println(time.Unix(0, int64(time.Millisecond)*basicos.FechaConsulta).
			Format(time.RFC3339))
		fmt.Println(time.Unix(0, int64(time.Millisecond)*basicos.FechaNacimientoPaciente).
			Format(time.RFC3339))
		fmt.Println(basicos)
	}
}

func activateCors(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "content-type")
		if r.Method == http.MethodOptions {
			return
		}
		f.ServeHTTP(w, r)
	}
}

func obtenerCredenciales(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		datosRegistro := &DatosRegistro{}

		err := json.NewDecoder(r.Body).Decode(datosRegistro)
		if err != nil {
			http.Error(w, "Error en los datos", http.StatusBadRequest)
			return
		}
		fmt.Printf("%s", datosRegistro)

		h := sha256.New()
		_, err = h.Write([]byte(datosRegistro.RepetirClave))

		if err != nil {
			http.Error(w, "Error Inesperado", http.StatusInternalServerError)
			return
		}

		token := base64.StdEncoding.EncodeToString(h.Sum(nil))

		credenciales := &Credenciales{token,
			Usuario{datosRegistro.Nombres,
				datosRegistro.Apellidos,
				datosRegistro.NombreUsuario,
				datosRegistro.Telefono}}

		encoder := json.NewEncoder(w)
		encoder.Encode(credenciales)
	}
}

func main() {
	port := "3030" //os.Getenv("PORT")
	if port == "" {
		log.Fatalln("$PORT must be set")
	}
	mux := http.NewServeMux()
	mux.HandleFunc(`/home`, mainPage)
	mux.HandleFunc("/url", activateCors(obtenerCredenciales))
	mux.HandleFunc("/", activateCors(datosBasicosHandler))
	fmt.Printf("Serving in %s", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), mux)
}
