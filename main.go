package main

import (
	"fmt"
	"image/color"
	"math"

	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"github.com/go-p5/p5"
)

type state int
type alphabet rune

type delta map[state]map[alphabet][]state

type position struct {
	x float64
	y float64
}

func (p position) normal() position {
	mag := math.Sqrt(p.x*p.x + p.y*p.y)
	return position{x: p.x / mag, y: p.y / mag}
}

type node struct {
	highlighted bool
	pos         position
}
type record struct {
	frame uint64
	pos   position
}

// TODO: use this data structure instead of global variable
type data struct {
	nodes       []node
	firstClick  record
	secondClick record
	dragging    bool
}

const radius = 55

func (n node) contains(point position) bool {
	distance_x := math.Abs(point.x - n.pos.x)
	distance_y := math.Abs(point.y - n.pos.y)
	return distance_x*distance_x+distance_y*distance_y <= radius*radius
}

func main() {
	p5.Run(setup, draw, mouse, keyboard)
}

func setup() {
	p5.Canvas(800, 800)
}

func keyboard(e key.Event) {
	fmt.Println("keyboard event happendd", e.String())
}

var nodes []node
var firstClick record
var secondClick record
var dragging bool

var instructions []func()

func draw() {
	for i, n := range nodes {
		if n.highlighted {
			p5.Stroke(color.RGBA{R: 0xFF, A: 0xFF})
		} else {
			p5.Stroke(color.Black)
		}
		p5.Circle(n.pos.x, n.pos.y, radius)
		p5.TextSize(24)
		p5.Text(fmt.Sprintf("q%d", i), n.pos.x-10, n.pos.y+10)
	}
	if dragging {
		p5.Line(secondClick.pos.x, secondClick.pos.y, p5.Event.Mouse.Position.X, p5.Event.Mouse.Position.Y)
	}
	for i, f := range instructions {
		fmt.Println("func-", i)
		f()
	}
}

func mouse(e pointer.Event) {
	currentPositoin := position{
		x: p5.Event.Mouse.Position.X,
		y: p5.Event.Mouse.Position.Y,
	}
	switch e.Type {
	case pointer.Press:
		if e.Buttons.Contain(pointer.ButtonPrimary) {
			dragging = true
			firstClick = secondClick
			secondClick = record{
				frame: p5.FrameCount(),
				pos:   currentPositoin,
			}
			if secondClick.frame-firstClick.frame < 20 {
				nodes = append(nodes, node{
					pos: currentPositoin,
				})
				return
			}
		}
	case pointer.Release:
		dragging = false
		// handle double click on Press
		if secondClick.frame-firstClick.frame < 20 {
			break
		}
		// if it is just a click (without drag) highlight the node
		if p5.FrameCount()-secondClick.frame < 40 {
			for i := range nodes {
				if nodes[i].contains(currentPositoin) {
					nodes[i].highlighted = !nodes[i].highlighted
				}
			}
		} else { // it is a drag
			pos := secondClick.pos
			instructions = append(instructions, func() {
				p5.Line(pos.x, pos.y, currentPositoin.x, currentPositoin.y)
				p5.Translate(pos.x, pos.y)
				p5.Triangle(0, 6, 12, 0, 0, -6)
				p5.Translate(0, 0)
			})
		}
	}
}
