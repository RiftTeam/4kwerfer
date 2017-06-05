// Rapunzel (Coming Up?)
// 4k exegfx
// Fell, 2016

#version 410

out vec2 uv;

void main(){
	uv=(gl_Vertex.xy+1)/2;
	gl_Position=gl_Vertex;
}
