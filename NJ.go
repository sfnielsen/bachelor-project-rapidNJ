package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
)

func canonicalNeighborJoining(Q [][]float64, r []float64, D [][]float64, n int) (int, int) {
	cur_val := math.MaxFloat64
	cur_i, cur_j := -1, -1

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {

			if i == j {
				Q[i][j] = 0
			} else {
				Q[i][j] = D[i][j] - r[i] - r[j]
				if Q[i][j] < cur_val {
					cur_val = Q[i][j]
					cur_i = i
					cur_j = j
				}
			}
		}
	}
	return cur_i, cur_j

}

func neighborJoin(D [][]float64, labels []string) (string, Tree) {

	var label_tree Tree = generateTreeForRapidNJ(labels)
	var tree Tree
	tree = append(tree, label_tree...)

	u := create_u(D)

	newick, tree := neighborJoinRec(D, labels, label_tree, tree, u)

	return newick, tree
}

func neighborJoinRec(D [][]float64, labels []string, array Tree, tree Tree, u []float64) (string, Tree) {

	n := len(D)
	Q := make([][]float64, n)
	for i := range Q {
		Q[i] = make([]float64, n)
	}

	cur_i, cur_j := canonicalNeighborJoining(Q, u, D, n)

	if NewickFlag {
		//Distance to new point where they meet
		v_iu := fmt.Sprintf("%f", D[cur_i][cur_j]/2+(u[cur_i]-u[cur_j])/2)
		v_ju := fmt.Sprintf("%f", D[cur_i][cur_j]/2+(u[cur_j]-u[cur_i])/2)
		//convert to string
		distance_to_x, _ := strconv.ParseFloat(v_iu, 64)
		distance_to_y, _ := strconv.ParseFloat(v_ju, 64)
		newNode := integrateNewNode(array[cur_i], array[cur_j], distance_to_x, distance_to_y)
		array[cur_i] = newNode
		tree = append(tree, newNode)
		array = append(array[:cur_j], array[cur_j+1:]...)

		//creating newick form
		labels[cur_i] = "(" + labels[cur_i] + ":" + v_iu + "," + labels[cur_j] + ":" + v_ju + ")"
		labels = append(labels[:cur_j], labels[cur_j+1:]...)
	}

	u = update_u_nj(D, u, cur_i, cur_j)

	D_new := createNewDistanceMatrixNJ(D, cur_i, cur_j)

	//stop maybe
	if len(D_new) > 2 {
		neighborJoinRec(D_new, labels, array, tree, u)
	} else {
		if NewickFlag {
			newick := "(" + labels[0] + ":" + fmt.Sprintf("%f", D_new[0][1]/2) + "," + labels[1] + ":" + fmt.Sprintf("%f", D_new[0][1]/2) + ");"

			err := ioutil.WriteFile("newick.txt", []byte(newick), 0644)
			if err != nil {
				panic(err)
			}

			new_edge_0 := new(Edge)
			new_edge_0.Distance = D_new[0][1]
			new_edge_0.Node = array[1]

			new_edge_1 := new(Edge)
			new_edge_1.Distance = D_new[0][1]
			new_edge_1.Node = array[0]

			array[0].Edge_array = append(array[0].Edge_array, new_edge_0)
			array[1].Edge_array = append(array[1].Edge_array, new_edge_1)

			array = remove(array, 0)

			return newick, tree
		}
	}
	return "error", tree
}

func createNewDistanceMatrixNJ(D [][]float64, p_i int, p_j int) [][]float64 {

	for k := 0; k < len(D); k++ {
		if p_i == k {
			continue
		}
		if p_j == k {
			continue
		} else {
			//Overwrite p_i as merge ij
			temp := (D[p_i][k] + D[p_j][k] - D[p_i][p_j]) / 2
			D[p_i][k] = temp
			D[k][p_i] = temp
		}
	}

	D_new := append(D[:p_j], D[p_j+1:]...)

	for i := 0; i < len(D_new); i++ {
		D_new[i] = append(D_new[i][:p_j], D_new[i][p_j+1:]...)
	}

	return D_new
}

func update_u_nj(D [][]float64, u []float64, i int, j int) []float64 {

	n := len(D)

	//update i with merge ij
	u[i] = 0

	for idx := range u {
		if idx == i || idx == j {
			continue
		}
		u[idx] = u[idx]*float64(n-2) - D[idx][i] - D[idx][j]
		new_dist := (D[i][idx] + D[j][idx] - D[i][j]) / 2.0

		u[idx] += new_dist
		u[idx] /= float64(n - 3)

		//also add value to ij merge
		u[i] += new_dist

	}
	u[i] /= float64(n - 3)

	//remove j from the array
	u = append(u[:j], u[j+1:]...)

	return u
}
