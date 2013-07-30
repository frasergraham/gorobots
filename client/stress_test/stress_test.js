var WebSocket = require('ws');

var stress_test = function gorobots(my, server){
    my.width = 800;
    my.height = 550;
    my.server = server;
    my.websocket = null;
    my.id = null;

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
            new_data = JSON.parse(e.data);

            if ('id' in new_data){

                // This is the handshake response, we've been assigned
                // an ID and are in the game.
                my.id = new_data['id'];
                my.setup_robot();
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
    console.log('\033[2J');
    console.log(server);
    for (client in clients){
        console.log(
            clients[client].id + "\t" +
            "[" + clients[client].health + "]" + "\t" +
            clients[client].delta_max() + " ms/f");
    }
}, 300);
