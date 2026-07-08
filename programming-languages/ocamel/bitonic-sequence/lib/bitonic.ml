let generate_bitonic n min_val max_val =
  if n > (max_val - min_val) * 2 + 1 then
    [-1]
  else
    let front = ref [max_val - 1] in
    let back_rev = ref [] in
    let size () = List.length !front + List.length !back_rev in
    (* Add decreasing part: for i = max_val downto min_val add to back (addLast) *)
    for i = max_val downto min_val do
      if size () < n then back_rev := i :: !back_rev
    done;
    (* Add increasing part: for i = max_val - 2 downto min_val add to front (addFirst) *)
    for i = max_val - 2 downto min_val do
      if size () < n then front := i :: !front
    done;
    (* concatenate front with reversed back_rev to produce deque contents *)
    !front @ (List.rev !back_rev)
