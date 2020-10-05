//este codigo abarca la conexión a base de datos, modelos y respuestas POST
//Cesar Gianluigi Figueroa Lima
//cesar0990@gmail.com

package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const nombreBDusuarios = "Usuarios"
const puerto = 800
const nombreColeccion = "usuario"
const nombreColeccionRecorda = "recordatorios"

//Conexiones base de datos

func iniciarMongoConexion() (*mongo.Client, error) {

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	return client, nil
}

func coleccionMongoDB(DbName string, CollectionName string) (*mongo.Collection, error) {
	client, err := iniciarMongoConexion()

	if err != nil {
		return nil, err
	}

	collection := client.Database(DbName).Collection(CollectionName)

	return collection, nil
}

//MODELOS
type Usuario struct {
	ID       string `json:"id,omitempty" bson:"_id,omitempty"`
	NOMBRE   string `json:"nombre,omitempty" bson:"nombre,omitempty"`
	APELLIDO string `json:"apellido,omitempty" bson:"apellido,omitempty"`
	CORREO   string `json:"correo,omitempty" bson:"correo,omitempty"`
}
type Recordatorio struct {
	ID        string `json:"id,omitempty" bson:"_id,omitempty"`
	TITUTLO   string `json:"titulo,omitempty" bson:"titulo,omitempty"`
	PRIORIDAD string `json:"prioridad,omitempty" bson:"prioridad,omitempty"`
	HECHO     string `json:"hecho,omitempty" bson:"hecho,omitempty"`
}

func getUsuarios(w http.ResponseWriter, req *http.Request) {
	var gente []Usuario

	collection, _ := coleccionMongoDB(nombreBDusuarios, nombreColeccion)
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var usuario Usuario
		cursor.Decode(&usuario)
		gente = append(gente, usuario)
	}
	if err := cursor.Err(); err != nil {
	}

	json.NewEncoder(w).Encode(gente)
}
func crearUsuario(w http.ResponseWriter, req *http.Request) {
	parametros := mux.Vars(req)
	var nuevoUsuario Usuario
	_ = json.NewDecoder(req.Body).Decode(&nuevoUsuario)
	nuevoUsuario.ID = parametros["id"] //asignacion de id
	collection, err := coleccionMongoDB(nombreBDusuarios, nombreColeccion)
	if err != nil {

	}
	res, err := collection.InsertOne(context.Background(), nuevoUsuario) //verificar que no de problemas el contexto

	response, _ := json.Marshal(res)
	w.Write(response)
}
func getUsuario(w http.ResponseWriter, req *http.Request) {

	parametros := mux.Vars(req)
	var usuario Usuario
	id := parametros["id"]
	collection, _ := coleccionMongoDB(nombreBDusuarios, nombreColeccion)

	err := collection.FindOne(context.Background(), Usuario{ID: id}).Decode(&usuario)
	if err != nil {
	}
	json.NewEncoder(w).Encode(usuario)

}
func borrarUsuario(w http.ResponseWriter, req *http.Request) {

	parametros := mux.Vars(req)
	var usuario Usuario
	id := parametros["id"]
	collection, _ := coleccionMongoDB(nombreBDusuarios, nombreColeccion)

	err, _ := collection.DeleteOne(context.Background(), Usuario{ID: id})
	if err != nil {
	}
	json.NewEncoder(w).Encode(usuario)

}

//Metodos REST para Recordatorios
func getRecordatorios(w http.ResponseWriter, req *http.Request) {
	var listadoRecordatorios []Recordatorio

	collection, _ := coleccionMongoDB(nombreBDusuarios, nombreColeccionRecorda)
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var recordatorio Recordatorio
		cursor.Decode(&recordatorio)
		listadoRecordatorios = append(listadoRecordatorios, recordatorio)
	}
	if err := cursor.Err(); err != nil {
	}

	json.NewEncoder(w).Encode(listadoRecordatorios)
}
func crearRecordatorio(w http.ResponseWriter, req *http.Request) {
	parametros := mux.Vars(req)
	var nuevoRecordatorio Recordatorio
	_ = json.NewDecoder(req.Body).Decode(&nuevoRecordatorio)
	nuevoRecordatorio.ID = parametros["id"] //asignacion de id
	collection, err := coleccionMongoDB(nombreBDusuarios, nombreColeccionRecorda)
	if err != nil {

	}
	res, err := collection.InsertOne(context.Background(), nuevoRecordatorio) //verificar que no de problemas el contexto

	response, _ := json.Marshal(res)
	w.Write(response)
}
func getRecordatorio(w http.ResponseWriter, req *http.Request) {
	parametros := mux.Vars(req)
	var nuevoRecordatorio Recordatorio
	id := parametros["id"]
	collection, _ := coleccionMongoDB(nombreBDusuarios, nombreColeccionRecorda)

	err := collection.FindOne(context.Background(), Recordatorio{ID: id}).Decode(&nuevoRecordatorio)
	if err != nil {
	}
	json.NewEncoder(w).Encode(nuevoRecordatorio)

}
func borrarRecordatorio(w http.ResponseWriter, req *http.Request) {
	parametros := mux.Vars(req)
	var recordatorioViejo Usuario
	id := parametros["id"]
	collection, _ := coleccionMongoDB(nombreBDusuarios, nombreColeccionRecorda)

	err, _ := collection.DeleteOne(context.Background(), Recordatorio{ID: id})
	if err != nil {
	}
	json.NewEncoder(w).Encode(recordatorioViejo)

}

func main() {
	enrutador := mux.NewRouter()

	//manejo de metodos
	enrutador.HandleFunc("/usuarios", getUsuarios).Methods("GET")
	enrutador.HandleFunc("/usuarios/{id}", getUsuario).Methods("GET")
	enrutador.HandleFunc("/usuarios/{id}", borrarUsuario).Methods("DELETE")
	enrutador.HandleFunc("/usuarios/{id}", crearUsuario).Methods("POST")

	//manejo de metodos de recordatorios
	enrutador.HandleFunc("/recordatorios", getRecordatorios).Methods("GET")
	enrutador.HandleFunc("/recordatorios/{id}", getRecordatorio).Methods("GET")
	enrutador.HandleFunc("/recordatorios/{id}", borrarRecordatorio).Methods("DELETE")
	enrutador.HandleFunc("/recordatorios/{id}", crearRecordatorio).Methods("POST")

	log.Fatal(http.ListenAndServe(":5001", enrutador))

}
