package node

import "github.com/hajimehoshi/ebiten/v2"

type Node interface {
	Update()
	Draw(screen *ebiten.Image)
}

func UpdateNode(n Node) {

}

func DrawNode(n Node, screen *ebiten.Image) {
	n.Draw(screen)
}
