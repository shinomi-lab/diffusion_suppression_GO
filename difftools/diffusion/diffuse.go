package diffusion

import (
	"m/difftools/funcs"
	"math/rand"
)

func Adjmat(adj [][]int, SeedSet []int, seed int64, prob_map [2][2][2][2]float64, pop_list [2]int, interest_list [][]int, assum_list [][]int) [][]int {
	//return value is infoごとに受け取ったノードたち(index)
	var n int
	n = len(adj)
	recieved_list := make([][]int, InfoTypes_n)
	for i := 0; i < InfoTypes_n; i++ {
		recieved_list[i] = make([]int, 0, n)
	}

	if seed != -1 {
		rand.Seed(seed)
	}

	current := make([][]int, InfoTypes_n)
	for i := 0; i < InfoTypes_n; i++ {
		current[i] = make([]int, 0, n)
	}
	var infotypes []int = []int{InfoType_F, InfoType_T}
	for j := 0; j < n; j++ {
		for _, info := range infotypes {
			if SeedSet[j] == info+1 {
				current[info] = append(current[info], j)
				recieved_list[info] = append(recieved_list[info], j)
			}

		}
	}

	//main loop
	for len(current[InfoType_F]) > 0 || len(current[InfoType_T]) > 0 {
		next := make([][]int, InfoTypes_n)
		for info, set := range current {
			for _, s_node := range set {
				pop := pop_list[info]
				interest := interest_list[s_node][pop]
				assum := assum_list[s_node][info]
				p := prob_map[pop][info][interest][assum]

				for j := 0; j < n; j++ {
					if adj[s_node][j] == 0 || funcs.Set_Has(recieved_list[InfoType_F], j) || funcs.Set_Has(recieved_list[InfoType_T], j) || funcs.Set_Has(next[InfoType_F], j) || funcs.Set_Has(next[info], j) {
						//道がないorすでに情報を受け取っているor次に偽の情報または同じ種類の情報を受け取ろうとしている
						continue
					}
					if p == 1 || p > rand.Float64() {
						next[info] = append(next[info], j)
						if info == InfoType_F && funcs.Set_Has(next[InfoType_F], j) {
							remove(next[InfoType_T], j)
						}
					}
				}
			}
		}
		current = make([][]int, InfoTypes_n)
		_ = copy(current, next)
		for _, info := range infotypes {
			recieved_list[info] = funcs.Set_Sum(recieved_list[info], next[info])
		}
	}
	return recieved_list
}

func remove(ints []int, search int) []int {
	result := []int{}
	for _, v := range ints {
		if v != search {
			result = append(result, v)
		}
	}
	return result
}
