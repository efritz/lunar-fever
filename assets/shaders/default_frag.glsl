#version 410

uniform sampler2D u_texture;

in vec2 vTexCoord;
in vec4 vColor;

out vec4 fragColor;

void main() {
    fragColor = vColor * texture(u_texture, vTexCoord);
}
