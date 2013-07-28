function init(){(function gorobots(my){
    my.width = 800;
    my.height = 400;
    my.server = "ws://twisted.local:8666/ws/";
    my.server = "ws://rs.mcquay.me:8666/ws/";
    my.websocket = null;
    my.id = null;

    my.debug = false;
    my.debug_draw = true;

    var colors = [
        "rgba(0, 0, 0, 0.5)",       // Grey
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

    my.connect = function(server){

        connection = new WebSocket(server, null);

        connection.onerror = function (error) {
          console.log('WebSocket Error ' + error);
        };

        connection.onopen = function(){
            console.log("Connected to " + server);

            // Handshake with the server
            connection.send(JSON.stringify({'move_to':{"x": 0, "y":0}}));
        };

        connection.onclose = function(){
            console.log("Lost Connection: " + server);
            my.id = null;

            // Retry every few seconds
            setTimeout(function(){
                my.websocket = my.connect(server);
            }, 3000);
        };

        connection.onmessage = function (e) {
            // console.log(e.data);
            new_data = JSON.parse(e.data);

            if ('id' in new_data){

                // This is the handshake response, we've been assigned
                // and ID and are in the game.
                my.id = new_data['id'];
                console.log("Assigned ID " + my.id + " by server");
            }

            my.process_packet(new_data);

        };

        return connection;
    };

    my.process_packet = function(new_data){
        var players = "";

        // Update and Draw the Robots
        var robots = new_data['robots'];
        for (var i=0; i < robots.length; i++){
            players += "&nbsp&nbsp" + robots[i]['id'];

            if (my.debug)
                console.log({"robot": robots[i]});

            if ("position" in robots[i]){
               my.draw(robots[i], i);
            }
            if (robots[i]['id'] == my.id)
                my.update_robot(robots[i], i);
        }

        // Update and Draw the projecticles
        var projectiles = new_data['projectiles'];
        for (i=0; i < projectiles.length; i++){
            if (my.debug)
                console.log({"projectiles": projectiles[i]});

            if ("position" in projectiles[i]){
                projectiles[i]['type'] = 'bullet';
                my.draw(projectiles[i], i+1);
            }
        }

        // Set the list of players
        var players_div = document.getElementById("players");
        players_div.innerHTML = players;
    };

    my.draw = function(data, index){

        var canvas = document.getElementById('battlefield');

        if (canvas.getContext){

            // ctx.canvas.width = my.width;
            // ctx.canvas.height = my.height;

            var ctx = canvas.getContext('2d');
            var x = data['position']['x'];
            var y = data['position']['y'];

            if (index === 0)
                ctx.clearRect ( 0 , 0 , 800 , 550 );

            ctx.save();
            ctx.beginPath();
            ctx.arc(x, y, 300, 0, 2 * Math.PI, false);
            ctx.clip();

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
                ctx.fillRect (x-5, y-5, 10, 10);

                if (my.id == data['id'])
                    ctx.font="bold 22px Helvetica";
                else
                    ctx.font="10px Helvetica";

                ctx.fillText(
                    data['id'] + "[" + data['health'] + "]",
                    x+12,y+10);

                if (my.debug_draw && 'move_to' in data) {
                    ctx.restore();
                    ctx.beginPath();
                    ctx.moveTo(x, y);
                    ctx.strokeStyle = colors[0];
                    ctx.lineTo(data.move_to.x, data.move_to.y);
                    ctx.stroke();
                }
            }

            ctx.restore();
        }
    };

    my.eval_input = function( input, output ){
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
            output.innerHTML = "OK";
        }

        return theResult;
    };

    my.update_robot = function(data){
        var robot_code = editor.getSession().getValue();
        var output = document.getElementById('output');

        var code = "( " + robot_code + " )";
        var rc = my.eval_input(code, output);

        var map = {"width": my.width, "height": my.height};
        var out = rc.call(this, data, map);

        out['id'] = my.id;

        if (websocket){
            var sent_ok = websocket.send(JSON.stringify(out));
            if (my.debug){
                if (sent_ok)
                    console.log("SENT: " + JSON.stringify(out));
                else
                    console.log("ERROR SENDING TO SERVER");
            }
        }
    };


    my.init = function(){
        console.log("Welcome to GoRobots");
        websocket = my.connect(my.server);
        editor = ace.edit("editor");
        editor.setTheme("ace/theme/monokai");
        editor.getSession().setMode("ace/mode/javascript");
        document.getElementById('check_code').onclick=my.eval_input;
    };

    my.init();
    return my;
})({});}
