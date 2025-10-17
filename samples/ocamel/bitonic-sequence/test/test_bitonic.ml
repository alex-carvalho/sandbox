open OUnit2
open Bitonic

let test_sequence1 _ =
  assert_equal [9; 10; 9; 8; 7] (generate_bitonic 5 3 10)

let test_sequence2 _ =
  assert_equal [2; 3; 4; 5; 4; 3; 2] (generate_bitonic 7 2 5)

let test_sequence3 _ =
  assert_equal [3; 4; 5; 4; 3; 2] (generate_bitonic 6 2 5)

let test_sequence4 _ =
  assert_equal [9; 10; 9] (generate_bitonic 3 1 10)

let suite =
  "BitonicTests" >::: [
    "test_sequence1" >:: test_sequence1;
    "test_sequence2" >:: test_sequence2;
    "test_sequence3" >:: test_sequence3;
    "test_sequence4" >:: test_sequence4;
  ]

let () = run_test_tt_main suite
