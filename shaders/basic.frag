#version 330

in vec2 fragUV;

out vec3 color;

uniform sampler2D tex;

void main() {
	color = texture(tex, fragUV).rgb;
}
