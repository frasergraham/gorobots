
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


var websocket;

var colors = [
    "rgba(0, 0, 200, 0.5)",
    "rgba(0, 200, 0, 0.5)",
    "rgba(200, 0, 0, 0.5)",
    "rgba(0, 200, 200, 0.5)",
    "rgba(200, 0, 200, 0.5)",
    "rgba(200, 200, 0, 0.5)"
];

var d = function(data, index){

    var canvas = document.getElementById('tutorial');

    if (canvas.getContext){
        var ctx = canvas.getContext('2d');
        var x = data['position']['x'];
        var y = data['position']['y'];

        if (index === 0)
            ctx.clearRect ( 0 , 0 , 450 , 350 );

        ctx.fillStyle = colors[index];
        ctx.fillRect (x, y, 10, 10);

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

    if (websocket){
        // console.log(out);
        websocket.send(JSON.stringify(out));
    }
};


function init(){
    websocket = establish_connection(eval_robot, d);
}

