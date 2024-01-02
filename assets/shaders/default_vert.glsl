#version 410

uniform mat4 u_projView;

in vec2 Position;
in vec2 TexCoord;
in vec4 Color;

out vec2 vTexCoord;
out vec4 vColor;

void main() {
    gl_Position = u_projView * vec4(Position, 0.0, 1.0);
    vTexCoord = TexCoord;
    vColor = Color;
}
