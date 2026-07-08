open Alcotest
open Task_manager

let test_create_task () =
  reset ();
  let task = create_task "Test task" in
  check int "task id" 1 task.id;
  check string "task title" "Test task" task.title;
  check bool "task completed" false task.completed

let test_get_all_tasks () =
  reset ();
  let _ = create_task "Task 1" in
  let _ = create_task "Task 2" in
  let tasks = get_all_tasks () in
  check int "tasks count" 2 (List.length tasks)

let test_complete_task () =
  reset ();
  let task = create_task "Complete me" in
  let completed = complete_task task.id in
  match completed with
  | Some t -> check bool "task completed" true t.completed
  | None -> fail "Task not found"

let test_delete_task () =
  reset ();
  let task = create_task "Delete me" in
  let deleted = delete_task task.id in
  check bool "task deleted" true deleted;
  let not_found = delete_task 999 in
  check bool "non-existent task" false not_found

let test_json_conversion () =
  reset ();
  let task = { id = 1; title = "JSON test"; completed = false } in
  let json = task_to_json task in
  let expected = `Assoc [
    ("id", `Int 1);
    ("title", `String "JSON test");
    ("completed", `Bool false)
  ] in
  check (testable Yojson.Safe.pp Yojson.Safe.equal) "json conversion" expected json

let () =
  run "Task Manager Tests" [
    "create_task", [ test_case "Create task" `Quick test_create_task ];
    "get_all_tasks", [ test_case "Get all tasks" `Quick test_get_all_tasks ];
    "complete_task", [ test_case "Complete task" `Quick test_complete_task ];
    "delete_task", [ test_case "Delete task" `Quick test_delete_task ];
    "json_conversion", [ test_case "JSON conversion" `Quick test_json_conversion ];
  ]