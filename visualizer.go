//go:build visualizer
// +build visualizer

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type visData struct {
	Rooms  []Room
	Links  [][2]string
	Start  string
	End    string
	Moves  [][]move
	Width  int
	Height int
}

type move struct {
	Ant  int    `json:"ant"`
	Room string `json:"room"`
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// split input lines into description and moves
	idx := 0
	for idx < len(lines) && strings.TrimSpace(lines[idx]) != "" {
		idx++
	}
	descLines := lines[:idx]
	var moveLines []string
	if idx < len(lines) {
		moveLines = lines[idx+1:]
	}

	g, _, _, err := Parse(strings.NewReader(strings.Join(descLines, "\n")))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var moves [][]move
	for _, l := range moveLines {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}
		parts := strings.Split(l, " ")
		var step []move
		for _, p := range parts {
			if !strings.HasPrefix(p, "L") {
				continue
			}
			sp := strings.SplitN(p[1:], "-", 2)
			if len(sp) != 2 {
				continue
			}
			a, err := strconv.Atoi(sp[0])
			if err != nil {
				continue
			}
			step = append(step, move{Ant: a, Room: sp[1]})
		}
		if len(step) > 0 {
			moves = append(moves, step)
		}
	}

	var rooms []Room
	maxX, maxY := 0, 0
	for _, r := range g.Rooms {
		rooms = append(rooms, *r)
		if r.X > maxX {
			maxX = r.X
		}
		if r.Y > maxY {
			maxY = r.Y
		}
	}
	var links [][2]string
	seen := make(map[[2]string]bool)
	for a, bs := range g.Links {
		for _, b := range bs {
			pair := [2]string{a, b}
			if a > b {
				pair = [2]string{b, a}
			}
			if !seen[pair] {
				seen[pair] = true
				links = append(links, pair)
			}
		}
	}

	data := visData{
		Rooms:  rooms,
		Links:  links,
		Start:  g.Start,
		End:    g.End,
		Moves:  moves,
		Width:  maxX*40 + 80,
		Height: maxY*40 + 80,
	}

	tmpl := template.Must(template.New("page").Parse(pageHTML))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		js, _ := json.Marshal(data)
		tmpl.Execute(w, template.JS(js))
	})

	fmt.Println("Visualizer running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

const pageHTML = `<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<title>lem-in visualizer</title>
<style>
  body { font-family: sans-serif; }
  canvas { border: 1px solid #333; }
</style>
</head>
<body>
<canvas id="c" width="{{.Width}}" height="{{.Height}}"></canvas>
<script>
const data = {{.}};
const scale = 40;
const radius = 10;
const canvas = document.getElementById('c');
const ctx = canvas.getContext('2d');

function pos(roomName){
    for (const r of data.Rooms){
        if (r.Name === roomName) return {x:r.X*scale+40, y:r.Y*scale+40};
    }
    return {x:0,y:0};
}

function draw(step){
    ctx.clearRect(0,0,canvas.width,canvas.height);
    ctx.lineWidth = 2;
    // links
    ctx.strokeStyle = '#555';
    for (const l of data.Links){
        const a = pos(l[0]);
        const b = pos(l[1]);
        ctx.beginPath();
        ctx.moveTo(a.x,a.y);
        ctx.lineTo(b.x,b.y);
        ctx.stroke();
    }
    // rooms
    for (const r of data.Rooms){
        const p = pos(r.Name);
        ctx.fillStyle = (r.Name===data.Start)?'green':(r.Name===data.End?'red':'white');
        ctx.strokeStyle = '#000';
        ctx.beginPath();
        ctx.arc(p.x,p.y,radius,0,Math.PI*2);
        ctx.fill();
        ctx.stroke();
        ctx.fillStyle = '#000';
        ctx.fillText(r.Name,p.x+radius,p.y);
    }
    // ants
    for (const ant in positions){
        const room=positions[ant];
        const p=pos(room);
        ctx.fillStyle='blue';
        ctx.beginPath();
        ctx.arc(p.x,p.y,radius/2,0,Math.PI*2);
        ctx.fill();
        ctx.fillStyle='white';
        ctx.fillText(ant,p.x-4,p.y+4);
    }
}

let positions = {};
let step = 0;
function advance(){
    if (step < data.Moves.length){
        for (const m of data.Moves[step]){
            positions[m.ant] = m.room;
            if(m.room===data.End){
                delete positions[m.ant];
            }
        }
        draw(step);
        step++;
        setTimeout(advance,800);
    } else {
        draw(step);
    }
}

draw(0);
advance();
</script>
</body>
</html>`
