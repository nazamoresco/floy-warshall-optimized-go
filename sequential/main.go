package main

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"
)

const MatrixSize = 3
const Infinity = math.MaxInt

func main() {
	//
	// Initialization of matrixes ~>
	//

	// Matrix with the minimum distances between each pair of vertices
	path_weight_matrix := make([]int, MatrixSize * MatrixSize)
	next_vertex_matrix := make([]int, MatrixSize * MatrixSize)

	// Adjencency matrix representing the graph
	graph_matrix := make([]int, MatrixSize * MatrixSize)

	// Initialize graph
	edges := float64(MatrixSize * (MatrixSize - 1)) * 0.7 // 70% Density
	for row_vertex := 0; row_vertex < MatrixSize; row_vertex++ {
		row_index_cached := row_vertex * MatrixSize
		for column_vertex := 0; column_vertex < MatrixSize; column_vertex ++ {
			cached_index := row_index_cached + column_vertex
			if(edges <= 0) {
				graph_matrix[cached_index] = -1
			} else {
				graph_matrix[cached_index] = rand.Intn(100)
				edges = edges - 1
			}
		}
	}

	//
	// Initialize result matrixes and calculate shortest path  ~>
	//

	start := time.Now()
	for stop_vertex := 0; stop_vertex < MatrixSize; stop_vertex++ {
		from_stop_index_cached := stop_vertex * MatrixSize
		for origin_vertex := 0; origin_vertex < MatrixSize; origin_vertex++ {
			origin_index_cached := origin_vertex * MatrixSize
			for target_vertex := 0; target_vertex < MatrixSize; target_vertex++ {
				index_cached := origin_index_cached + target_vertex
				to_stop_index_cached := origin_index_cached + stop_vertex
				if(stop_vertex == 0) {
					// Each vertex has a distance of zero to itself
					if target_vertex == origin_vertex {
						path_weight_matrix[index_cached] = 0
						next_vertex_matrix[index_cached] = target_vertex
						continue
					}

					// If there is no path the path weight is positive Infinity
					if graph_matrix[index_cached] == -1 {
						path_weight_matrix[index_cached] = Infinity
						next_vertex_matrix[index_cached] = -1
						continue
					}

					// Default minimum path between the two vertices is the direct path
					path_weight_matrix[index_cached] = graph_matrix[index_cached]
					next_vertex_matrix[index_cached] = target_vertex
					continue
				}

				// if there is no path to the stop or from the stop go to next iteration
				path_weight_from_stop := path_weight_matrix[from_stop_index_cached + target_vertex]
				path_weight_to_stop := path_weight_matrix[to_stop_index_cached]
				if path_weight_to_stop == Infinity || path_weight_from_stop == Infinity {
					continue
				}

				previous_path := path_weight_matrix[index_cached]
				current_path := path_weight_to_stop + path_weight_from_stop
				if current_path < previous_path {
					path_weight_matrix[index_cached] = current_path
					next_vertex_matrix[index_cached] = next_vertex_matrix[to_stop_index_cached]
				}
			}
		}
	}


	//
	// Return values  ~>
	//
	total_time := time.Since(start)

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
	fmt.Println(total_time)
}


// printMatrix: Format and prints a given matrix
// Assumes first index is row, second index is column
// Prepared for 4 digit numbers, more digits would break this
func printMatrix(matrix [MatrixSize * MatrixSize]int) {
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