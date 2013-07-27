
var establish_connection = function (callback){

    var websocket_server = "ws://midna.local:8666/ws/";
    connection = new WebSocket(websocket_server, null);

    connection.onerror = function (error) {
      console.log('WebSocket Error ' + error);
    };

    connection.onopen = function(){
        console.log("Connected to " + websocket_server);
    };

    connection.onclose = function(){
        console.log("Lost Connection: " + websocket_server);
        setTimeout(function(){
            // console.log("Trying");
            establish_connection(callback);
        }, 5000);
    };

    connection.onmessage = function (e) {
        new_data = JSON.parse(e.data);
        // console.log(new_data);

        callback(new_data);
    };

};

