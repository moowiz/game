#version 330

in vec2 fragUV;

out vec3 color;

uniform sampler2D tex;
uniform vec3 diffuseColor;

void main() {
	color = texture(tex, fragUV).rgb;
}
