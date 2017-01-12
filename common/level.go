package common

import (
	"engo.io/engo"
	"engo.io/gl"
	"fmt"
)

// Level is a parsed TMX level containing all layers and default Tiled attributes
type Level struct {
	// Orientation is the parsed level orientation from the TMX XML, like orthogonal, isometric, etc.
	Orientation string
	// MapToPoint is a pointer to a function which maps map tile coordinates to display points
	MapToPoint func(*engo.Point) (*engo.Point)
	// PointToMap is a pointer to a function which maps display points to map coordinates
	// return value uses engo.Point rather than int to provide sub-block accuracy.
	PointToMap func(*engo.Point) (*engo.Point)
	// Return the value of the Max bounds point of the map before offset
	MapMaxBounds func() (*engo.Point)
	// RenderOrder is the in Tiled specified TileMap render order, like right-down, right-up, etc.
	RenderOrder string
	width       int
	height      int
	// CanvasOffset transforms the map from 0,0
	Offset engo.Point
	// TileWidth defines the width of each tile in the level
	TileWidth int
	// TileHeight defines the height of each tile in the level
	TileHeight int
	// NextObjectId is the next free Object ID defined by Tiled
	NextObjectId int
	// TileLayers contains all TileLayer of the level
	TileLayers []*TileLayer
	// ImageLayers contains all ImageLayer of the level
	ImageLayers []*ImageLayer
	// ObjectLayers contains all ObjectLayer of the level
	ObjectLayers []*ObjectLayer
	// Properties represents properties about this map.
	Properties map[string]Property
}

// MapToPosition maps a map coordinate with subtile accuracy and return position
// coordinate on the screen. If necessary, calls the setupOrientation function to parse the
// level orientation string and setup the correct mapping.
func (lvl *Level)MapToPosition(mp engo.Point) (*engo.Point, error) {
	if lvl.MapToPoint == nil {
		// Only attempt setup if the function pointer is nil
		var err error
		if err = lvl.setupOrientation(); err != nil {
			return nil, err
		}
	}
	mpc := &engo.Point{mp.X, mp.Y}
	mpc = lvl.MapToPoint(mpc)
	mpc.X += lvl.Offset.X
	mpc.Y += lvl.Offset.Y
	return mpc, nil
}

// PositionToMap maps a screen position coordinate returns a map coordinate with subtile accuracy
// If necessary, calls the setupOrientation function to parse the level orientation string and setup 
// the correct mapping.
func (lvl *Level)PositionToMap(p engo.Point) (*engo.Point, error) {
	if lvl.PointToMap == nil {
		// Only attempt setup if the function pointer is nil
		var err error
		if err = lvl.setupOrientation(); err != nil {
			return nil, err
		}
	}
	mp := &engo.Point{p.X, p.Y}
	mp.X -= lvl.Offset.X
	mp.Y -= lvl.Offset.Y
	return lvl.PointToMap(mp), nil
}

func (l *TileLayer)GetTile(p *engo.Point) (*tile) {
	mp, _ := l.Level.PositionToMap(*p)
	tidx := int(mp.X) + int(mp.Y)*l.width
	return l.Tiles[tidx]
}


// setupOrientation is a function to setup defualt helper functions based on the
// level orientation as defined by tmx.
func (lvl *Level)setupOrientation() error {
	// tile (half-)widths for isometric tilesets.
	tw := float32(lvl.TileWidth)
	th := float32(lvl.TileHeight)
	hw := float32(lvl.TileWidth/2)
	hh := float32(lvl.TileHeight/2)

	// Do the string comparisons once and setup helper functions
	if lvl.Orientation == "orthogonal" {
		lvl.MapToPoint = func(m *engo.Point) (*engo.Point) {
			m.X = m.X * tw
			m.Y = m.Y * th
			return m
		}
		lvl.PointToMap = func(p *engo.Point) (*engo.Point) {
			p.X = p.X / tw
			p.Y = p.Y / th
			return p
		}
		lvl.MapMaxBounds = func() (*engo.Point) {
			return &engo.Point{
				float32(lvl.TileWidth * lvl.width),
				float32(lvl.TileHeight * lvl.height),
			}
		}
	} else if lvl.Orientation == "isometric" {
		lvl.MapToPoint = func(m *engo.Point) (*engo.Point) {
			m.X = (m.X - m.Y) * hw
			m.Y = (m.X + m.Y) * hh
			return m
		}
		lvl.PointToMap = func(p *engo.Point) (*engo.Point) {
			p.X = (p.X + p.Y) / tw
			p.Y = (p.Y - p.X) / th
			return p
		}
		lvl.MapMaxBounds = func() (*engo.Point) {
			return &engo.Point{
				float32(lvl.TileWidth * lvl.width) + float32(lvl.TileWidth/2),
				float32(lvl.TileHeight/2 * lvl.height) + float32(lvl.TileHeight/2),
			}
		}
	} else if lvl.Orientation == "staggered" {
		lvl.MapToPoint = func(m *engo.Point) (*engo.Point) {
			staggerX := float32(0) // no offset on even rows
			if int(m.Y)%2 == 1 {   // odd row?
				staggerX = hw
			}
			m.X = (m.X * tw) + staggerX
			m.Y = m.Y * hh
			return m
		}
		lvl.PointToMap = func(p *engo.Point) (*engo.Point) {
			Y := p.Y
			p.Y = (p.Y - p.X) / th
			staggerX := float32(0) // no offset on even rows
			if int(p.Y)%2 == 1 {   // odd row?
				staggerX = hw
			}
			p.X = (p.X + Y - staggerX) / tw
			return p
		}
		lvl.MapMaxBounds = func() (*engo.Point) {
			return &engo.Point{
				float32(lvl.TileWidth * lvl.width) + float32(lvl.TileWidth/2),
				float32(lvl.TileHeight/2 * lvl.height) + float32(lvl.TileHeight/2),
			}
		}
	} else {
		return fmt.Errorf(
			"Level: Unsupported orientation %v",
			lvl.Orientation)
	}
	return nil
}


// TileLayer contains a list of its tiles plus all default Tiled attributes
type TileLayer struct {
	// Name defines the name of the tile layer given in the TMX XML / Tiled
	Name string
	// Width is the integer width of each tile in this layer
	Width int
	// Height is the integer height of each tile in this layer
	Height int
	// Tiles contains the list of tiles
	Tiles []*tile
	// Level contains a link back to the level we are part of for
	*Level
	// Properties represents properties about this objectLayer
	Properties map[string]Property
}

// ImageLayer contains a list of its images plus all default Tiled attributes
type ImageLayer struct {
	// Name defines the name of the image layer given in the TMX XML / Tiled
	Name string
	// Width is the integer width of each image in this layer
	Width int
	// Height is the integer height of each image in this layer
	Height int
	// Source contains the original image filename
	Source string
	// Images contains the list of all image tiles
	Images []*tile
}

// ObjectLayer contains a list of its standard objects as well as a list of all its polyline objects
type ObjectLayer struct {
	// Name defines the name of the object layer given in the TMX XML / Tiled
	Name string
	// OffSetX is the parsed X offset for the object layer
	OffSetX float32
	// OffSetY is the parsed Y offset for the object layer
	OffSetY float32
	// Objects contains the list of (regular) Object objects
	Objects []*Object
	// PolyObjects contains the list of PolylineObject objects
	PolyObjects []*PolylineObject
	// Properties represents properties about this objectLayer
	Properties map[string]Property
}

// Object is a standard TMX object with all its default Tiled attributes
type Object struct {
	// Id is the unique ID of each object defined by Tiled
	Id int
	// Name defines the name of the object given in Tiled
	Name string
	// Type contains the string type which was given in Tiled
	Type string
	// X holds the X float64 coordinate of the object in the map
	X float64
	// X holds the X float64 coordinate of the object in the map
	Y float64
	// Width is the integer width of the object
	Width int
	// Height is the integer height of the object
	Height int
	// Properties represents properties about this object
	Properties map[string]Property
}

// PolylineObject is a TMX polyline object with all its default Tiled attributes
type PolylineObject struct {
	// Id is the unique ID of each polyline object defined by Tiled
	Id int
	// Name defines the name of the polyline object given in Tiled
	Name string
	// Type contains the string type which was given in Tiled
	Type string
	// X holds the X float64 coordinate of the polyline in the map
	X float64
	// Y holds the Y float64 coordinate of the polyline in the map
	Y float64
	// Points contains the original, unaltered points string from the TMZ XML
	Points string
	// LineBounds is the list of engo.Line objects generated from the points string
	LineBounds []*engo.Line
}

// Bounds returns the level boundaries as an engo.AABB object
func (l *Level) Bounds() engo.AABB {
	max := l.MapMaxBounds()
	max.Add(l.Offset)
	return engo.AABB{
		Min: l.Offset,
		Max: *max,
	}
}

// Width returns the integer width of the level
func (l *Level) Width() int {
	return l.width
}

// Height returns the integer height of the level
func (l *Level) Height() int {
	return l.height
}

// Height returns the integer height of the tile
func (t *tile) Height() float32 {
	return t.Image.Height()
}

// Width returns the integer width of the tile
func (t *tile) Width() float32 {
	return t.Image.Width()
}

// Texture returns the tile's Image texture
func (t *tile) Texture() *gl.Texture {
	return t.Image.id
}

// Close deletes the stored texture of a tile
func (t *tile) Close() {
	t.Image.Close()
}

// View returns the tile's viewport's min and max X & Y
func (t *tile) View() (float32, float32, float32, float32) {
	return t.Image.View()
}

func (t *tile) IsWalkable() bool {
	p, ok := t.Properties["walkable"]
	if ok && p.Value == "true" {
		return true
	}
	return false
}

type tile struct {
	engo.Point
	Image *Texture
	// Properties represents properties about this map.
	Properties map[string]Property
}

type Property struct {
	Type  string
	Value string
}

type tilesheet struct {
	Image    *TextureResource
	TileWidth int
	TileHeight int
	Firstgid int
	Properties map[string]Property
}

type layer struct {
	Name        string
	Width       int
	Height      int
	TileMapping []uint32
	// Properties represents properties about this layer
	Properties map[string]Property
}

func createTileset(lvl *Level, sheets []*tilesheet) []*tile {
	tileset := make([]*tile, 0)

	for _, sheet := range sheets {
		tw := float32(sheet.TileWidth)
		th := float32(sheet.TileHeight)
		setWidth := sheet.Image.Width / tw
		setHeight := sheet.Image.Height / th
		totalTiles := int(setWidth * setHeight)

		for i := 0; i < totalTiles; i++ {
			t := &tile{}
			x := float32(i%int(setWidth)) * tw
			y := float32(i/int(setWidth)) * th

			invTexWidth := 1.0 / float32(sheet.Image.Width)
			invTexHeight := 1.0 / float32(sheet.Image.Height)

			u := float32(x) * invTexWidth
			v := float32(y) * invTexHeight
			u2 := float32(x+tw) * invTexWidth
			v2 := float32(y+th) * invTexHeight
			t.Image = &Texture{
				id:     sheet.Image.Texture,
				width:  tw,
				height: th,
				viewport: engo.AABB{
					engo.Point{u, v},
					engo.Point{u2, v2},
				},
			}
			tileset = append(tileset, t)
		}
	}

	return tileset
}

func createLevelTiles(lvl *Level, layers []*layer, ts []*tile) []*TileLayer {
	var levelTileLayers []*TileLayer

	// Create a TileLayer for each provided layer
	for _, layer := range layers {
		tilemap := make([]*tile, 0)
		tileLayer := &TileLayer{}
		mapping := layer.TileMapping

		// Append tiles to map
		for i := 0; i < lvl.height; i++ {
			for x := 0; x < lvl.width; x++ {
				idx := x + i*lvl.width
				t := &tile{}

				if tileIdx := int(mapping[idx]) - 1; tileIdx >= 0 {
					t.Image = ts[tileIdx].Image
					tp, _ := lvl.MapToPosition(engo.Point{float32(x), float32(i)})
					// Align tiles to bottom for oversize tile layering
					// XXX: bug for unusual draw orders? configurable?
					tp.Y -= t.Image.Height() - float32(lvl.TileHeight)
					t.Point = *tp
				}
				tilemap = append(tilemap, t)
			}
		}

		tileLayer.Name = layer.Name
		tileLayer.Width = layer.Width
		tileLayer.Height = layer.Height
		tileLayer.Tiles = tilemap
		tileLayer.Level = lvl

		levelTileLayers = append(levelTileLayers, tileLayer)
	}

	return levelTileLayers
}
