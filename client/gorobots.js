function init(){(function gorobots(my){
    my.width = 800;
    my.height = 550;
    my.server = "ws://twisted.local:8666/ws/";
    // my.server = "ws://rs.mcquay.me:8666/ws/";
    my.websocket = null;
    my.id = null;
    my.ctx = null;
    my.observe_only = false;
    my.y_base = 16; // The top status bar
    my.connection_retry = 3000;

    var d = new Date();
    my.last_frame_time = d.getTime();
    my.delta = 0;
    my.delta_history = [];

    my.debug = false;
    my.debug_draw = true;

    var grey = "rgba(0, 0, 0, 0.1)";
    var white = "rgba(255, 255, 255, 1)";
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


    my.toggle_debug = function(){
        if (my.debug_draw)
            my.debug_draw = false;
        else
            my.debug_draw = true;

        console.log("toggling debug status to " + my.observe_only);
    };

    my.toggle_observer = function(){
        if (my.observe_only)
            my.observe_only = false;
        else
            my.observe_only = true;

        console.log("toggling observer status to " + my.observe_only);
    };

    my.connect = function(server){

        connection = new WebSocket(server, null);

        connection.onerror = function (error) {
          console.log('WebSocket Error ' + error);
        };

        connection.onopen = function(){
            console.log("Connected to " + server);
        };

        connection.onclose = function(){
            my.id = null;

            if (my.connection_retry > 0){
                // Retry every few seconds
                console.log("Lost Connection: " + server);
                console.log(my);
                setTimeout(function(){
                    my.websocket = my.connect(my.server);
                }, my.connection_retry);
            }
        };

        connection.onmessage = function (e) {

            var d = new Date();
            var time = d.getTime();
            my.delta = time - my.last_frame_time;
            my.delta_history.unshift(my.delta);
            if (my.delta_history.length > 50){
                my.delta_history.pop();
            }

            my.last_frame_time = time;

            // console.log(e.data);
            new_data = JSON.parse(e.data);

            if ('id' in new_data){

                // This is the handshake response, we've been assigned
                // an ID and are in the game.
                my.id = new_data['id'];
                console.log("Assigned ID " + my.id + " by server");
                my.setup_robot();
            }
            if('success' in new_data) {
                console.log(new_data['success']);
                if (!new_data.success){
                    alert("invalid config!!");
                    my.websocket = null
                    return false;
                }
            }

            my.process_packet(new_data);
        };

        return connection;
    };

    my.process_packet = function(new_data){
        var players = "";

        if (new_data.reset){
            my.setup_robot();
        }

        my.clear();

        // Status Bar Text
        var slowest = Math.max.apply(null, my.delta_history);
        var debug_string = slowest + "ms" + "    " + my.server;
        my.ctx.fillStyle = colors[0];
        my.ctx.font="12px Helvetica";
        my.ctx.fillText(debug_string,2,12);

        // Status Bar Line
        my.ctx.beginPath();
        my.ctx.moveTo(0, my.y_base);
        my.ctx.strokeStyle = colors[0];
        my.ctx.lineTo(my.width, my.y_base);
        my.ctx.stroke();

        my.ctx.save();
        my.ctx.beginPath();
        my.ctx.rect(0,my.y_base,my.width,my.height);
        my.ctx.clip();

        // Update and Draw the Robots
        var robots = new_data['robots'];
        var i = 0;

        if (robots){
            for (i=0; i < robots.length; i++){
                if ("position" in robots[i]){
                    if (robots[i]['id'] == my.id)
                        my.clip(robots[i]);
                }
            }

            for (i=0; i < robots.length; i++){
                var col = colors[i+1];
                if (robots[i].health === 0){
                    col = colors[0];
                }
                players += ("<span style='color: " + col + ";'>&nbsp&nbsp" + robots[i]['id'] +
                    " [" + robots[i]['health'] + "]</span>");

                if (my.debug)
                    console.log(JSON.stringify(robots[i]));

                if ("position" in robots[i]){
                    my.draw(robots[i], i);
                }

                // Update my the robot for this client
                if (robots[i]['id'] == my.id)
                    my.update_robot(robots[i], i);
            }
        }

        // Draw the projecticles
        var projectiles = new_data['projectiles'];
        if (projectiles){
            for (i=0; i < projectiles.length; i++){
                if (my.debug)
                    console.log(JSON.stringify(projectiles[i]));

                if ("position" in projectiles[i]){
                    projectiles[i]['type'] = 'bullet';
                    my.draw(projectiles[i], i+1);
                }
            }
        }

        // Draw the projecticles
        var splosions = new_data['splosions'];
        if (splosions){
            for (i=0; i < splosions.length; i++){
                if (my.debug)
                    console.log(JSON.stringify(splosions[i]));

                if ("position" in splosions[i]){
                    splosions[i]['type'] = 'explosion';
                    my.draw(splosions[i], i+1);
                }
            }
        }

        // Set the list of players
        var players_div = document.getElementById("players");
        players_div.innerHTML = players;
    };

    my.clip = function(robot){
        if (my.observe_only){
            var x_scale = my.ctx.canvas.width / my.width;
            var y_scale = (my.ctx.canvas.height - my.y_base )/ my.height;

            my.ctx.fillStyle = grey;
            my.ctx.fillRect (0, 0+my.y_base, my.width, my.height);

            my.ctx.beginPath();
            my.ctx.fillStyle = white;
            my.ctx.arc(robot.position.x * x_scale, robot.position.y * y_scale + my.y_base, robot.stats.scanner_radius, 0, 2 * Math.PI, false);
            my.ctx.fill();

            my.ctx.save();
            my.ctx.beginPath();
            my.ctx.arc(robot.position.x * x_scale, robot.position.y * y_scale + my.y_base, robot.stats.scanner_radius, 0, 2 * Math.PI, false);
            my.ctx.clip();
        }
    };

    my.clear = function(){
        my.ctx.restore();
        my.ctx.clearRect ( 0 , 0 , my.width , my.height );
    };

    my.draw = function(data, index){

        var x_scale = my.ctx.canvas.width / my.width;
        var y_scale = (my.ctx.canvas.height - my.y_base )/ my.height;

        var x = data['position']['x'] * x_scale;
        var y = data['position']['y'] * y_scale + my.y_base;

        if ('type' in data && data['type'] == 'bullet'){
            my.ctx.fillStyle = colors[0];
            my.ctx.fillRect (x, y, 5, 5);
        }
        else if ('type' in data && data['type'] == 'explosion'){
            my.ctx.beginPath();
            my.ctx.arc(x, y, data.radius * x_scale, 0, 2 * Math.PI, false);
            my.ctx.fillStyle = colors[3];
            my.ctx.fill();
        }
        else{
            if (data.health > 0)
                my.ctx.fillStyle = colors[index+1];
            else
                my.ctx.fillStyle = colors[0];

            my.ctx.fillRect (x-5, y-5, 10, 10);

            if (my.id == data['id'])
                my.ctx.font="bold 22px Helvetica";
            else
                my.ctx.font="10px Helvetica";

            my.ctx.fillText(
                data['id'] + "[" + data['health'] + "]",
                x+12,y+10);

            if (my.debug_draw && 'move_to' in data && my.id == data['id']) {
                // my.ctx.restore();
                my.ctx.beginPath();
                my.ctx.moveTo(x, y);
                my.ctx.strokeStyle = colors[0];
                my.ctx.lineTo(data.move_to.x * x_scale, data.move_to.y * y_scale + my.y_base);
                my.ctx.stroke();
            }
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

    my.get_robot_code = function(){
        var robot_code = editor.getSession().getValue();
        var output = document.getElementById('output');

        var code = "( " + robot_code + " )";
        var rc = my.eval_input(code, output);

        return rc.call(this);
    };

    my.setup_robot = function(){
        var robot = my.get_robot_code();
        console.log(robot);
        if ('setup' in robot){
            var map = {"width": my.width, "height": my.height};
            var config = {
                "stats": robot.setup(map),
                "id": my.id
            };
            console.log(config);
            var sent_ok = my.websocket.send(JSON.stringify(config));
            if (true){
                if (sent_ok)
                    console.log("SENT CONFIG: " + JSON.stringify(config));
                else
                    console.log("ERROR SENDING CONFIG TO SERVER");
            }
        }
    };

    my.update_robot = function(data){
        var robot = my.get_robot_code();
        var map = {"width": my.width, "height": my.height};

        var instructions = null;
        if ('update' in robot){
            instructions = robot.update(data, map);
        }
        instructions['id'] = my.id;

        if (my.websocket){
            var sent_ok = my.websocket.send(JSON.stringify(instructions));
            if (my.debug){
                if (sent_ok)
                    console.log("SENT: " + JSON.stringify(instructions));
                else
                    console.log("ERROR SENDING TO SERVER");
            }
        }
    };


    my.init = function(){
        console.log("Welcome to GoRobots");
        // my.websocket = my.connect(my.server);
        editor = ace.edit("editor");
        editor.setTheme("ace/theme/monokai");
        editor.getSession().setMode("ace/mode/javascript");

        document.getElementById('debug_toggle').onclick=function(){
                my.toggle_debug();
            };

        document.getElementById('fov_toggle').onclick=function(){
                my.toggle_observer();
            };

        var options = decodeURIComponent(window.location.search.slice(1));
        options = options.split('=');

        if (options[0] == 'server'){
            my.server = options[1];
            my.websocket = my.connect(my.server);
        }

        var server_name = document.getElementById("server");
        server_name.value = my.server;
        var form = document.getElementById("form");
        form.onsubmit = function(e){
            e.preventDefault();
            my.server = server_name.value;
            if (my.websocket){
                console.log("Switching Server: " + my.server);
                my.websocket.close();
            }
            else{
                console.log("Setting Server: " + my.server);
                my.websocket = my.connect(my.server);
            }
            return false;
        }

        var canvas = document.getElementById('battlefield');
        if (canvas.getContext){
            my.ctx = canvas.getContext('2d');
            my.ctx.canvas.width = my.width;
            my.ctx.canvas.height = my.height;
        }
        else{
            console.log("Canvas Error");
        }
    };

    my.init();
    return my;
})({});}
