#### \[WIP\]

### How to run
1. Install go
2. Go to root path of this proyect
3. Run `go build secuencial/main.go` or `go build --race concurrent/main.go`
4. Execute `./main`


#### Reconstruct path
It's possible to use the `next_vertice_matrix` matrix to reconstruct the minimum path between two vertices.
Here is the pseudo-code:
```go
start_vertex := x
end_vertex := y
path := []

if next_vertice_matrix[start_vertex][end_vertex] == nil {
  return path;
}

path.append(start_vertex)
current_step := start_vertex
for(current_step != end_vertex) {
  current_step = next[current_step][end_vertex]
  path.append(current_step)
}

return path
```




## TODO:

* Test performance using slices for the matrixes
  * https://stackoverflow.com/questions/39804861/what-is-a-concise-way-to-create-a-2d-slice-in-go
* Optimize sequential version
* Implement concurrent version