
#version 330

layout(location = 0) in vec3 vert;
layout(location = 1) in vec2 UV;

out vec2 fragUV;

uniform mat4 projection;

void main() {
    fragUV = UV;
    gl_Position = projection * vec4(vert, 1);
}
