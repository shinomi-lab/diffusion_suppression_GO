package optimization

import(
  diff "m/difftools/diffusion"
)

func Greedy(seed int64,sample_size int , adj [][]int , Seed_set []int, prob_map [2][2][2][2]float64, pop [2]int, interest_list [][]int, assum_list [][]int, ans_len int, Count_true bool,sample_size2 int)([]int,float64,float64){
  //sample_size2はグリーディで求めた買いをより詳しくやる
  var n int = len(adj)
  var max float64 = 0
  var result float64
  var index int
  var ans []int

  ans = make([]int, 0, ans_len)
  S := make([]int, len(Seed_set))
  _ = copy(S,Seed_set)
  S_test := make([]int ,len(Seed_set))
  _ = copy(S_test, Seed_set)

  var info_num int

  if(Count_true){
    info_num = 2
  }else{
    info_num = 1
  }

  for i:=0; i<ans_len; i++{
    max = 0
    for j:=0; j<n; j++{
      _ = copy(S_test,S)
      if(S_test[j] != 0){//すでに発信源のユーザだったら
        continue
      }
      S_test[j] = info_num

      dist := Infl_prop_exp(seed, sample_size, adj, S_test, prob_map, pop, interest_list, assum_list)
      if (Count_true){
        result = dist[diff.InfoType_T]
      }else{
        result = dist[diff.InfoType_F]
      }

      if (result > max){
        max = result
        index = j
      }
    }//subloop end

    ans = append(ans, index)
    S[index] = info_num

  }//mainloop end

  var max_2 float64
  dist2 := Infl_prop_exp(seed, sample_size2, adj, S, prob_map, pop, interest_list, assum_list)
  if (Count_true){
    max_2 = dist2[diff.InfoType_T]
  }else{
    max_2 = dist2[diff.InfoType_F]
  }
  return ans,max,max_2
}
