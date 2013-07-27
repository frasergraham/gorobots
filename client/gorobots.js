var id = null;

var establish_connection = function (update_callback, draw_callback){

    var websocket_server = "ws://twisted.local:8666/ws/";
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
        // console.log(e.data);
        new_data = JSON.parse(e.data);

        if ('id' in new_data){
            id = new_data['id'];
            console.log("Got ID: " + id);
        }

        var players = "";
        var robots = new_data['robots'];

        for (var i=0; i < robots.length; i++){
            players += "\n" + robots[i]['id'];
            // console.log(robots[i]);
            if ("position" in robots[i]){
               draw_callback(robots[i], i);
            }
            if (robots[i]['id'] == id)
                update_callback(robots[i], i);
        }

        var projectiles = new_data['projectiles'];

        for (i=0; i < projectiles.length; i++){
            // console.log(projectiles[i]);
            if ("position" in projectiles[i]){
                projectiles[i]['type'] = 'bullet';
                draw_callback(projectiles[i], i+1);
            }
        }
        var players_div = document.getElementById("players");
        players_div.innerHTML = players;
    };

    return connection;

};


var websocket;

var colors = [
    "rgba(0, 0, 0, 0.5)",
    "rgba(0, 0, 200, 0.5)",
    "rgba(0, 200, 0, 0.5)",
    "rgba(200, 0, 0, 0.5)",
    "rgba(0, 200, 200, 0.5)",
    "rgba(200, 0, 200, 0.5)",
    "rgba(200, 200, 0, 0.5)",
    "rgba(0, 0, 150, 0.5)",
    "rgba(0, 150, 0, 0.5)",
    "rgba(150, 0, 0, 0.5)",
    "rgba(0, 150, 150, 0.5)",
    "rgba(150, 0, 150, 0.5)",
    "rgba(150, 150, 0, 0.5)",
    "rgba(0, 0, 50, 0.5)",
    "rgba(0, 50, 0, 0.5)",
    "rgba(50, 0, 0, 0.5)",
    "rgba(0, 50, 50, 0.5)",
    "rgba(50, 0, 50, 0.5)",
    "rgba(50, 50, 0, 0.5)",
    "rgba(0, 0, 222, 0.5)",
    "rgba(0, 222, 0, 0.5)",
    "rgba(222, 0, 0, 0.5)",
    "rgba(0, 222, 222, 0.5)",
    "rgba(222, 0, 222, 0.5)",
    "rgba(222, 222, 0, 0.5)"
];

// This is going to be the main module
var gorobots = function(my){
    my.width = 800;
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
            ctx.clearRect ( 0 , 0 , 800 , 350 );

        if ('type' in data && data['type'] == 'bullet'){
            ctx.fillStyle = colors[0];
            ctx.fillRect (x, y, 5, 5);
        }
        else if ('type' in data && data['type'] == 'explosion'){
            ctx.beginPath();
            ctx.arc(x, y, 40, 0, 2 * Math.PI, false);
            ctx.fillStyle = colors[index+1];
            ctx.fill();
        }
        else{
            ctx.fillStyle = colors[index+1];
            ctx.fillRect (x, y, 10, 10);

            if (id == data['id'])
                ctx.font="bold 22px Helvetica";
            else
                ctx.font="10px Helvetica";

            ctx.fillText(
                data['id'] + "[" + data['health'] + "]",
                x+12,y+10);
        }
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
    var map = {"width": 800, "height": 350};
    var out = rc.call(this, data, map);

    out['id'] = id;

    if (websocket){
        websocket.send(JSON.stringify(out));
        //console.log("SENDING: " + JSON.stringify(out));
    }
};


function init(){
    websocket = establish_connection(eval_robot, draw);
}

