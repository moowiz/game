
#version 330

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;
uniform vec3 light;

layout(location = 0) in vec3 vert;

void main() {
    gl_Position = projection * camera * model * vec4(vert, 1);
}
