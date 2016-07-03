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
uniform vec3 specularColor;
uniform int illum;

void main() {
	vec3 LightColor = vec3(1,1,1);
	float LightPower = 50.0f;

	// Material properties
	vec3 MaterialDiffuseColor = texture(tex, fragUV).rgb * diffuseColor;
	vec3 MaterialAmbientColor = ambientColor * diffuseColor;

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
	
	// Eye vector (towards the camera)
	vec3 E = normalize(eyeDirectionInCamera);
	// Direction in which the triangle reflects the light
	vec3 R = reflect(-l,n);
	// Cosine of the angle between the Eye vector and the Reflect vector,
	// clamped to 0
	//  - Looking into the reflection -> 1
	//  - Looking elsewhere -> < 1
	float cosAlpha = clamp( dot( E,R ), 0,1 );
	
	color = 
		// Ambient : simulates indirect lighting
		MaterialAmbientColor +
		// Diffuse : "color" of the object
		diffuseColor * LightColor * LightPower * cosTheta / (distance*distance) +
		// Specular : reflective highlight, like a mirror
		specularColor * LightColor * LightPower * pow(cosAlpha,5) / (distance*distance);
}
