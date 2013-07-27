function robot(data){
        //console.log(data);
        var instructions = {};
        var x,y;
        if (data && 'move_to' in data){
            x = data['move_to']['x'];
            y = data['move_to']['y'];

            if (data['position']['x'] == data['move_to']['x'] &&
                data['position']['y'] == data['move_to']['y']){
                x = Math.floor(Math.random() * 450);
                y = Math.floor(Math.random() * 400);
                console.log("X: " + x + ", Y: " + y);
                instructions['move_to'] = {"x": x, "y": y};
            }
            else{
                instructions['move_to'] = {"x": x, "y": y};
            }
        }
        else{
            x = Math.floor(Math.random() * 450);
            y = Math.floor(Math.random() * 400);
            instructions['move_to'] = {"x": x, "y": y};
        }

        return instructions;
    }

