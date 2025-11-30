package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"os"
	"slices"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	PlayerR = 30
)

type Score struct {
	Time     time.Duration
	UserName string
}

type SaveData struct {
	Top5 []Score
}
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
	Level                 int
	Particles             []Particle
	Barrels               []BarrelsS
	Obstacles             []ObstaclesS
	Bouncers              []BouncersS
	Top5Bestplayers       []Score
	Save                  SaveData
	StartTime             time.Time
	State                 int
	endTime               time.Duration
	EndTimeIsSet          bool
	DataAreSave           bool
	Opacity               float64
	Radius1               float64
	Radius2               float64
	Radius3               float64
	MOL1                  bool
	MOL2                  bool
	MOL3                  bool
	OpacityPlusOrNegative bool
	ChangeLevelAnimation  bool
	EndOfRun              bool
	TimeSaveAnimation     float64
	TimeBeforeLevelDown   int
	SpaceCNT              int
	PlayerX               float64
	PlayerY               float64
	PlayerSpeed           float64
	currentUserName       string
	PlayerBarrel          int
	PlayerLife            int
	PlayerMoved           bool
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
			{50, 215, 100, 50, true, false, false, false, 0, 0, 100, color.RGBA{139, 69, 19, 255}},
			{400, 215, 100, 50, true, false, false, true, 330, 75, 100, color.RGBA{139, 69, 19, 255}},
			{250, 50, 100, 50, true, false, false, false, 0, 0, 100, color.RGBA{139, 69, 19, 255}},
			{490, 50, 100, 50, true, false, false, false, 0, 0, 100, color.RGBA{139, 69, 19, 255}},
		}
		g.Obstacles = []ObstaclesS{}
		g.Bouncers = []BouncersS{}

	case 7:
		g.Barrels = []BarrelsS{
			{50, 375, 100, 50, true, false, false, false, 0, 0, 100, color.RGBA{139, 69, 19, 255}},
			{350, 215, 100, 50, true, true, false, true, 530, 240, 100, color.RGBA{139, 69, 19, 255}},
			{450, 215, 100, 50, false, true, false, false, 0, 0, 100, color.RGBA{139, 69, 19, 255}},
			{490, 50, 100, 50, true, false, false, false, 0, 0, 100, color.RGBA{139, 69, 19, 255}},
		}
		g.Obstacles = []ObstaclesS{
			{270, 300, 50, 50, false, false, color.RGBA{178, 34, 34, 255}},
		}
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
func (g *Game) DrawTimer(screen *ebiten.Image) error {
	if g.SpaceCNT == 0 {
		ebitenutil.DebugPrintAt(screen, "Time: 0:0:00", 5, 5)
	} else if g.Level != 8 {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Time: %s", time.Since(g.StartTime).Round(10*time.Millisecond)), 5, 5)
	} else if g.State == 4 {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Time: %s", g.endTime), 5, 5)
	}
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
func SaveToDisk(data SaveData, filename string) error {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, bytes, 0644)
}

func LoadFromDisk(filename string) (SaveData, error) {
	var data SaveData

	bytes, err := os.ReadFile(filename)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(bytes, &data)
	return data, err
}

func Within(px, py, rx, ry, rw, rh float64) bool {
	return px >= rx && px <= rx+rw && py >= ry && py <= ry+rh
}

func CmpTime(a, b Score) int {
	if a.Time < b.Time {
		return -1
	} else if a.Time == b.Time {
		return 0
	}

	return 1
}

func (g *Game) Update() error {
	if g.TimeSaveAnimation > 0 && g.State == 1 {
		g.TimeSaveAnimation--
	}
	if g.TimeSaveAnimation == 0 && g.State == 1 {
		g.State++
	}
	if ebiten.IsKeyPressed(ebiten.KeyControlLeft) {
		os.Remove("save.json")
		g.Save = SaveData{}
	}
	if g.State > 1 {
		xC, yC := ebiten.CursorPosition()
		x, y := float64(xC), float64(yC)
		if g.ChangeLevelAnimation && g.Opacity > 255 {
			g.OpacityPlusOrNegative = false
		}
		if !g.OpacityPlusOrNegative {
			g.Opacity -= 2
			g.Generate_Level(g.Level)
		}
		if g.ChangeLevelAnimation && g.Opacity < 255 && g.OpacityPlusOrNegative {
			g.OpacityPlusOrNegative = true
		}
		if g.OpacityPlusOrNegative && g.ChangeLevelAnimation {
			g.Opacity += 2
		}
		if g.State == 4 {
			if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
				if Within(x, y, 170, 183, 312, 63) {
					g.Level = 1
					g.PlayerY = 240
					g.PlayerX = 160
					g.PlayerSpeed = 15
					g.OpacityPlusOrNegative = true
					g.PlayerLife = 3
					g.TimeBeforeLevelDown = -67
					g.Particles = nil
					g.Barrels = nil
					g.Obstacles = nil
					g.Bouncers = nil
					g.SpaceCNT = 0
					g.PlayerBarrel = 0
					g.PlayerMoved = false
					g.StartTime = time.Time{}
					g.endTime = 0
					g.EndTimeIsSet = false
					g.Opacity = 0
					g.ChangeLevelAnimation = false
					g.Generate_Level(g.Level)
					g.State = 2
				}
			}
		}
		if g.ChangeLevelAnimation && g.Opacity <= 0 {
			g.ChangeLevelAnimation = false
			g.Opacity = 0
			if g.Level != 7 {
				g.PlayerX = 160
				g.PlayerY = 240
			} else {
				g.PlayerX = 160
				g.PlayerY = 400
			}
			g.PlayerSpeed = 15
			g.OpacityPlusOrNegative = true
		}
		if g.State == 4 && !g.EndTimeIsSet {
			g.endTime = time.Since(g.StartTime).Round(10 * time.Millisecond)
			g.EndTimeIsSet = true
		}
		if g.State == 4 && !g.DataAreSave {
			// Ajouter le temps du joueur
			g.Save.Top5 = append(g.Save.Top5, Score{
				Time:     g.endTime,
				UserName: g.currentUserName,
			})

			// Trier et garder seulement les 5 meilleurs
			slices.SortFunc(g.Save.Top5, CmpTime)
			if len(g.Save.Top5) > 5 {
				g.Save.Top5 = g.Save.Top5[:5]
			}

			// Mettre à jour le joueur
			SaveToDisk(g.Save, "save.json")
			g.DataAreSave = true

		}
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
		if g.TimeBeforeLevelDown == 0 && !g.ChangeLevelAnimation {
			g.Level--
			if g.Level != 7 {
				g.PlayerX = 160
				g.PlayerY = 240
			} else {
				g.PlayerX = 160
				g.PlayerY = 400
			}
			g.PlayerSpeed = 15
			defer g.Generate_Level(g.Level)
			g.TimeBeforeLevelDown = -67
		}
		// --- BORDS DE L'ÉCRAN ---
		if (g.PlayerX-PlayerR < 0 ||
			g.PlayerX+PlayerR > 640 ||
			g.PlayerY-PlayerR < 0 ||
			g.PlayerY+PlayerR > 480) &&
			!g.ChangeLevelAnimation {

			if g.Level != 7 {
				g.PlayerX = 160
				g.PlayerY = 240
			} else {
				g.PlayerX = 160
				g.PlayerY = 400
			}
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
				if b.CoolDown <= 0 && !g.ChangeLevelAnimation {
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
			}
			if CircleRectCollision(g.PlayerX, g.PlayerY, PlayerR, o.x, o.y, o.w, o.h) && !g.ChangeLevelAnimation {
				if g.Level != 7 {
					g.PlayerX = 160
					g.PlayerY = 240
				} else {
					g.PlayerX = 160
					g.PlayerY = 400
				}
				g.PlayerLife--
				g.PlayerSpeed = 15
			}
		}

		// --- MOUVEMENT DU PLAYER ---
		if ebiten.IsKeyPressed(ebiten.KeySpace) && !g.PlayerMoved {
			g.PlayerMoved = true
			g.SpaceCNT++
			if g.SpaceCNT == 1 {
				g.State++
				g.StartTime = time.Now()
			}
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
				if b.x == 490 && !g.ChangeLevelAnimation {
					g.Level++
					g.ChangeLevelAnimation = true
					if g.Level != 7 {
						g.PlayerX = 160
						g.PlayerY = 240
					} else {
						g.PlayerX = 160
						g.PlayerY = 400
					}
					g.PlayerSpeed = 15
				}
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
	}
	if g.State == 0 {
		// 1. Récupérer le texte tapé cette frame
		typed := ebiten.InputChars()
		if len(typed) > 0 {
			g.currentUserName += string(typed)
		}

		// 2. Supprimer un caractère avec Backspace
		if ebiten.IsKeyPressed(ebiten.KeyBackspace) && len(g.currentUserName) > 0 {
			g.currentUserName = g.currentUserName[:len(g.currentUserName)-1]
		}

		// 3. Fin de la saisie si l’utilisateur appuie sur Enter
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			fmt.Println("Nom terminé:", g.currentUserName)
			g.State++
		}
	}
	if g.Level == 8 {
		g.State = 4
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.State > 1 {
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
		ebitenutil.DebugPrintAt(screen, "version 1.1", 550, 460)
		//draw timer
		g.DrawTimer(screen)
		//draw change level animation
		ebitenutil.DrawRect(screen, 0, 0, 640, 480, color.RGBA{0, 0, 0, uint8(g.Opacity)})
		//draw Restart button
		if g.State == 4 {
			ebitenutil.DrawRect(screen, 170, 183, 312, 63, color.RGBA{0, 255, 0, 255})
			op := &text.DrawOptions{}
			op.GeoM.Translate(float64(171), float64(190))
			op.ColorScale.ScaleWithColor(color.RGBA{255, 255, 255, 255})
			text.Draw(screen, "Restart", &text.GoTextFace{
				Source: mplusFaceSource,
				Size:   45,
			}, op)
		}
	}
	if g.State == 0 {
		op := &text.DrawOptions{}
		op.GeoM.Translate(float64(50), float64(200))
		op.ColorScale.ScaleWithColor(color.RGBA{255, 255, 255, 255})
		text.Draw(screen, fmt.Sprintf("UserName: %s", g.currentUserName), &text.GoTextFace{
			Source: mplusFaceSource,
			Size:   30,
		}, op)
	}
	if g.State == 2 {
		ebitenutil.DrawRect(screen, 50, 300, 500, 100, color.RGBA{0, 255, 0, 255})
		for i, s := range g.Save.Top5 {
			op := &text.DrawOptions{}
			op.GeoM.Translate(float64(100), float64(20*i+300))
			op.ColorScale.ScaleWithColor(color.RGBA{255, 255, 255, 255})
			text.Draw(screen, fmt.Sprintf("%d. %v. - %s", i+1, s.Time, s.UserName), &text.GoTextFace{
				Source: mplusFaceSource,
				Size:   20,
			}, op)
		}
	}
	if g.State == 1 {
		op := &text.DrawOptions{}
		op.GeoM.Translate(float64(50), float64(230))
		op.ColorScale.ScaleWithColor(color.RGBA{255, 255, 255, 255})
		text.Draw(screen, "Downloading UserName...", &text.GoTextFace{
			Source: mplusFaceSource,
			Size:   20,
		}, op)
	}
	if g.State == 1 {
		centerX, centerY := 320.0, 420.0 // position centrale (modifie selon ton UI)
		spacing := 34.0                  // espacement horizontal entre points
		baseR := 8.0                     // rayon de base
		speed := 6.45                    // vitesse des pulsations

		for i := 0; i < 3; i++ {
			// phase décallée pour chaque point
			phase := float64(i) * 0.9

			// sin pour aller de -1..1, on transforme en 0..1 puis en scale 0.6..1.4
			s := (math.Sin(g.TimeSaveAnimation*speed+phase) + 1.0) / 2.0
			scale := 0.6 + 0.8*s // nombre entre 0.6 et 1.4

			r := baseR * scale
			x := centerX + (float64(i)-1.0)*spacing
			// couleur : blanc, tu peux changer
			ebitenutil.DrawCircle(screen, x, centerY, r, color.RGBA{255, 255, 255, 255})
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello World")
	save, err := LoadFromDisk("save.json")
	if err != nil {
		fmt.Println("Aucune sauvegarde trouvée, valeurs par défaut...")
		save = SaveData{
			Top5: []Score{}, // initialiser le slice vide
		}
	}
	g := &Game{
		Save:                  save,
		PlayerY:               240,
		PlayerX:               160,
		PlayerSpeed:           15,
		Level:                 1,
		OpacityPlusOrNegative: true,
		PlayerLife:            3,
		TimeBeforeLevelDown:   -67,
		TimeSaveAnimation:     70,

		Top5Bestplayers: save.Top5, // récupérer le top5 depuis la sauvegarde
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
