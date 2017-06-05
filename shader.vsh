#version 410
uniform vec3 u;
in vec4 v;
out vec2 uv;

void main(){
  uv=(v.xy+1)/2;
  gl_Position = v;
}
