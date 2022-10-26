package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	diff "m/difftools/diffusion"
	opt "m/difftools/optimization"
	"os"
	"strings"
	"time"
)

type Parameter struct {
	GraphPath        string
	Node_n           int
	Random_seed      int64
	K_F              int
	K_T              int
	Mont_sample_size int
	Prob_map         [2][2][2][2]float64
	Seq              [16]float64
	Pop_list         [2]int
	SeedSet_F        []int
	Interest_list    [][]int
	Assum_list       [][]int
}

func sample1() {
	var n int = 50
	var seed int64 = 1
	var K_F int = 5
	var K_T int = 10
	var sample_size int = 1000
	var pop_list [2]int
	pop_list[0] = diff.Pop_high
	pop_list[1] = diff.Pop_high

	fmt.Println(K_T, K_F, diff.InfoType_F, sample_size, pop_list)
	adjFilePath := "adj_json50node.txt"
	bytes, err := ioutil.ReadFile(adjFilePath)
	if err != nil {
		panic(err)
	}

	// fmt.Println(string(bytes))

	var dataJson string = string(bytes)

	arr := make(map[int]map[int]int)
	// var arr []string
	_ = json.Unmarshal([]byte(dataJson), &arr)
	// fmt.Println(arr)

	// fmt.Println(arr[0][1])

	var adj [][]int = make([][]int, n)

	for i := 0; i < n; i++ {
		adj[i] = make([]int, n)
		for j := 0; j < n; j++ {
			adj[i][j] = arr[j][i]
		}
	}

	// fmt.Println(adj)

	var SeedSet_F []int = diff.Make_seedSet_F(n, 1, seed, adj)

	var interest_list [][]int = diff.Make_interest_list(n, seed)

	var assum_list [][]int = diff.Make_assum_list(n, seed)

	var seq [16]float64 = diff.Make_probability()

	var prob_map [2][2][2][2]float64 = diff.Map_probagbility(seq)

	fmt.Println("Seedsetf")
	fmt.Println(SeedSet_F)

	// fmt.Println((seq))
	//
	// fmt.Println(prob_map)
	seedsetf_pram := make([]int, 0, n)
	for i := 0; i < n; i++ {
		if SeedSet_F[i] == 1 {
			seedsetf_pram = append(seedsetf_pram, i)
		}
	}
	pram_data := Parameter{GraphPath: adjFilePath, Node_n: n, Random_seed: seed, K_F: K_F, K_T: K_T, Mont_sample_size: sample_size, Pop_list: pop_list, SeedSet_F: seedsetf_pram, Interest_list: interest_list, Assum_list: assum_list, Seq: seq, Prob_map: prob_map}

	fmt.Println(pram_data)

	// jsonData, err := json.Marshal(pram_data)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// pram_file, err := os.Create("result/param_json.txt")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//  _, err1 := pram_file.WriteString(jsonData)
	//  if err1 != nil {
	// 	log.Fatal(err)
	// }
	const layout2 = "2006-01-02 15:04:05"
	str := strings.Replace(time.Now().Format(layout2), ":", "-", -1)
	folder_path := "result/" + str
	err = os.Mkdir(folder_path, os.ModePerm)
	if err != nil {
		fmt.Println("error create" + folder_path)
		log.Fatal(err)
	}

	file, _ := json.MarshalIndent(pram_data, "", " ")
	_ = ioutil.WriteFile(folder_path+"/param_json.txt", file, 0644)
	//test part
	var S []int
	var hist [][]float64
	//

	SeedSet_F_strong := make([]int, len(adj))
	SeedSet_F_strong[0] = 1 //here
	// SeedSet_Greedy[1] = 1 //here
	//偽情報の発信源を色々と

	greedy_ans1, greedy_value1, greedy_value21 := opt.Greedy(1, 100, adj, SeedSet_F_strong, prob_map, pop_list, interest_list, assum_list, 3, true, 100)
	fmt.Println(greedy_ans1,greedy_value1,greedy_value21)
	os.Exit(0)


	sample_size = 1000
	S, hist = sim_submod(adj, sample_size, pop_list, interest_list, assum_list, SeedSet_F_strong, K_T, prob_map, folder_path)

	fmt.Println("End Check_submod")

	// new_folder_path := folder_path + "/FocusLoop"
	// err = os.Mkdir(new_folder_path, os.ModePerm)
	// if err != nil {
	// 	fmt.Println("error create FocusLoop")
	// 	log.Fatal(err)
	// }
	// //[0 1 2 6 8 15 18 20 21 22 28 32 33 37 48 49 61 67 76 93]
	// loop_n := 1000
	// sample_size = 1000
	// list1 := []int{0, 6, 8}
	// list2 := []int{0, 20}
	// opt.FocusLoop(loop_n, list1, list2, SeedSet_F, 1, sample_size, adj, prob_map, pop_list, interest_list, assum_list, new_folder_path)

	// time.Sleep(time.Second * 2)
	// list1 = []int{68, 48, 2, 8}
	// list2 = []int{8, 48, 37, 93}
	// opt.FocusLoop(loop_n, list1, list2, SeedSet_F, 1, sample_size, adj, prob_map, pop_list, interest_list, assum_list, new_folder_path)

	// time.Sleep(time.Second * 2)
	// list1 = []int{20, 15, 6}
	// list2 = []int{48, 0, 18}
	// opt.FocusLoop(loop_n, list1, list2, SeedSet_F, 1, sample_size, adj, prob_map, pop_list, interest_list, assum_list, new_folder_path)

	// fmt.Println("End FocusLoop")

	filename := folder_path + "/GreedyAndStrict2.csv"
	f, err2 := os.Create(filename)
	if err2 != nil {
		fmt.Println("error create" + filename)
		log.Fatal(err)
	}

	w := csv.NewWriter(f)

	colmns := []string{"greedy_ans", "strict_ans", "greedy_value", "greedy_value2", "strict_value", "strict_value2", "kinjiritu", "random_seed"}
	w.Write(colmns)

	//start loop
	sample_size = 1000
	sample_size2 := 1000
	var random_seed int64
	random_seed = 0
	_ = random_seed

	//選ばれうる0 1 2 6 8 15 18 20 37 48

	seedsetfs := []int{0, 1, 2, 6}
	// var seedsetfs []int = make([]int, len(diff.Set))
	// _ = copy(seedsetfs, diff.Set)
	// fmt.Println(seedsetfs)
	// os.Exit(0)
	fmt.Println((len(adj)))
	for i := 0; i < 2; i++ {
		SeedSet_Greedy := make([]int, len(adj))
		SeedSet_Greedy[seedsetfs[i]] = 1 //here
		//偽情報の発信源を色々と
		for num := 0; num < 10; num++ {
			greedy_ans, greedy_value, greedy_value2 := opt.Greedy(1, sample_size, adj, SeedSet_Greedy, prob_map, pop_list, interest_list, assum_list, num, true, sample_size2)

			fmt.Println("greedy_ans")
			fmt.Println(greedy_ans, greedy_value, greedy_value2)

			// strict_ans, strict_value, strict_value2 := opt.Strict(random_seed, sample_size, adj, SeedSet_Greedy, prob_map, pop_list, interest_list, assum_list, 3, true, sample_size2)

			// fmt.Println("strict_ans")
			// fmt.Println(strict_ans, strict_value, strict_value2)
			//
			// fmt.Println("近似率")
			// fmt.Println(greedy_value2 / strict_value2)

			// Sets_string := make([][]string, 2)
			Sets_string := make([][]string, 1)
			Sets_string[0] = opt.Int_to_String(greedy_ans)
			// Sets_string[1] = opt.Int_to_String(strict_ans)

			// part0 := []string{strings.Join(Sets_string[0], "-"), strings.Join(Sets_string[1], "-")} //here
			part0 := []string{strings.Join(Sets_string[0], "-")} //here

			// a := []float64{greedy_value, greedy_value2, strict_value, strict_value2, greedy_value2 / strict_value2, float64(random_seed)}
			a := []float64{greedy_value, greedy_value2}

			part1 := opt.Float_to_String(a)

			retu := append(part0, part1...)

			w.Write(retu)
		}
	}

	//loop end

	w.Flush()

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
	//SeedSet_Greedy := make([]int,len(adj))
	//SeedSet_Greedy[15] = 2

	fmt.Println(S, hist)
}

func sim_submod(adj [][]int, sample_size int, pop_list [2]int, interest_list [][]int, assum_list [][]int, SeedSet_F []int, K_T int, prob_map [2][2][2][2]float64, folder_path string) ([]int, [][]float64) {
	var S []int
	var hist [][]float64
	S, hist = opt.Check_submod(1, K_T, sample_size, adj, SeedSet_F, prob_map, pop_list, interest_list, assum_list, folder_path)

	return S, hist
}

func main() {
	sample1()
}
