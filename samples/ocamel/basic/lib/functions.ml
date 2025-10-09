(* Basic OCaml function examples *)

(* Simple function with two parameters *)
let add x y = x + y

(* Function with type annotations *)
let multiply (x : int) (y : int) : int = x * y

(* Anonymous function (lambda) *)
let double = fun x -> x * 2

(* Recursive function *)
let rec factorial n =
  if n <= 1 then 1
  else n * factorial (n - 1)

(* Mutual recursive functions *)
let rec is_even n =
  if n = 0 then true
  else is_odd (n - 1)
and is_odd n =
  if n = 0 then false
  else is_even (n - 1)

(* Pattern matching function *)
let describe_number = function
  | 0 -> "zero"
  | 1 -> "one"
  | 2 -> "two"
  | n when n < 0 -> "negative"
  | _ -> "many"

(* Function with pattern matching on lists *)
let rec sum_list = function
  | [] -> 0
  | head :: tail -> head + sum_list tail

(* Higher-order function that takes a function as parameter *)
let apply_twice f x = f (f x)

(* Function composition operator *)
let (>>) f g x = g (f x)
let add_one x = x + 1
let multiply_by_two x = x * 2
let composed = add_one >> multiply_by_two

(* Partial application *)
let add_five = add 5

(* Function that returns a function (currying) *)
let make_adder x = fun y -> x + y

(* Function with optional parameter *)
let greet ?(greeting="Hello") name =
  Printf.printf "%s, %s!\n" greeting name

(* Function with labeled parameters *)
let divide ~numerator ~denominator =
  if denominator = 0 then None
  else Some (float_of_int numerator /. float_of_int denominator)

