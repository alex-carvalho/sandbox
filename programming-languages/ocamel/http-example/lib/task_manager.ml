type task = {
  id: int;
  title: string;
  completed: bool;
}

let tasks = ref []
let next_id = ref 1

let reset () =
  tasks := [];
  next_id := 1

let create_task title =
  let task = { id = !next_id; title; completed = false } in
  incr next_id;
  tasks := task :: !tasks;
  task

let get_all_tasks () = List.rev !tasks

let complete_task id =
  tasks := List.map (fun t -> if t.id = id then { t with completed = true } else t) !tasks;
  List.find_opt (fun t -> t.id = id) !tasks

let delete_task id =
  let old_tasks = !tasks in
  tasks := List.filter (fun t -> t.id <> id) !tasks;
  List.length old_tasks <> List.length !tasks

let task_to_json task =
  `Assoc [
    ("id", `Int task.id);
    ("title", `String task.title);
    ("completed", `Bool task.completed)
  ]

let tasks_to_json tasks =
  `List (List.map task_to_json tasks)

let json_to_task json =
  match json with
  | `Assoc assoc ->
    let title = match List.assoc "title" assoc with
      | `String s -> s
      | _ -> failwith "Invalid title"
    in
    Some title
  | _ -> None