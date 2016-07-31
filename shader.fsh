#version 410
out vec3 outputColor;
uniform vec3 foo;
uniform vec2 bar;
uniform float baz;
void main() {
    outputColor = vec3(.60,.9,foo.x + bar.y * baz);//texture(tex, fragTexCoord);
}
