let redis_conn = Redis_lwt.Client.connect { host = "127.0.0.1"; port = 6379 }

let () =
  Dream.run
  @@ Dream.logger
  @@ Dream.router [
    Dream.post "/bitonic" (fun request ->
      let%lwt body = Dream.body request in
      try
        let json = Yojson.Safe.from_string body in
        let open Yojson.Safe.Util in
        let n = json |> member "n" |> to_int in
        let min_val = json |> member "min" |> to_int in
        let max_val = json |> member "max" |> to_int in
        
        if n < 0 || min_val > max_val then
          Dream.json ~status:`Bad_Request 
            {|{"error":"Invalid parameters"}|}
        else
          let sequence = Bitonic.generate_bitonic n min_val max_val in
          let result = `Assoc [
            ("sequence", `List (List.map (fun x -> `Int x) sequence));
            ("length", `Int (List.length sequence));
            ("is_bitonic", `Bool (sequence <> [-1]))
          ] in
          let result_str = Yojson.Safe.to_string result in
          let key = Printf.sprintf "bitonic:%d:%d:%d" n min_val max_val in
          let%lwt conn = redis_conn in
          let%lwt _ = Redis_lwt.Client.set conn key result_str in
          let%lwt _ = Redis_lwt.Client.expire conn key 3600 in
          Dream.json result_str
      with _ ->
        Dream.json ~status:`Bad_Request 
          {|{"error":"Invalid JSON format"}|}
    );
    
    Dream.get "/health" (fun _ ->
      Dream.json {|{"status":"ok"}|}
    );
  ]
