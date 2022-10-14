package optimization

import(
  diff "m/difftools/diffusion"
  // "fmt"
)

func Strict(seed int64, sample_size int , adj [][]int , Seed_set []int, prob_map [2][2][2][2]float64, pop [2]int, interest_list [][]int, assum_list [][]int, ans_len int, Count_true bool,sample_size2 int)([]int,float64,float64){
  var n int = len(diff.Set)
  var max float64 = 0
  var result float64
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

  for i:=0; i<n; i++{
    if(S_test[diff.Set[i]] != 0){//すでに発信源のユーザだったら
      continue
    }
    for j:=i;j<n;j++{
      if(S_test[diff.Set[j]] != 0){//すでに発信源のユーザだったら
        continue
      }
      for k:=j;k<n;k++{
        if(S_test[diff.Set[k]] != 0){//すでに発信源のユーザだったら
          continue
        }
        //main loop
        _ = copy(S_test,S)
        S_test[diff.Set[i]] = info_num
        S_test[diff.Set[j]] = info_num
        S_test[diff.Set[k]] = info_num
        //complete set Seedsets

        dist := Infl_prop_exp(seed, sample_size, adj, S_test, prob_map, pop, interest_list, assum_list)
        if (Count_true){
          result = dist[diff.InfoType_T]
        }else{
          result = dist[diff.InfoType_F]
        }
        if (result > max){
          max = result
          ans = []int{diff.Set[i], diff.Set[j], diff.Set[k]}//いける？
        }

      }
    }
  }//mainloop end
  _ = copy(S_test,S)
  for i:=0;i<len(ans);i++{
    S_test[ans[i]] = info_num
  }

  //complete set Seedsets

  var max2 float64
  dist2 := Infl_prop_exp(seed, sample_size2, adj, S_test, prob_map, pop, interest_list, assum_list)
  if (Count_true){
    max2 = dist2[diff.InfoType_T]
  }else{
    max2 = dist2[diff.InfoType_F]
  }

  return ans,max,max2
}
