function robot(data, map){

        //console.log(data, map);
        var instructions = {};
        var x,y;
        if (data && 'move_to' in data){
            x = data['move_to']['x'];
            y = data['move_to']['y'];

            if (Math.abs(data['position']['x'] - data['move_to']['x']) < 3 &&
                Math.abs(data['position']['y'] - data['move_to']['y']) < 3){
                x = Math.floor(Math.random() * map['width']);
                y = Math.floor(Math.random() * map['height']);
                console.log("X: " + x + ", Y: " + y);
                instructions['move_to'] = {"x": x, "y": y};
            }
            else{
                instructions['move_to'] = {"x": x, "y": y};
            }
        }
        else{
            x = Math.floor(Math.random() * map['width']);
            y = Math.floor(Math.random() * map['height']);
            instructions['move_to'] = {"x": x, "y": y};
        }

        return instructions;
    }

    