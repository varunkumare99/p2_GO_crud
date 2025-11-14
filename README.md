1 - Add a /health endpoint returning { "status": "ok" }.

2 - Initialize todos with 2 default items.

3 - Add filtering: GET /todos?done=true should return only completed todos.

4 - Add validation: if title is empty in POST, return 400.

5 - Add middleware to log request time for every API call.
