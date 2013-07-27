var id = null;

var establish_connection = function (update_callback, draw_callback){

    var websocket_server = "ws://midna.local:8666/ws/";
    connection = new WebSocket(websocket_server, null);

    connection.onerror = function (error) {
      console.log('WebSocket Error ' + error);
    };

    connection.onopen = function(){
        console.log("Connected to " + websocket_server);
        connection.send(JSON.stringify({'move_to':{"x": 0, "y":0}}));
    };

    connection.onclose = function(){
        console.log("Lost Connection: " + websocket_server);
        id = null;
        setTimeout(function(){
            websocket = establish_connection(update_callback, draw_callback);
        }, 5000);
    };

    connection.onmessage = function (e) {
        console.log(e.data);
        new_data = JSON.parse(e.data);

        if ('id' in new_data){
            id = new_data['id'];
            console.log("Got ID: " + id);
        }

        for (var i=0; i < new_data.length; i++){
            // console.log(new_data[i]);
            if ("position" in new_data[i]){
               draw_callback(new_data[i], i);
            }
            if (new_data[i]['id'] == id)
                update_callback(new_data[i], i);
        }
    };

    return connection;

};


var websocket;

var colors = [
    "rgba(0, 0, 200, 0.5)",
    "rgba(0, 200, 0, 0.5)",
    "rgba(200, 0, 0, 0.5)",
    "rgba(0, 200, 200, 0.5)",
    "rgba(200, 0, 200, 0.5)",
    "rgba(200, 200, 0, 0.5)"
];

// This is going to be the main module
var gorobots = function(my){
    my.width = 450;
    my.height = 400;

    return my;
}({});


var draw = function(data, index){

    var canvas = document.getElementById('tutorial');

    if (canvas.getContext){
        var ctx = canvas.getContext('2d');
        var x = data['position']['x'];
        var y = data['position']['y'];

        if (index === 0)
            ctx.clearRect ( 0 , 0 , 450 , 350 );

        ctx.fillStyle = colors[index];
        ctx.fillRect (x, y, 10, 10);

        // ctx.fillStyle = "blue";
        // ctx.font = "bold 16px Arial";
        // ctx.fillText(data['id'], x, y);

    }
};

function evalInput( input, output ){
    var theResult, evalSucceeded;
    try{
        theResult = eval( input );
        evalSucceeded = true;
    }
    catch(e){
        output.innerHTML = e;
    }
    if ( evalSucceeded )
    {
        // output.innerHTML = (theResult+"").replace( /&/g, '&amp;' ).replace( /</g, '&lt;' ).replace( />/g, '&gt;' ).replace( /\r\n|\r|\n/g, '<br>' );
    }

    return theResult;
}


var eval_robot = function(data){
    var robot_code = document.getElementById('robot');
    var output = document.getElementById('output');
    var code = "( " + robot_code.value + " )";
    var rc = evalInput(code, output);
    var out = rc.call(this, data);

    out['id'] = id;

    if (websocket){
        websocket.send(JSON.stringify(out));
        console.log("SENDING: " + JSON.stringify(out));
    }
};


function init(){
    websocket = establish_connection(eval_robot, draw);
}

