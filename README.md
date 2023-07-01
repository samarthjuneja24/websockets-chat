This application will help you create a websocket server. You can connect different clients to it and send messages either to all of the connected clients or specific clients

Steps to run this

1. Run `go mod tidy` to install all dependencies
2. Run `go run main.go` to start the application
3. Open three tabs in your browser's console
4. Use the following commands to connect to the server

   `let socket = new WebSocket("ws://localhost:8080/ws?number=123")`
   `socket.onmessage = (event) => console.log(event.data)`
5. Run these commands in all three tabs with different number param
6. Run this command in any one of the tabs to send a message to a specific tab
   `socket.send("{\"message\": \"hello from 123\", \"receiver\": \"456\"}")`

   This is assuming the receiver number in one of the other tabs is `456`
7. To send a message in all the tabs/to all the clients, replace the number in receiver with the keyword `all`