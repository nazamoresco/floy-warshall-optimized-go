package main

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"sync"
	"time"
)

const (
	MatrixSize = 4
  Threads = 4
  Infinity = math.MaxInt
)

var (
	path_weight_matrix []int
	next_vertex_matrix []int

	// Vars
	edges float64
	row_vertex int
	row_index_cached int
	column_vertex int
	index_cached int
	stop_vertex int
	origin_vertex int
	target_vertex int
	from_stop_index_cached int
	origin_index_cached int
	to_stop_index_cached int
	path_weight_from_stop int
	path_weight_to_stop int
	start time.Time
	total_time time.Duration
	previous_path int
	current_path int
	rows_per_thread int
	threads_group sync.WaitGroup
)



func process_rows(thread_index int, threads_group *sync.WaitGroup) {
	var (
		index_cached int
		origin_vertex int
		target_vertex int
		origin_index_cached int
		to_stop_index_cached int
		path_weight_from_stop int
		path_weight_to_stop int
		previous_path int
		current_path int
	)

	for origin_vertex = thread_index * rows_per_thread; origin_vertex < (thread_index + 1) * rows_per_thread; origin_vertex++ {
		origin_index_cached = origin_vertex * MatrixSize
		for target_vertex = 0; target_vertex < MatrixSize; target_vertex++ {
			index_cached = origin_index_cached + target_vertex
			to_stop_index_cached = origin_index_cached + stop_vertex

			// if there is no path to the stop or from the stop go to next iteration
			path_weight_from_stop = path_weight_matrix[from_stop_index_cached + target_vertex]
			path_weight_to_stop = path_weight_matrix[to_stop_index_cached]
			if path_weight_to_stop == Infinity || path_weight_from_stop == Infinity {
				continue
			}

			previous_path = path_weight_matrix[index_cached]
			current_path = path_weight_to_stop + path_weight_from_stop
			if current_path < previous_path {
				path_weight_matrix[index_cached] = current_path
				next_vertex_matrix[index_cached] = next_vertex_matrix[to_stop_index_cached]
			}
		}
	}
	threads_group.Done()
}


func main() {
	//
	// Initialization of matrixes and variables ~>
	//
	// Matrix with the minimum distances between each pair of vertices
	// Origin are rows, targes are columns
	path_weight_matrix = make([]int, MatrixSize * MatrixSize)
	next_vertex_matrix = make([]int, MatrixSize * MatrixSize)

	// Adjencency matrix representing the graph
	graph_matrix := make([]int, MatrixSize * MatrixSize)

	// Initialize graph
	edges = float64(MatrixSize * (MatrixSize - 1)) * 0.7 // 70% Density
	for row_vertex = 0; row_vertex < MatrixSize; row_vertex++ {
		row_index_cached = row_vertex * MatrixSize
		for column_vertex = 0; column_vertex < MatrixSize; column_vertex ++ {
			index_cached = row_index_cached + column_vertex
			if(row_vertex == column_vertex) {
				graph_matrix[index_cached] = 0
			} else if(edges < 1) {
				graph_matrix[index_cached] = -1
			} else {
				graph_matrix[index_cached] = rand.Intn(100)
				edges = edges - 1
			}
		}
	}

	// Initialize path and next
	for i := range path_weight_matrix {
		if(graph_matrix[i] == -1) {
			path_weight_matrix[i] = Infinity
			next_vertex_matrix[i] = -1
		} else {
			path_weight_matrix[i] = graph_matrix[i]
			next_vertex_matrix[i] = i % MatrixSize
		}
	}


	//
	// Initialize result matrixes and calculate shortest path  ~>
	//

	// IMPORTANT AND DEFINITION OF THREADS
	start = time.Now()
	rows_per_thread = MatrixSize/Threads
	for stop_vertex = 0; stop_vertex < MatrixSize; stop_vertex++ {
		from_stop_index_cached = stop_vertex * MatrixSize

		threads_group.Add(Threads)
		for thread_index := 1; thread_index < Threads; thread_index++ {
			go process_rows(thread_index, &threads_group)
		}

		process_rows(0, &threads_group)
		threads_group.Wait()
	}

	//
	// Return values  ~>D
	//
	total_time = time.Since(start)

	fmt.Println()
	fmt.Println("Input graph")
	printMatrix(graph_matrix)

	fmt.Println()
	fmt.Println("Path weight matrix")
	printMatrix(path_weight_matrix)

	fmt.Println()
	fmt.Println("Next vertex matrix")
	printMatrix(next_vertex_matrix)

	fmt.Println()
	fmt.Println("Time:")
	fmt.Println(total_time.Seconds())
}


// printMatrix: Format and prints a given matrix
// Assumes first index is row, second index is column
// Prepared for 4 digit numbers, more digits would break this
func printMatrix(matrix []int) {
	// Compact version
	if(MatrixSize > 25) {
		fmt.Println("[ ")
		for row_index := 0; row_index < MatrixSize; row_index ++ {
			fmt.Print("  ", row_index, ": { ")
			for column_index := 0; column_index < MatrixSize; column_index ++ {
				fmt.Print(column_index, ":", matrix[row_index * MatrixSize + column_index], ", ")
			}
			fmt.Println("}")
		}
		fmt.Println("]")
		return
	}


	// Pretty version
	fmt.Print(" ")
	fmt.Print(strings.Repeat("-", 4 + MatrixSize * 8))

	fmt.Println()

	fmt.Print("|")
	fmt.Print(strings.Repeat(" ", 4))
	for column_index := 0; column_index < MatrixSize; column_index ++ {
		fmt.Print("|    ")
		if(column_index < 10) {
			fmt.Print(" ")
		}
		fmt.Print(column_index)
		fmt.Print(" ")
	}
	fmt.Print("|")
	fmt.Println()

	fmt.Print(" ")
	fmt.Print(strings.Repeat("-", 4 + MatrixSize * 8))
	fmt.Println()

	for row_index := 0; row_index < MatrixSize; row_index ++ {
		fmt.Print("| ")
		if(row_index < 10) {
			fmt.Print(" ")
		}
		fmt.Print(row_index)
		fmt.Print(" ")

		for column_index := 0; column_index < MatrixSize; column_index ++ {
			fmt.Print("| ")
			switch value := matrix[row_index * MatrixSize + column_index]; value {
				case Infinity:
					fmt.Print("   +∞")
				default:
					spaces := 1
					abs_value := value
					if value < 0 {
						abs_value = -1 * value
						spaces = 0
					}

					if(abs_value < 10) {
						spaces += 3
					} else if (abs_value < 100) {
						spaces += 2
					} else if (abs_value < 1000) {
						spaces += 1
					}

					fmt.Print(strings.Repeat(" ", spaces))
					fmt.Print(value)
			}
			fmt.Print(" ")
		}
		fmt.Print("|")
		fmt.Println()
	}

	fmt.Print(" ")
	fmt.Print(strings.Repeat("-", 4 + MatrixSize * 8))
	fmt.Println()
}