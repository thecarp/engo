# IsoAdventure Demo

## What does it do?
* Everything adventure does but with a boring isometric tileset
* Use a crude Y row based sorting to sort the isometric draw order

## What are important aspects of the code?
These lines are key in this demo:

* `zo := tile.RenderComponent.Drawable.Height() - float32(levelData.TileHeight)` the Y offset to the top of the tile from the top of its actual tile image
* `tile.RenderComponent.SetZIndex(tileElement.Point.Y + zo)` This is the tile ZIndex mapping that our character sorts into
* `eh := e.RenderComponent.Drawable.Height` Offest to the bottom of an entity because its bottom position is where it sits in relation to the map ZIndex
* e.RenderComponent.SetZIndex(e.SpaceComponent.Position.Y + eh) sort the Hero into the map


