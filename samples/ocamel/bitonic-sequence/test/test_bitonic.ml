open OUnit2
open Bitonic

let test_empty _ =
  assert_equal [] (generate_bitonic 0 1 10)

let test_single _ =
  assert_equal [1] (generate_bitonic 1 1 10)

let test_basic_sequence _ =
  let seq = generate_bitonic 5 1 10 in
  assert_equal 5 (List.length seq);
  assert_bool "should be bitonic" (is_bitonic seq)

let test_even_length _ =
  let seq = generate_bitonic 6 0 100 in
  assert_equal 6 (List.length seq);
  assert_bool "should be bitonic" (is_bitonic seq)

let test_range_bounds _ =
  let seq = generate_bitonic 7 5 15 in
  assert_bool "min in range" (List.hd seq >= 5);
  assert_bool "max in range" (List.fold_left max min_int seq <= 15)

let test_is_bitonic_valid _ =
  assert_bool "ascending then descending" (is_bitonic [1; 3; 5; 4; 2]);
  assert_bool "only ascending" (is_bitonic [1; 2; 3; 4]);
  assert_bool "only descending" (is_bitonic [5; 4; 3; 2])

let test_is_bitonic_invalid _ =
  assert_bool "not bitonic" (not (is_bitonic [1; 3; 2; 4]))

let suite =
  "BitonicTests" >::: [
    "test_empty" >:: test_empty;
    "test_single" >:: test_single;
    "test_basic_sequence" >:: test_basic_sequence;
    "test_even_length" >:: test_even_length;
    "test_range_bounds" >:: test_range_bounds;
    "test_is_bitonic_valid" >:: test_is_bitonic_valid;
    "test_is_bitonic_invalid" >:: test_is_bitonic_invalid;
  ]

let () = run_test_tt_main suite
