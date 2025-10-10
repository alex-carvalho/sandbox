let generate_bitonic n min_val max_val =
  if n <= 0 then []
  else if n = 1 then [min_val]
  else
    let peak_pos = (n + 1) / 2 in
    let increasing_len = peak_pos in
    let decreasing_len = n - peak_pos in
    
    let range = max_val - min_val in
    
    let increasing = 
      List.init increasing_len (fun i ->
        min_val + (range * i) / (increasing_len - 1)
      )
    in
    
    let decreasing = 
      if decreasing_len > 0 then
        List.init decreasing_len (fun i ->
          max_val - (range * (i + 1)) / decreasing_len
        )
      else []
    in
    
    increasing @ decreasing

let is_bitonic seq =
  let rec find_peak i increasing =
    if i >= List.length seq - 1 then true
    else
      let curr = List.nth seq i in
      let next = List.nth seq (i + 1) in
      if increasing then
        if curr < next then find_peak (i + 1) true
        else if curr = next then false
        else find_peak (i + 1) false
      else
        if curr > next then find_peak (i + 1) false
        else false
  in
  match seq with
  | [] | [_] -> true
  | _ -> find_peak 0 true
