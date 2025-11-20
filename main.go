package main

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	PlayerR = 30
	finishX = 490
	finishY = 215
)

type BarrelsS struct {
	x            float64
	y            float64
	w            float64
	h            float64
	BarrelsDirY  bool
	BarrelsIndex int
	color        color.Color
}

type BouncersS struct {
	x     float64
	y     float64
	w     float64
	h     float64
	color color.Color
}

type ObstaclesS struct {
	x              float64
	y              float64
	w              float64
	h              float64
	ObstaclesIndex int
	ObstaclesDirY  bool
	color          color.Color
}
type Game struct {
	Level        int
	Barrels      []BarrelsS
	Obstacles    []ObstaclesS
	Bouncers     []BouncersS
	PlayerX      float64
	PlayerY      float64
	PlayerSpeed  float64
	PlayerBarrel int
	PlayerMoved  bool
}

func CircleRectCollision(cx, cy, cr, rx, ry, rw, rh float64) bool {
	closestX := math.Max(rx, math.Min(cx, rx+rw))
	closestY := math.Max(ry, math.Min(cy, ry+rh))
	dx := cx - closestX
	dy := cy - closestY

	return dx*dx+dy*dy <= cr*cr
}
func (g *Game) Generate_Level(L int) {
	switch L {
	case 1:
		g.Barrels = []BarrelsS{
			{50, 215, 100, 50, true, 1, color.RGBA{139, 69, 19, 255}},
			{490, 215, 100, 50, true, 2, color.RGBA{139, 69, 19, 255}},
		}
		return
	case 2:
		g.Barrels = []BarrelsS{
			{50, 215, 100, 50, true, 3, color.RGBA{139, 69, 19, 255}},
			{270, 215, 100, 50, true, 4, color.RGBA{139, 69, 19, 255}},
			{490, 400, 100, 50, true, 5, color.RGBA{139, 69, 19, 255}},
		}
		g.Bouncers = []BouncersS{
			{407, 300, 50, 50, color.RGBA{58, 110, 165, 255}},
		}
		return
	case 3:
		g.Barrels = []BarrelsS{
			{50, 215, 100, 50, true, 6, color.RGBA{139, 69, 19, 255}},
			{270, 81, 100, 50, true, 7, color.RGBA{139, 69, 19, 255}},
			{490, 400, 100, 50, true, 8, color.RGBA{139, 69, 19, 255}},
		}
		g.Obstacles = []ObstaclesS{
			{407, 350, 50, 50, 1, true, color.RGBA{160, 82, 45, 255}},
		}
		return
	case 4:
		g.Barrels = []BarrelsS{
			{50, 215, 100, 50, true, 9, color.RGBA{139, 69, 19, 255}},
			{270, 81, 100, 50, true, 10, color.RGBA{139, 69, 19, 255}},
			{490, 400, 100, 50, true, 11, color.RGBA{139, 69, 19, 255}},
		}
		g.Obstacles = []ObstaclesS{
			{205, 81, 50, 50, 2, true, color.RGBA{160, 82, 45, 255}},
		}
		g.Bouncers = []BouncersS{
			{450, 215, 50, 50, color.RGBA{58, 110, 165, 255}},
		}
	}
}
func (g *Game) Update() error {
	if g.PlayerX-PlayerR < 0 {
		g.PlayerX = 160
		g.PlayerY = 240
		g.PlayerSpeed = 15
	}
	if g.PlayerX+PlayerR > 640 {
		g.PlayerX = 160
		g.PlayerY = 240
		g.PlayerSpeed = 15
	}
	if g.PlayerY-PlayerR < 0 {
		g.PlayerX = 160
		g.PlayerY = 240
		g.PlayerSpeed = 15
	}
	if g.PlayerY+PlayerR > 480 {
		g.PlayerX = 160
		g.PlayerY = 240
		g.PlayerSpeed = 15
	}
	for i, b := range g.Barrels {
		if b.BarrelsIndex == 4 || b.BarrelsIndex == 7 || b.BarrelsIndex == 8 || b.BarrelsIndex == 10 {
			if g.Barrels[i].y > 400 {
				g.Barrels[i].BarrelsDirY = false
			}
			if g.Barrels[i].y < 80 {
				g.Barrels[i].BarrelsDirY = true
			}
			if !g.Barrels[i].BarrelsDirY {
				g.Barrels[i].y -= 2
			}
			if g.Barrels[i].BarrelsDirY {
				g.Barrels[i].y += 2
			}
			if CircleRectCollision(g.PlayerX, g.PlayerY, PlayerR, b.x, b.y, b.w, b.h) {
				g.PlayerY = g.Barrels[i].y + 25
			}
		}
	}
	for i, o := range g.Obstacles {
		if o.ObstaclesIndex == 2 {
			if g.Obstacles[i].y > 400 {
				g.Obstacles[i].ObstaclesDirY = false
			}
			if g.Obstacles[i].y < 80 {
				g.Obstacles[i].ObstaclesDirY = true
			}
			if !g.Obstacles[i].ObstaclesDirY {
				g.Obstacles[i].y -= 2
			}
			if g.Obstacles[i].ObstaclesDirY {
				g.Obstacles[i].y += 2
			}
			if CircleRectCollision(g.PlayerX, g.PlayerY, PlayerR, o.x, o.y, o.w, o.h) {
				g.PlayerX = 160
				g.PlayerY = 240
			}
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) && !g.PlayerMoved {
		g.PlayerMoved = true
	}
	if g.PlayerMoved {
		g.PlayerX += g.PlayerSpeed
	}
	for _, b := range g.Barrels {
		x := b.x + 125
		w := b.w - 125

		if g.PlayerSpeed == -15 {
			x = b.x - 10
			w = b.w
		}

		if CircleRectCollision(g.PlayerX, g.PlayerY, PlayerR, x, b.y, w, b.h) {
			g.PlayerMoved = false
			g.PlayerSpeed = 15
			if g.Level == 1 {
				if CircleRectCollision(g.PlayerX, g.PlayerY, PlayerR, finishX, finishY, b.w, b.h) {
					g.Level++
					g.PlayerX = 160
					g.PlayerY = 240
					g.PlayerSpeed = 15
					g.Generate_Level(g.Level)
				}
			} else if g.Level == 2 || g.Level == 3 {
				if CircleRectCollision(g.PlayerX, g.PlayerY, PlayerR, 490, 400, 100, 50) {
					g.Level++
					g.PlayerX = 160
					g.PlayerY = 240
					g.PlayerSpeed = 15
					g.Generate_Level(g.Level)
				}
			}
		}
	}
	for _, o := range g.Obstacles {
		if CircleRectCollision(g.PlayerX, g.PlayerY, PlayerR, o.x, o.y, o.w, o.h) {
			g.PlayerX = 160
			g.PlayerY = 240
			g.PlayerSpeed = 15
		}
	}
	for _, boun := range g.Bouncers {
		if CircleRectCollision(g.PlayerX, g.PlayerY, PlayerR, boun.x, boun.y, boun.w, boun.h) {
			g.PlayerSpeed = -15
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	//draw barrels
	for _, b := range g.Barrels {
		ebitenutil.DrawRect(screen, b.x, b.y, b.w, b.h, b.color)
	}
	//draw obstacles
	for _, o := range g.Obstacles {
		ebitenutil.DrawRect(screen, o.x, o.y, o.w, o.h, o.color)
	}
	//draw Bouncers
	for _, boun := range g.Bouncers {
		ebitenutil.DrawRect(screen, boun.x, boun.y, boun.w, boun.h, boun.color)
	}
	//draw player
	ebitenutil.DrawCircle(screen, g.PlayerX, g.PlayerY, PlayerR, color.RGBA{255, 255, 0, 255})
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello World - Ebiten")
	g := &Game{
		PlayerY:     240,
		PlayerX:     160,
		PlayerSpeed: 15,
		Level:       1,
		Barrels: []BarrelsS{
			{
				x:            50,
				y:            215,
				w:            100,
				h:            50,
				BarrelsDirY:  true,
				BarrelsIndex: 1,
				color:        color.RGBA{139, 69, 19, 255},
			},
			{
				x:            490,
				y:            215,
				w:            100,
				h:            50,
				BarrelsDirY:  true,
				BarrelsIndex: 2,
				color:        color.RGBA{139, 69, 19, 255},
			},
		},
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
