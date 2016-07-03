#version 330

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;
uniform vec3 light;

in vec3 vert;
in vec2 UV;
in vec3 normal;

out vec2 fragUV;
out vec3 positionWorld;
out vec3 eyeDirectionInCamera;
out vec3 lightDirectionInCamera;
out vec3 normalInCamera;

void main() {
    fragUV = UV;
    gl_Position = projection * camera * model * vec4(vert, 1);

	positionWorld = (model * vec4(vert, 1)).xyz;

	vec3 vertexToCamera = (camera * model * vec4(vert, 1)).xyz;
	eyeDirectionInCamera = vec3(0, 0, 0) - vertexToCamera;

	vec3 vertexToLightInCamera = (camera * vec4(light, 1)).xyz;
	lightDirectionInCamera = vertexToLightInCamera + eyeDirectionInCamera;

	normalInCamera = (camera * model * vec4(normal, 0)).xyz;
}
