/* Crear un programa para gestionar estudiantes en una universidad. El sistema debe permitir:
1. Registrar estudiantes con su información básica. (Nombre, Edad, Carrera, Calificaciones)
Tip: Las calificaciones que sea un map en dónde el key es la materia y el value es la calificación.
2. Almacenar sus calificaciones por materia.
3. Calcular el promedio de calificaciones de cada estudiante.
4. Determinar si aprueban o no usando condicionales. (80, 70)
5. Agrupar a los estudiantes por carrera. (Maestria, Univeridad)
6. Permitir modificar sus calificaciones usando punteros.
7. Calcular estadísticas generales del sistema. */

package main

import (
	"fmt"
)

type materia string

type Estudiantes struct {
	Nombre         string
	Edad           int
	Carrera        string
	Calificaciones map[string]int
}

type EstudianteMaestria struct {
	Estudiantes
}

type EstudiantePregrado struct {
	Estudiantes
}

type Estudiante interface {
	DeterminarAprobacion()
	GetNombre() string
	CalcularPromedio() float32
	GetCalificaciones() map[string]int
}

const VALOR_MINIMO_MAESTRIA = 80
const VALOR_MINIMO_PREGRADO = 70

// 1 Pregrado
// 2 Maestria
func main() {
	var arrEstudiantesMaestria []EstudianteMaestria
	var arrEstudiantesPregrado []EstudiantePregrado
	var estudiantesActivos []Estudiante

	//Guardar Estudiante Maaestria
	e1 := EstudianteMaestria{
		Estudiantes{
			"Joel",
			21,
			"Maestria Seguridad",
			make(map[string]int),
		},
	}

	e1.Calificaciones["Matematica"] = 70
	e1.Calificaciones["Fisica"] = 60

	arrEstudiantesMaestria = append(arrEstudiantesMaestria, e1)

	//Guarda Estudiante Pregrado
	e2 := EstudiantePregrado{
		Estudiantes{
			"Marlon",
			21,
			"Ingenieria",
			make(map[string]int),
		},
	}

	e2.Calificaciones["Matematica"] = 80
	e2.Calificaciones["Fisica"] = 90

	arrEstudiantesPregrado = append(arrEstudiantesPregrado, e2)

	//Agregar estudiantes general
	estudiantesActivos = append(estudiantesActivos, &e1)
	estudiantesActivos = append(estudiantesActivos, &e2)

	for _, e := range estudiantesActivos {
		fmt.Printf("Estudiante: %s\n", e.GetNombre())
		e.DeterminarAprobacion()
		e.CalcularPromedio()
	}

	calcularPromedioGeneral(&estudiantesActivos)
}

func (e *Estudiantes) CalcularPromedio() float32 {
	sumatoria := 0
	for _, v := range e.Calificaciones {
		sumatoria += v
	}

	if sumatoria > 0 {
		sumatoria = sumatoria / len(e.Calificaciones)
		fmt.Printf("Promedio: %d \n", sumatoria)
	}

	return float32(sumatoria)
}

func (e *EstudianteMaestria) DeterminarAprobacion() {
	for k, v := range e.Calificaciones {
		if v >= VALOR_MINIMO_MAESTRIA {
			fmt.Printf("Estudiante aprobado en la materia de %s \n", k)
		} else {
			fmt.Printf("Estudiante no aprobado en la materia de %s\n", k)
		}
	}
}

func (e *EstudiantePregrado) DeterminarAprobacion() {
	for _, v := range e.Calificaciones {
		if v >= VALOR_MINIMO_PREGRADO {
			fmt.Println("Estudiante aprobado")
		} else {
			fmt.Println("Estudiante no aprobado")
		}
	}
}

func (e *Estudiantes) ModificarNotas(materia string, nuevaNota int) {
	_, ok := e.Calificaciones[materia]
	if ok {
		e.Calificaciones[materia] = nuevaNota
		fmt.Println("Nota modificada correctamente")
	} else {
		fmt.Println("Materia no asociada al estudiante")
	}
}

func calcularPromedioGeneral(estudiantesActivos *[]Estudiante) {
	sumatoria := 0
	counter := 0
	for _, ea := range *estudiantesActivos {
		for _, c := range ea.GetCalificaciones() {
			sumatoria += c
			counter += 1
		}
	}

	if counter > 0 {
		sumatoria = sumatoria / counter
		fmt.Printf("Promedio General: %d \n", sumatoria)
	} else {
		fmt.Println("No se pudo calcular el promedio")
	}
}

func (e *EstudianteMaestria) GetNombre() string {
	return e.Nombre
}

func (e *EstudiantePregrado) GetNombre() string {
	return e.Nombre
}

func (e *EstudianteMaestria) CalcularPromedio() float32 {
	return e.Estudiantes.CalcularPromedio()
}

func (e *EstudiantePregrado) CalcularPromedio() float32 {
	return e.Estudiantes.CalcularPromedio()
}

func (e *EstudianteMaestria) GetCalificaciones() map[string]int {
	return e.Estudiantes.Calificaciones
}

func (e *EstudiantePregrado) GetCalificaciones() map[string]int {
	return e.Estudiantes.Calificaciones
}
