package node

import "github.com/hajimehoshi/ebiten/v2"

type Node interface {
	Update()
	Draw(screen *ebiten.Image)
}

func Update(n Node) {
	n.Update()
}

func Draw(n Node, screen *ebiten.Image) {
	n.Draw(screen)
}
