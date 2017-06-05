#version 320

// For ease of gamma adjustment at some party when the Mercury gamma slide comes on the beamer!
#define GAMMA 2.2													// AUTOREP

uniform vec3 u;														// {xres, yres, frame count}
uniform sampler2D fr;												// previous frame
out vec3 c;															// output pixel

void main(){
	vec2 uv=gl_FragCoord.xy/u.xy;									// calc uv
	c=texture(fr,uv).xyz/max(1,u.z);								// grab the accumulated frame & divide by number of frames rendered
	
	////////////////////////////////////
	// Do any popro you'd like here :)
	////////////////////////////////////
	
	c=pow(c,vec3(1/GAMMA));											// gamma correction
}
