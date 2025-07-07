// scene package provides functionality for switching
// between scenes, such as level (gamplay) and pause menu
// TODO: the whole package?
package scene

import "github.com/pyrdelic/goball/node"

type SceneManager struct {
	Buffer    []*node.Node
	CurrScene int
}

// Current returns a pointer to currently selected scene (node.Node).
func (sm *SceneManager) Current() *node.Node {
	return sm.Buffer[sm.CurrScene]
}

func (sm *SceneManager) SwitchCurr(sceneIndex int) {
	if sm.Buffer[sceneIndex] == nil {

	}
}

func (sm *SceneManager) LoadScene(sceneIndex int, n *node.Node) {

}

func (sm *SceneManager) CloseScene(sceneIndex int) {

}
