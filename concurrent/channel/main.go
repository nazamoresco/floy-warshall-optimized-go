package main

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"
)

const (
	MatrixSize = 16384
  Threads = 2
  Infinity = math.MaxInt
	rows_per_thread = MatrixSize/Threads
	MainThreadIndex = 0
)

func process_rows(
	rows_per_thread int,
	thread_index int,
	origin_weight_matrix []int,
	origin_next_matrix []int,
	thread_id_channel chan int,
	path_rows_channel chan []int,
	weight_rows_channel chan []int,
	stop_vertex_weight_matrix_channel chan []int,
	) {
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
		stop_weight_matrix []int
		stop_vertex int
	)

	for stop_vertex = 0; stop_vertex < MatrixSize; stop_vertex++ {
		stop_weight_matrix = <- stop_vertex_weight_matrix_channel

		for origin_vertex = 0; origin_vertex < rows_per_thread; origin_vertex++ {
			origin_index_cached = origin_vertex * MatrixSize
			for target_vertex = 0; target_vertex < MatrixSize; target_vertex++ {
				index_cached = origin_index_cached + target_vertex
				to_stop_index_cached = origin_index_cached + stop_vertex

				// if there is no path to the stop or from the stop go to next iteration
				path_weight_from_stop = stop_weight_matrix[target_vertex]
				path_weight_to_stop = origin_weight_matrix[to_stop_index_cached]
				if path_weight_to_stop == Infinity || path_weight_from_stop == Infinity {
					continue
				}

				// fmt.Println(index_cached, origin_weight_matrix, stop_weight_matrix)
				previous_path = origin_weight_matrix[index_cached]
				current_path = path_weight_to_stop + path_weight_from_stop

				if current_path < previous_path {
					origin_weight_matrix[target_vertex] = current_path
					origin_next_matrix[index_cached] = origin_next_matrix[to_stop_index_cached]
				}
			}
		}

		// solo un hilo
		thread_id_channel <- thread_index
		weight_rows_channel <- origin_weight_matrix
		path_rows_channel <- origin_next_matrix
	}
}

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
	communication_seconds float64
	start_comm time.Time
	end_comm time.Duration
)


func main() {
	communication_seconds = 0

	// Channels
	// Public channels
	thread_id_channel := make(chan int, Threads)
	stop_vertex_weight_matrix_channel := make(chan []int, Threads)

	// Private channels
	weight_rows_channels := make([]chan []int, Threads)
	path_rows_channels := make([]chan []int, Threads)
	for i := range weight_rows_channels {
		weight_rows_channels[i] = make(chan []int, 1)
		path_rows_channels[i] = make(chan []int, 1)
	}

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
	for thread_index := 1; thread_index < Threads; thread_index++ {
		slice_start := (thread_index * rows_per_thread) * MatrixSize
		slice_end := slice_start + MatrixSize * rows_per_thread

		go process_rows(
			rows_per_thread,
			thread_index,
			path_weight_matrix[slice_start:slice_end],
			next_vertex_matrix[slice_start:slice_end],
			thread_id_channel,
			path_rows_channels[thread_index],
			weight_rows_channels[thread_index],
			stop_vertex_channel,
			stop_vertex_weight_matrix_channel,
		)
	}


	for stop_vertex = 0; stop_vertex < MatrixSize; stop_vertex++ {
		from_stop_index_cached = stop_vertex * MatrixSize

		for thread_index := 1; thread_index < Threads; thread_index++ {
			stop_vertex_weight_matrix_channel <- path_weight_matrix[from_stop_index_cached:from_stop_index_cached + MatrixSize]
		}

		// Calculations main thread
		slice_start := (MainThreadIndex * rows_per_thread) * MatrixSize
		slice_end := slice_start + MatrixSize * rows_per_thread
		stop_weight_matrix := path_weight_matrix[from_stop_index_cached:from_stop_index_cached + MatrixSize]
		origin_weight_matrix := path_weight_matrix[slice_start:slice_end]
		origin_next_matrix := next_vertex_matrix[slice_start:slice_end]

		for origin_vertex = 0; origin_vertex < rows_per_thread; origin_vertex++ {
			origin_index_cached = origin_vertex * MatrixSize
			for target_vertex = 0; target_vertex < MatrixSize; target_vertex++ {
				index_cached = origin_index_cached + target_vertex
				to_stop_index_cached = origin_index_cached + stop_vertex

				// if there is no path to the stop or from the stop go to next iteration
				path_weight_from_stop = stop_weight_matrix[target_vertex]
				path_weight_to_stop = origin_weight_matrix[to_stop_index_cached]
				if path_weight_to_stop == Infinity || path_weight_from_stop == Infinity {
					continue
				}

				// fmt.Println(index_cached, origin_weight_matrix, stop_weight_matrix)
				previous_path = origin_weight_matrix[index_cached]
				current_path = path_weight_to_stop + path_weight_from_stop

				if current_path < previous_path {
					origin_weight_matrix[target_vertex] = current_path
					origin_next_matrix[index_cached] = origin_next_matrix[to_stop_index_cached]
				}
			}
		}

		copy(next_vertex_matrix[slice_start:slice_end], origin_next_matrix)
		copy(path_weight_matrix[slice_start:slice_end], origin_weight_matrix)
		// MAIN THREAD

		// Fetch new calculations
		for i := 1; i < Threads; i++ {
			start_comm = time.Now()
			thread_id := <- thread_id_channel
			weight_rows := <- weight_rows_channels[thread_id]
			path_rows := <- path_rows_channels[thread_id]
			end_comm = time.Since(start_comm)
			communication_seconds += end_comm.Seconds()

			slice_start := (thread_id * rows_per_thread) * MatrixSize
			slice_end := slice_start + MatrixSize * rows_per_thread

			copy(next_vertex_matrix[slice_start:slice_end], path_rows)
			copy(path_weight_matrix[slice_start:slice_end], weight_rows)
		}
	}

	//
	// Return values  ~>D
	//
	total_time = time.Since(start)

	// fmt.Println()
	// fmt.Println("Input graph")
	// printMatrix(graph_matrix)

	// fmt.Println()
	// fmt.Println("Path weight matrix")
	// printMatrix(path_weight_matrix)

	// fmt.Println()
	// fmt.Println("Next vertex matrix")
	// printMatrix(next_vertex_matrix)

	fmt.Println()
	fmt.Println("Time:")
	fmt.Println(total_time.Seconds())
	fmt.Println()
	fmt.Println("Communication Time:")
	fmt.Println(communication_seconds)
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