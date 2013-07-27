
var establish_connection = function (update_callback, draw_callback){

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
            establish_connection(update_callback, draw_callback);
        }, 5000);
    };

    connection.onmessage = function (e) {
        // console.log(e.data);
        new_data = JSON.parse(e.data);

        for (var i=0; i < new_data.length; i++){
            // console.log(new_data[i]);
            update_callback(new_data[i], i);
            draw_callback(new_data[i], i);
        }
    };

    return connection;

};

