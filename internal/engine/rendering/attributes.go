package rendering

const (
	null          = "\x00"
	attrPosition  = "Position" + null
	attrTextCoord = "TexCoord" + null
	attrColor     = "Color" + null
)

var vertexAttributes = []struct {
	location  uint32
	name      string
	numValues int32
}{
	{location: 0, name: attrPosition, numValues: 2},
	{location: 1, name: attrTextCoord, numValues: 2},
	{location: 2, name: attrColor, numValues: 4},
}
