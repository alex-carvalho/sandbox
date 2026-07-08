(* Basic OCaml variable declarations *)

(* Integer *)
let my_int = 42

(* Float *)
let my_float = 3.14
let scientific = 2.5e-3

(* String *)
let my_string = "Hello, OCaml!"
let multiline_string = "This is a
multiline string"

(* Boolean *)
let my_bool = true
let another_bool = false

(* Character *)
let my_char = 'A'

(* List *)
let int_list = [1; 2; 3; 4; 5]
let string_list = ["apple"; "banana"; "orange"]

(* Tuple *)
let my_tuple = (42, "OCaml", true)
(* Access tuple elements *)
let first, second, third = my_tuple
(* first = 42, second = "OCaml", third = true *)

(* Option type *)
let maybe_int = Some 42
let no_value = None

(* Unit type *)
let unit_value = ()

(* Compound declaration *)
let (x, y) = (10, 20)

(* Function type *)
let add a b = a + b

