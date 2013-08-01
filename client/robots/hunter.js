function robot(){
    var setup = function(map){
        return {
            "hp": 100,
            "speed": 200,
            "weapon_radius": 30,
            "scanner_radius": 500
        }
    }

    var distance = function(pos1,pos2){
        var delta_x = pos1.x - pos2.x;
        var delta_y = pos1.y - pos2.y;
        return Math.sqrt(delta_x * delta_x + delta_y * delta_y);
    };

    var update = function(data, map){
        //console.log(data, map);
        var instructions = {};
        var x,y;
        if (data && 'move_to' in data){
            x = data.move_to.x;
            y = data.move_to.y;

            if (Math.abs(data.position.x - data.move_to.x) < 3 &&
                Math.abs(data.position.y - data.move_to.y) < 3){
                x = Math.floor(0.01 * map.width);
                y = Math.floor(Math.random() * map.height);
                instructions.move_to = {"x": x, "y": y};
                //instructions.fire_at = {"x": Math.floor(Math.random() * map.width), "y": Math.floor(Math.random() * map.height)};
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
            var d = distance(data.position, data.scanners[0].position);

            if (typeof(last_x) !== 'undefined'){
                dx = data.scanners[0].position.x - last_x;
                dy =  data.scanners[0].position.y - last_y;
            }

            last_x = data.scanners[0].position.x;
            last_y = data.scanners[0].position.y

            // console.log(data.scanners[0], dx, dy);
            if (d > 50){
                instructions.fire_at = {
                    "x": data.scanners[0].position.x + 15*dx,
                    "y": data.scanners[0].position.y + 15*dy,
                }
            }

            if (d > 100)
                instructions.move_to = data.scanners[0].position;
            else
                instructions.move_to = {
                    "x": data.scanners[0].position.y + d,
                    "y": data.scanners[0].position.y + d,
                }
        }
        else{
            instructions.fire_at = {"x": 0, "y": 0};
        }
        return instructions;
    }

    return {
        "update": update,
        "setup": setup,
    }
}

