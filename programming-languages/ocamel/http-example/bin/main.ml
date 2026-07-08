open Lwt.Syntax
open Cohttp_lwt_unix
open Task_manager

let json_header = Cohttp.Header.init_with "content-type" "application/json"

let handle_request req body =
  let uri = Cohttp.Request.uri req in
  let meth = Cohttp.Request.meth req in
  let path = Uri.path uri in
  
  match (meth, path) with
  | (`GET, "/tasks") ->
    let tasks = get_all_tasks () in
    let json = tasks_to_json tasks in
    Server.respond_string ~status:`OK ~headers:json_header ~body:(Yojson.Safe.to_string json) ()
    
  | (`POST, "/tasks") ->
    let* body_string = Cohttp_lwt.Body.to_string body in
    (try
       let json = Yojson.Safe.from_string body_string in
       match json_to_task json with
       | Some title ->
         let task = create_task title in
         let json = task_to_json task in
         Server.respond_string ~status:`Created ~headers:json_header ~body:(Yojson.Safe.to_string json) ()
       | None ->
         Server.respond_string ~status:`Bad_request ~body:"Invalid JSON" ()
     with _ ->
       Server.respond_string ~status:`Bad_request ~body:"Invalid JSON" ())
       
  | (`PUT, path) when String.starts_with ~prefix:"/tasks/" path ->
    let id_str = String.sub path 7 (String.length path - 7) in
    (try
       let id = int_of_string id_str in
       match complete_task id with
       | Some task ->
         let json = task_to_json task in
         Server.respond_string ~status:`OK ~headers:json_header ~body:(Yojson.Safe.to_string json) ()
       | None ->
         Server.respond_string ~status:`Not_found ~body:"Task not found" ()
     with _ ->
       Server.respond_string ~status:`Bad_request ~body:"Invalid task ID" ())
       
  | (`DELETE, path) when String.starts_with ~prefix:"/tasks/" path ->
    let id_str = String.sub path 7 (String.length path - 7) in
    (try
       let id = int_of_string id_str in
       if delete_task id then
         Server.respond_string ~status:`No_content ~body:"" ()
       else
         Server.respond_string ~status:`Not_found ~body:"Task not found" ()
     with _ ->
       Server.respond_string ~status:`Bad_request ~body:"Invalid task ID" ())
       
  | _ ->
    Server.respond_string ~status:`Not_found ~body:"Not found" ()

let start_server port =
  let callback _conn req body = handle_request req body in
  let server = Server.create ~mode:(`TCP (`Port port)) (Server.make ~callback ()) in
  Printf.printf "Server running on http://localhost:%d\n" port;
  server

let () =
  let port = 8080 in
  Lwt_main.run (start_server port)