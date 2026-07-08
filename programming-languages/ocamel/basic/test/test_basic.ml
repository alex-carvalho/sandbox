open Basic

let test_add () =
  assert (Functions.add 3 4 = 7);
  assert (Functions.add 0 0 = 0);
  assert (Functions.add (-1) 1 = 0)

let test_multiply () =
  assert (Functions.multiply 5 6 = 30);
  assert (Functions.multiply 0 10 = 0)

let test_factorial () =
  assert (Functions.factorial 0 = 1);
  assert (Functions.factorial 1 = 1);
  assert (Functions.factorial 5 = 120)

let test_is_even () =
  assert (Functions.is_even 0 = true);
  assert (Functions.is_even 2 = true);
  assert (Functions.is_even 3 = false)

let test_sum_list () =
  assert (Functions.sum_list [] = 0);
  assert (Functions.sum_list [1; 2; 3] = 6)

let test_variables () =
  assert (Variables.my_int = 42);
  assert (Variables.my_float = 3.14);
  assert (Variables.my_string = "Hello, OCaml!")

let () =
  test_add ();
  test_multiply ();
  test_factorial ();
  test_is_even ();
  test_sum_list ();
  test_variables ();
  print_endline "All tests passed!"