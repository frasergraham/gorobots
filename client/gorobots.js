var websocket_server = "midna.local"


var establish_connection = function establish_connection(){

    my_connection.connection = new WebSocket(websocket_server, null);

    my_connection.onerror = function (error) {
      console.log('WebSocket Error ' + error);
    };

    my_connection.connection.onopen = function(){
        console.log("Connected to " + websocket_server);
    };

    my_connection.connection.onclose = function(){
        console.log("Lost Connection: " + websocket_server);
        setTimeout(function(){
            establish_connection();
        }, 5000);
    };

    my_connection.connection.onmessage = function (e) {
        new_data = JSON.parse(e.data);


    };

}();

