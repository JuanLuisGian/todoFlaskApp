from flask import Flask, render_template,jsonify, request,url_for, redirect
from flask_restful import Resource, Api
from flask_cors import CORS
import requests
app = Flask(__name__)
CORS(app)

class recordatorio ():
    ID = "0"
    TITULO = ""
    PRIORIDAD = ""
    HECHO = ""
TotalDatos = 5
todos_los_recordatorios  = []
@app.route('/')
def home():
    
    return render_template('index.html',todos_los_recordatorios=todos_los_recordatorios, TotalDatos =TotalDatos )





@app.route('/create-task', methods=['POST'])
def create():
    id_ingresada = request.form['content']
    titulo_ingresado = request.form["titulo"]
    prioridad_ingresada = request.form['prioridad']
    terminado = request.form['hecho_form']
    r = requests.post('http://localhost:5001/recordatorios/'+ id_ingresada,json={"id":"0","titulo":titulo_ingresado,"prioridad":prioridad_ingresada,"hecho":terminado})
    
    return redirect(url_for('home'))  


@app.route('/recordatorios', methods=['GET'])
def recordatorios():
    r = requests.get('http://localhost:5001/recordatorios')
    todos_los_recordatorios.clear()
    for x in range(len(r.json())):
        item_en_json = r.json()[x] 
        nuevodato = recordatorio()
        nuevodato.ID = item_en_json['id']
        nuevodato.PRIORIDAD = item_en_json['prioridad']
        nuevodato.TITULO = item_en_json['titulo']
        nuevodato.HECHO = item_en_json['hecho']
        todos_los_recordatorios.append(nuevodato)    
    global TotalDatos 
    TotalDatos = len(r.json())

    return redirect(url_for('home')) 
#@app.route('/delete/<id>')
#def delete(id):
   


if __name__ == '__main__':
    app.run(debug = True)
