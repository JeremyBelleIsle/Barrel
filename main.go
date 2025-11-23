package main

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	PlayerR = 30
)

type Particle struct {
	x, y   float64
	vx, vy float64
	life   int
	color  color.Color
}
type BarrelsS struct {
	x           float64
	y           float64
	w           float64
	h           float64
	BarrelsDirY bool
	Moved       bool
	fragile     bool
	Teleporter  bool
	TeleporterX float64
	TeleporterY float64
	CoolDown    int
	color       color.Color
}

type BouncersS struct {
	x     float64
	y     float64
	w     float64
	h     float64
	color color.Color
}

type ObstaclesS struct {
	x             float64
	y             float64
	w             float64
	h             float64
	Moved         bool
	ObstaclesDirY bool
	color         color.Color
}
type Game struct {
	Level               int
	Particles           []Particle
	Barrels             []BarrelsS
	Obstacles           []ObstaclesS
	Bouncers            []BouncersS
	TimeBeforeLevelDown int
	PlayerX             float64
	PlayerY             float64
	PlayerSpeed         float64
	PlayerBarrel        int
	PlayerLife          int
	PlayerMoved         bool
}

var (
	mplusFaceSource *text.GoTextFaceSource
)

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.PressStart2P_ttf))
	if err != nil {
		log.Fatal(err)
	}
	mplusFaceSource = s
}
func CircleRectCollision(cx, cy, cr, rx, ry, rw, rh float64) bool {
	closestX := math.Max(rx, math.Min(cx, rx+rw))
	closestY := math.Max(ry, math.Min(cy, ry+rh))
	dx := cx - closestX
	dy := cy - closestY

	return dx*dx+dy*dy <= cr*cr
}
func (g *Game) SpawnBarrelExplosion(x, y float64) {
	for i := 0; i < 20; i++ {
		g.Particles = append(g.Particles, Particle{
			x:     x + 50, // milieu du baril
			y:     y + 25,
			vx:    (rand.Float64()*4 - 2), // -2 à +2
			vy:    (rand.Float64()*4 - 2),
			life:  25 + rand.Intn(10),
			color: color.RGBA{169, 99, 29, 255}, // brun clair éclats bois
		})
	}
}
func (g *Game) Generate_Level(L int) {
	switch L {
	case 1:
		g.Barrels = []BarrelsS{
			{50, 215, 100, 50, true, false, false, false, 0, 0, 100, color.RGBA{139, 69, 19, 255}},
			{490, 215, 100, 50, true, false, false, false, 0, 0, 100, color.RGBA{139, 69, 19, 255}},
		}
		g.Bouncers = []BouncersS{}
		g.Obstacles = []ObstaclesS{}
		return
	case 2:
		g.Barrels = []BarrelsS{
			{50, 215, 100, 50, true, false, true, false, 0, 0, 100, color.RGBA{139, 69, 19, 255}},
			{270, 215, 100, 50, true, true, false, false, 0, 0, 100, color.RGBA{139, 69, 19, 255}},
			{490, 400, 100, 50, true, false, false, false, 0, 0, 100, color.RGBA{139, 69, 19, 255}},
		}
		g.Bouncers = []BouncersS{
			{407, 300, 50, 50, color.RGBA{58, 110, 165, 255}},
		}
		g.Obstacles = []ObstaclesS{}
		return
	case 3:
		g.Barrels = []BarrelsS{
			{50, 215, 100, 50, true, false, false, false, 0, 0, 100, color.RGBA{139, 69, 19, 255}},
			{270, 81, 100, 50, true, true, false, false, 0, 0, 100, color.RGBA{139, 69, 19, 255}},
			{490, 400, 100, 50, true, true, false, false, 0, 0, 100, color.RGBA{139, 69, 19, 255}},
		}
		g.Obstacles = []ObstaclesS{
			{407, 350, 50, 50, false, true, color.RGBA{178, 34, 34, 255}},
		}
		g.Bouncers = []BouncersS{}
		return
	case 4:
		g.Barrels = []BarrelsS{
			{50, 215, 100, 50, true, false, false, false, 0, 0, 100, color.RGBA{139, 69, 19, 255}},
			{270, 81, 100, 50, true, true, false, false, 0, 0, 100, color.RGBA{139, 69, 19, 255}},
			{490, 400, 100, 50, true, false, false, false, 0, 0, 100, color.RGBA{139, 69, 19, 255}},
		}
		g.Obstacles = []ObstaclesS{
			{205, 81, 50, 50, true, true, color.RGBA{178, 34, 34, 255}},
		}
		g.Bouncers = []BouncersS{
			{450, 215, 50, 50, color.RGBA{58, 110, 165, 255}},
		}
	case 5:
		g.Barrels = []BarrelsS{
			{50, 215, 100, 50, true, true, false, false, 0, 0, 100, color.RGBA{139, 69, 19, 255}},
			{270, 0, 100, 50, true, true, true, false, 0, 0, 100, color.RGBA{139, 69, 19, 255}},
			{375, 250, 100, 50, false, true, false, false, 0, 0, 100, color.RGBA{139, 69, 19, 255}},
			{490, 50, 100, 50, true, false, false, false, 0, 0, 100, color.RGBA{139, 69, 19, 255}},
		}
		g.Obstacles = []ObstaclesS{
			{205, 300, 50, 50, true, false, color.RGBA{178, 34, 34, 255}},
			{310, 50, 50, 50, false, false, color.RGBA{178, 34, 34, 255}},
		}
		g.Bouncers = []BouncersS{}
	case 6:
		g.Barrels = []BarrelsS{
			{50, 215, 100, 50, true, false, false, false, 330, 60, 100, color.RGBA{139, 69, 19, 255}},
			{400, 215, 100, 50, true, false, false, true, 330, 75, 100, color.RGBA{139, 69, 19, 255}},
			{250, 50, 100, 50, true, false, false, false, 0, 0, 100, color.RGBA{139, 69, 19, 255}},
			{490, 50, 100, 50, true, false, false, false, 0, 0, 100, color.RGBA{139, 69, 19, 255}},
		}
		g.Obstacles = []ObstaclesS{}
		g.Bouncers = []BouncersS{}
	}
}
func (g *Game) DrawLifes(screen *ebiten.Image) error {
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(200), float64(50))
	op.ColorScale.ScaleWithColor(color.RGBA{222, 49, 99, 0})
	text.Draw(screen, fmt.Sprintf("Vies :%d", g.PlayerLife), &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   34,
	}, op)
	return nil
}
func (g *Game) DrawLevel(screen *ebiten.Image) error {
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(200), float64(90))
	op.ColorScale.ScaleWithColor(color.RGBA{222, 49, 99, 0})
	text.Draw(screen, fmt.Sprintf("Level :%d", g.Level), &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   34,
	}, op)
	return nil
}
func (g *Game) Update() error {
	if g.PlayerLife == 0 {
		g.Level = 1
		g.PlayerX = 160
		g.PlayerY = 240
		g.PlayerLife = 3
		g.PlayerSpeed = 15
		defer g.Generate_Level(g.Level)
	}
	if g.TimeBeforeLevelDown > 0 {
		g.TimeBeforeLevelDown--
	}
	if g.TimeBeforeLevelDown == 0 {
		g.Level--
		g.PlayerX = 160
		g.PlayerY = 240
		g.PlayerSpeed = 15
		defer g.Generate_Level(g.Level)
		g.TimeBeforeLevelDown = -67
	}
	// --- BORDS DE L'ÉCRAN ---
	if g.PlayerX-PlayerR < 0 ||
		g.PlayerX+PlayerR > 640 ||
		g.PlayerY-PlayerR < 0 ||
		g.PlayerY+PlayerR > 480 {

		g.PlayerX = 160
		g.PlayerY = 240
		g.PlayerSpeed = 15
		g.PlayerLife--
	}

	// --- SUPPRESSION DES BARILS ---
	deleteBarrels := []int{} // tableau vide OK

	for i := range g.Barrels {
		b := &g.Barrels[i] // pointeur = on modifie réellement le slice

		// --- Barils qui bougent verticalement ---
		if b.Moved {

			if b.y > 400 {
				b.BarrelsDirY = false
			}
			if b.y < 80 {
				b.BarrelsDirY = true
			}

			if b.BarrelsDirY {
				b.y += 2
			} else {
				b.y -= 2
			}

			// Player sur le baril = il reste dessus
			if CircleRectCollision(g.PlayerX, g.PlayerY, PlayerR, b.x, b.y, b.w, b.h) {
				g.PlayerY = b.y + 25
			}
		}

		// ---- BARILS FRAGILES ----
		if b.fragile && CircleRectCollision(g.PlayerX, g.PlayerY, PlayerR, b.x, b.y, b.w, b.h) {
			b.CoolDown--
			b.color = color.RGBA{255, 0, 0, 255}
			if b.CoolDown <= 0 {
				g.TimeBeforeLevelDown = 45
				g.SpawnBarrelExplosion(b.x, b.y)
				deleteBarrels = append(deleteBarrels, i)
				b.CoolDown = 100
			}
		}
		if b.Teleporter && CircleRectCollision(g.PlayerX, g.PlayerY, PlayerR, b.x, b.y, b.w, b.h) {
			g.PlayerX = b.TeleporterX
			g.PlayerY = b.TeleporterY
		}
		if !CircleRectCollision(g.PlayerX, g.PlayerY, PlayerR, b.x, b.y, b.w, b.h) {
			b.color = color.RGBA{139, 69, 19, 255}
		}
	}

	// --- ON SUPPRIME EN PARTANT DE LA FIN ---
	for i := len(deleteBarrels) - 1; i >= 0; i-- {
		g.Barrels = slices.Delete(g.Barrels, deleteBarrels[i], deleteBarrels[i]+1)
	}

	// --- OBSTACLES ---
	for i := range g.Obstacles {
		o := &g.Obstacles[i]

		if o.Moved {

			if o.y > 400 {
				o.ObstaclesDirY = false
			}
			if o.y < 80 {
				o.ObstaclesDirY = true
			}

			if o.ObstaclesDirY {
				o.y += 2
			} else {
				o.y -= 2
			}

			if CircleRectCollision(g.PlayerX, g.PlayerY, PlayerR, o.x, o.y, o.w, o.h) {
				g.PlayerX = 160
				g.PlayerY = 240
				g.PlayerLife--
			}
		}
	}

	// --- MOUVEMENT DU PLAYER ---
	if ebiten.IsKeyPressed(ebiten.KeySpace) && !g.PlayerMoved {
		g.PlayerMoved = true
	}
	if g.PlayerMoved {
		g.PlayerX += g.PlayerSpeed
	}

	// --- COLLISIONS AVEC BARILS ---
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

			// --- CHANGER DE NIVEAU ---
			if b.x == 490 {
				g.Level++
				g.PlayerX = 160
				g.PlayerY = 240
				g.PlayerSpeed = 15
				g.Generate_Level(g.Level)
			}
		}
	}

	// --- OBSTACLES NORMAUX ---
	for _, o := range g.Obstacles {
		if CircleRectCollision(g.PlayerX, g.PlayerY, PlayerR, o.x, o.y, o.w, o.h) {
			g.PlayerX = 160
			g.PlayerY = 240
			g.PlayerSpeed = 15
		}
	}

	// --- BOUNCERS ---
	for _, boun := range g.Bouncers {
		if CircleRectCollision(g.PlayerX, g.PlayerY, PlayerR, boun.x, boun.y, boun.w, boun.h) {
			g.PlayerSpeed = -15
		}
	}
	for i := 0; i < len(g.Particles); i++ {
		p := &g.Particles[i]

		p.x += p.vx
		p.y += p.vy
		p.vy += 0.1 // gravité légère

		p.life--
		if p.life <= 0 {
			g.Particles = append(g.Particles[:i], g.Particles[i+1:]...)
			i--
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	//draw Lifes
	g.DrawLifes(screen)
	//draw Level
	g.DrawLevel(screen)
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
	//draw particules
	for _, p := range g.Particles {
		ebitenutil.DrawRect(screen, p.x, p.y, 4, 4, p.color)
	}
	//draw teleporter lines
	for _, b := range g.Barrels {
		if b.Teleporter {
			ebitenutil.DrawLine(screen, b.x, b.y, b.TeleporterX, b.TeleporterY, color.White)
		}
	}
	//draw player
	ebitenutil.DrawCircle(screen, g.PlayerX, g.PlayerY, PlayerR, color.RGBA{255, 255, 0, 255})
	//draw version
	ebitenutil.DebugPrintAt(screen, "version 0.6", 640-120, 480-20)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello World - Ebiten")
	g := &Game{
		PlayerY:             240,
		PlayerX:             160,
		PlayerSpeed:         15,
		Level:               1,
		PlayerLife:          3,
		TimeBeforeLevelDown: -67,
		Barrels: []BarrelsS{
			{
				x:           50,
				y:           215,
				w:           100,
				h:           50,
				BarrelsDirY: true,
				Moved:       false,
				fragile:     false,
				Teleporter:  false,
				TeleporterX: 0,
				TeleporterY: 0,
				CoolDown:    100,
				color:       color.RGBA{139, 69, 19, 255},
			},
			{
				x:           490,
				y:           215,
				w:           100,
				h:           50,
				BarrelsDirY: true,
				Moved:       false,
				fragile:     false,
				Teleporter:  false,
				TeleporterX: 0,
				TeleporterY: 0,
				CoolDown:    100,
				color:       color.RGBA{139, 69, 19, 255},
			},
		},
	}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
