#version 330

in vec2 fragUV;
in vec3 positionWorld;
in vec3 eyeDirectionInCamera;
in vec3 lightDirectionInCamera;
in vec3 normalInCamera;

out vec3 color;

uniform sampler2D tex;
uniform vec3 light;
uniform vec3 diffuseColor;
uniform vec3 ambientColor;
uniform float LightPower;

void main() {
	vec3 LightColor = vec3(1,1,1);

	// Material properties
	vec3 MaterialDiffuseColor = texture(tex, fragUV).rgb;
	vec3 MaterialAmbientColor = ambientColor * MaterialDiffuseColor;

	// Distance to the light
	float distance = length(light - positionWorld);

	// Normal of the computed fragment, in camera space
	vec3 n = normalize(normalInCamera);
	// Direction of the light (from the fragment to the light)
	vec3 l = normalize(lightDirectionInCamera);
	// Cosine of the angle between the normal and the light direction, 
	// clamped above 0
	//  - light is at the vertical of the triangle -> 1
	//  - light is perpendicular to the triangle -> 0
	//  - light is behind the triangle -> 0
	float cosTheta = clamp( dot( n,l ), 0,1 );
	
	color = 
		// Ambient : simulates indirect lighting
		MaterialAmbientColor +
		// Dif2se : "color" of the object
		MaterialDiffuseColor * LightColor * LightPower * cosTheta / (distance*distance);
}
