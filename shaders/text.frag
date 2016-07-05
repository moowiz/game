#version 330

in vec2 fragUV;

layout(location = 0) out vec4 color;

uniform sampler2D tex;

void main() {
	color = texture(tex, fragUV).rgba;
}
