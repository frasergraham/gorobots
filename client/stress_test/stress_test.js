var WebSocket = require('ws');
var fs = require('fs');

var stress_test = function gorobots(my, server){
    my.width = null;
    my.height = null;
    my.server = server;
    my.websocket = null;
    my.id = null;
    my.state = null;

    var d = new Date();
    my.last_frame_time = d.getTime();
    my.delta = 0;
    my.delta_history = [];

    my.delta_max = function(){
        return Math.max.apply(null, my.delta_history);
    }

    my.health = 0;

    my.connect = function(server){
        if (!server || server.length == 0){
            setTimeout(function(){
                my.websocket = my.connect(server);
            }, 1000);
            return null;
        }

        connection = new WebSocket(server, null);

        connection.onclose = function(){
            // console.log("Lost Connection: " + server);
            my.id = null;
            my.state = null;

            // Retry every few seconds
            setTimeout(function(){
                my.websocket = my.connect(server);
            }, 3000);
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
            // console.log(my.state);
            new_data = JSON.parse(e.data);

            if (my.state == null) {
                console.log("initial response");
                if (new_data.type == "idreq") {
                    if (my.send_client_id()) {
                        my.state = "gameparam";
                    } else {
                        my.state = null;
                    }
                }
            } else if (my.state == "gameparam"){
                if (new_data.type == "gameparam") {
                    // > [OK | FULL | NOT AUTH], board size, game params
                    if (my.parse_game_params(new_data)) {
                        my.state = "handshake";
                    }
                    my.setup_robot();
                }
            } else if (my.state == "handshake") {
                my.state = "play";
            } else if(my.state == "play") {
                if (new_data.type == "handshake") {
                    // This is the handshake response, we've been assigned an ID
                    // and are in the game (or TODO: in a lobby).
                    console.log("setting up game");
                    if('success' in new_data) {
                        console.log(new_data['success']);
                        if (!new_data.success){
                            alert("invalid config!!");
                            return false;
                        }
                    }
                    if ('id' in new_data){
                        my.id = new_data['id'];
                        console.log("Assigned ID " + my.id + " by server");
                    } else {
                        console.log("server failed to send us an id")
                    }
                } else if (new_data.type == "boardstate") {
                    my.process_gameplay_packet(new_data);
                }
            }
        };

        return connection;
    };

    my.process_gameplay_packet = function(new_data){
        if (my.id){
           // fs.appendFile(my.id + ".log", JSON.stringify(new_data));
           // fs.appendFile(my.id + ".log", "\n\n\n\n\n\n\n==========\n\n\n\n\n");
        }

        var players = "";

        if (new_data.reset){
            my.setup_robot();
        }

        var robots = new_data['robots'];
        var i = 0;
        if (robots){
            for (i=0; i < robots.length; i++){
                if (robots[i]['id'] == my.id){
                    my.update_robot(robots[i], i);
                    my.health = robots[i].health;
                }
            }
        }
    };

    my.eval_input = function( input, output ){
        var theResult, evalSucceeded;

        try{
            theResult = eval( input );
        }
        catch(e){
            return null;
        }

        return theResult;
    };

    my.get_robot_code = function(){

        //======================================
        //========== ROBOT CODE ================
        //======================================

        var code = function robot(){
            var setup = function(map){
                return {
                    "hp": 100,
                    "speed": 100,
                    "weapon_radius": 35,
                    "scanner_radius": 250
                }
            }

            var update = function(data, map){
                //console.log(data, map);
                var instructions = {};
                var x,y;
                if (data && "move_to" in data){
                    x = data.move_to.x;
                    y = data.move_to.y;

                    if (Math.abs(data.position.x - data.move_to.x) < 3 &&
                        Math.abs(data.position.y - data.move_to.y) < 3){
                        x = Math.floor(Math.random() * map.width);
                        y = Math.floor(Math.random() * map.height);
                        instructions.move_to = {"x": x, "y": y};
                        instructions.fire_at = {"x": Math.floor(Math.random() * map.width), "y": Math.floor(Math.random() * map.height)};
                    }
                    else{
                        instructions.move_to = {"x": x, "y": y};
                    }
                }
                else{
                    x = Math.floor(Math.random() * map.width);
                    y = Math.floor(Math.random() * map.height);
                    instructions.move_to = {"x": x, "y": y};
                }

                if (data.scanners.length > 0){
                    instructions.fire_at = data.scanners[0].position;
                }
                else{
                    instructions.fire_at = {"x": Math.floor(Math.random() * map.width), "y": Math.floor(Math.random() * map.height)};
                }
                return instructions;
            }

            return {
                "update": update,
                "setup": setup,
            }
        }();

        //======================================
        //========== END ROBOT CODE ============
        //======================================

        return code;
    };

    my.setup_robot = function(){
        var robot = my.get_robot_code();
        if ('setup' in robot){
            var map = {"width": my.width, "height": my.height};
            var config = {
                "stats": robot.setup(map),
                "id": my.id
            };
            my.websocket.send(JSON.stringify(config));
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
            my.websocket.send(JSON.stringify(instructions));
        }
    };


    my.send_client_id = function() {
        client_id = {
            "type": "robot",
            "name": "dummy",
            "id": "24601",
            "useragent": "gorobots.js",
        };
        var sent_ok = my.websocket.send(JSON.stringify(client_id));
        if (sent_ok) {
            console.log("sent clientid: " + JSON.stringify(client_id));
        } else {
            console.log("error sending clientid to server");
        }
        return sent_ok;
    };

    my.parse_game_params = function(params) {
        // TODO: flesh out validation?
        my.width = new_data.boardsize.width;
        my.height = new_data.boardsize.height;
        console.log(my.width);
        console.log(my.height);
        return true;
    }


    my.init = function(){
        my.websocket = my.connect(my.server);
    };

    my.init();
    return my;
};

var server = process.argv[2];
var count = process.argv[3];

clients = [];
for (var i = 0; i < count; i++){
    clients.push(stress_test({}, server));
}

setInterval(function(){
    // console.log('\033[2J');
    console.log(server);
    for (client in clients){
        console.log(
            clients[client].id + "\t" +
            "[" + clients[client].health + "]" + "\t" +
            clients[client].delta_max() + " ms/f");
    }
}, 300);
